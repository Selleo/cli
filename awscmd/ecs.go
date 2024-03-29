package awscmd

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/wzshiming/ctc"
)

type InputEcsDeploy struct {
	Region      string
	Cluster     string
	Service     string
	DockerImage string
	OneOffs     []string
}

type OuputEcsDeploy struct {
	Service               string
	MonitoredDeploymentID string
}

func EcsDeploy(ctx context.Context, input *InputEcsDeploy, w io.Writer) (*OuputEcsDeploy, error) {
	sess, err := NewSession(input.Region)
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)
	identity := sts.New(sess)

	caller, err := identity.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	taskDef, err := ecsTaskDefinitionFromService(svc, ctx, input.Cluster, input.Service, w)
	if err != nil {
		return nil, err
	}
	taskName := *taskDef.Family

	fmt.Fprintf(w, "Registering new task [%s]\n", idOneOff(taskName, ""))
	newTask, err := ecsRegisterTaskDefinition(svc, ctx, &inputRegisterTaskDefinition{
		TaskDefinitionOrARN: *taskDef.TaskDefinitionArn,
		DockerImage:         input.DockerImage,
	}, w)
	if err != nil {
		return nil, err
	}
	for _, cmd := range input.OneOffs {
		fmt.Fprintf(w, "Registering new one-off task [%s]\n", idOneOff(taskName, cmd))
		taskDef := fmt.Sprintf(
			"arn:aws:ecs:%s:%s:task-definition/%s",
			input.Region,
			*caller.Account,
			idOneOff(taskName, cmd),
		)
		_, err := ecsRegisterTaskDefinition(svc, ctx, &inputRegisterTaskDefinition{
			TaskDefinitionOrARN: taskDef,
			DockerImage:         input.DockerImage,
		}, w)
		if err != nil {
			return nil, err
		}
	}

	// 4. update ecs service with new task arn
	fmt.Fprintf(w, "Updating service\n")
	updateOut, err := svc.UpdateServiceWithContext(ctx, &ecs.UpdateServiceInput{
		Cluster:        aws.String(input.Cluster),
		Service:        aws.String(input.Service),
		TaskDefinition: aws.String(newTask.ARN),
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to update service with new task revision: %w", err)
	}
	var monitoredDeployment *ecs.Deployment = nil
	for _, deployment := range updateOut.Service.Deployments {
		if *deployment.Status == "PRIMARY" {
			monitoredDeployment = deployment
			break
		}
	}

	if monitoredDeployment == nil {
		return &OuputEcsDeploy{
			Service: *updateOut.Service.ServiceName,
		}, fmt.Errorf("Service %s deployed but couldn't fetch primary deployment status", *updateOut.Service.ServiceName)
	}

	return &OuputEcsDeploy{
		Service:               *updateOut.Service.ServiceName,
		MonitoredDeploymentID: *monitoredDeployment.Id,
	}, nil
}

type InputEcsDeployWait struct {
	Region       string
	Cluster      string
	Service      string
	DeploymentID string
}

type OuputEcsDeployWait struct {
}

func EcsDeployWait(ctx context.Context, input *InputEcsDeployWait, w io.Writer) (*OuputEcsDeployWait, error) {
	sess, err := NewSession(input.Region)
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)

	attempt := 1

	for {
		serviceOut, err := svc.DescribeServicesWithContext(ctx, &ecs.DescribeServicesInput{
			Cluster:  aws.String(input.Cluster),
			Services: []*string{aws.String(input.Service)},
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch service: %w", err)
		}
		if len(serviceOut.Services) == 0 {
			return nil, fmt.Errorf("No service definition found")
		}
		if len(serviceOut.Services) != 1 {
			// this should not happen because we defined only 1 service in input but stays for sanity check
			return nil, fmt.Errorf("Ambigious match, found more than 1 service")
		}

		var monitoredDeployment *ecs.Deployment = nil
		for _, deployment := range serviceOut.Services[0].Deployments {
			if *deployment.Id == input.DeploymentID {
				monitoredDeployment = deployment
				break
			}
		}

		if monitoredDeployment == nil {
			return nil, fmt.Errorf("Failed to monitor service deployment")
		}

		fmt.Fprintf(w, "%s%d%s Running | %s%d%s Pending | %s%d%s Desired (Check %d, retrying in 10s)\n",
			ctc.ForegroundGreen, *monitoredDeployment.RunningCount, ctc.Reset,
			ctc.ForegroundYellow, *monitoredDeployment.PendingCount, ctc.Reset,
			ctc.ForegroundRed, *monitoredDeployment.DesiredCount, ctc.Reset,
			attempt,
		)

		if *monitoredDeployment.RolloutState != "IN_PROGRESS" {
			completed := (*monitoredDeployment.RolloutState == "COMPLETED")
			if completed {
				return &OuputEcsDeployWait{}, nil
			} else {
				return nil, fmt.Errorf("Deployment failed")
			}
		}

		time.Sleep(10 * time.Second)
		attempt++
	}
}

func EcsRunTask(ctx context.Context, input *InputEcsRunTask, w io.Writer) (*OutputEcsRunTask, error) {
	sess, err := NewSession(input.Region)
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)

	taskDef, err := ecsTaskDefinitionFromService(svc, ctx, input.Cluster, input.Service, w)
	taskName := idOneOff(*taskDef.Family, input.OneOffCommand)

	out, err := svc.RunTask(&ecs.RunTaskInput{
		Cluster:        aws.String(input.Cluster),
		TaskDefinition: aws.String(taskName),
		Count:          aws.Int64(1),
	})
	if err != nil {
		return &OutputEcsRunTask{}, fmt.Errorf("Can't runTask: %v", err)
	}
	err = failuresToError(out.Failures)
	if err != nil {
		return &OutputEcsRunTask{}, fmt.Errorf("Can't runTask: %v", err)
	}

	arn := *out.Tasks[0].TaskArn
	splits := strings.Split(arn, "/")
	id := splits[len(splits)-1]

	return &OutputEcsRunTask{
		ARN: arn,
		ID:  id,
	}, nil
}

type InputEcsRunTask struct {
	Region        string
	Cluster       string
	Service       string
	OneOffCommand string
}

type OutputEcsRunTask struct {
	ARN string
	ID  string
}

func EcsRestartTask(ctx context.Context, input *InputEcsRestartTask, w io.Writer) (*OutputEcsRestartTask, error) {
	sess, err := NewSession(input.Region)
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)

	// 1. update ecs service with new task arn
	fmt.Fprintf(w, "Updating service\n")
	updateOut, err := svc.UpdateServiceWithContext(ctx, &ecs.UpdateServiceInput{
		Cluster:            aws.String(input.Cluster),
		Service:            aws.String(input.Service),
		ForceNewDeployment: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to restart service: %w", err)
	}
	var monitoredDeployment *ecs.Deployment = nil
	for _, deployment := range updateOut.Service.Deployments {
		if *deployment.Status == "PRIMARY" {
			monitoredDeployment = deployment
			break
		}
	}

	if monitoredDeployment == nil {
		return &OutputEcsRestartTask{
			Service: *updateOut.Service.ServiceName,
		}, fmt.Errorf("Service %s deployed but couldn't fetch primary deployment status", *updateOut.Service.ServiceName)
	}

	return &OutputEcsRestartTask{
		Service:               *updateOut.Service.ServiceName,
		MonitoredDeploymentID: *monitoredDeployment.Id,
	}, nil
}

type InputEcsRestartTask struct {
	Region        string
	Cluster       string
	Service       string
	OneOffCommand string
}

type OutputEcsRestartTask struct {
	Service               string
	MonitoredDeploymentID string
}

func EcsTaskWait(ctx context.Context, input *InputEcsTaskWait, w io.Writer) (*OuputEcsTaskWait, error) {
	sess, err := NewSession(input.Region)
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)

	attempt := 1
	for {
		result, err := svc.DescribeTasksWithContext(ctx, &ecs.DescribeTasksInput{
			Cluster: aws.String(input.Cluster),
			Tasks:   []*string{aws.String(input.ARN)},
		})
		if err != nil {
			return nil, err
		}

		task := result.Tasks[0]
		status := *task.LastStatus

		if status == "STOPPED" {
			success := true
			for _, c := range task.Containers {
				fmt.Fprintf(w, "Container `%s` exit code is %d\n",
					*c.Name,
					*c.ExitCode,
				)
				if *c.ExitCode != 0 {
					success = false
				}
			}
			if !success {
				fmt.Fprintf(w, "%sStopped reason: %s%s\n",
					ctc.ForegroundRed, *task.StoppedReason, ctc.Reset,
				)
				return nil, fmt.Errorf("Task exited with non-zero code")
			}
			break
		} else {
			fmt.Fprintf(w, "Waiting to finish (`%s%s%s`).. (Check %d, retrying in 3s)\n",
				ctc.ForegroundYellow, status, ctc.Reset,
				attempt,
			)
		}

		time.Sleep(3 * time.Second)
	}

	return &OuputEcsTaskWait{}, nil
}

type InputEcsTaskWait struct {
	Region  string
	Cluster string
	ARN     string
}

type OuputEcsTaskWait struct {
}

func ecsRegisterTaskDefinition(svc *ecs.ECS, ctx context.Context, input *inputRegisterTaskDefinition, w io.Writer) (*outputRegisterTaskDefinition, error) {
	fmt.Fprintf(w, "  Fetching task definition\n")
	taskDefinitionOut, err := svc.DescribeTaskDefinitionWithContext(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(input.TaskDefinitionOrARN),
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch task definition: %v", err)
	}

	containerDefinitions := []*ecs.ContainerDefinition{}
	for _, container := range taskDefinitionOut.TaskDefinition.ContainerDefinitions {
		container.Image = aws.String(input.DockerImage)
		containerDefinitions = append(containerDefinitions, container)
	}

	// 3. register new task revision
	fmt.Fprintf(w, "  Registering new task revision with new docker image\n")
	registerOut, err := svc.RegisterTaskDefinitionWithContext(ctx, &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    containerDefinitions,
		Cpu:                     taskDefinitionOut.TaskDefinition.Cpu,
		EphemeralStorage:        taskDefinitionOut.TaskDefinition.EphemeralStorage,
		ExecutionRoleArn:        taskDefinitionOut.TaskDefinition.ExecutionRoleArn,
		Family:                  taskDefinitionOut.TaskDefinition.Family,
		InferenceAccelerators:   taskDefinitionOut.TaskDefinition.InferenceAccelerators,
		IpcMode:                 taskDefinitionOut.TaskDefinition.IpcMode,
		Memory:                  taskDefinitionOut.TaskDefinition.Memory,
		NetworkMode:             taskDefinitionOut.TaskDefinition.NetworkMode,
		PidMode:                 taskDefinitionOut.TaskDefinition.PidMode,
		PlacementConstraints:    taskDefinitionOut.TaskDefinition.PlacementConstraints,
		ProxyConfiguration:      taskDefinitionOut.TaskDefinition.ProxyConfiguration,
		RequiresCompatibilities: taskDefinitionOut.TaskDefinition.RequiresCompatibilities,
		RuntimePlatform:         taskDefinitionOut.TaskDefinition.RuntimePlatform,
		TaskRoleArn:             taskDefinitionOut.TaskDefinition.TaskRoleArn,
		Volumes:                 taskDefinitionOut.TaskDefinition.Volumes,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to register new task revision: %w", err)
	}

	return &outputRegisterTaskDefinition{
		ARN: *registerOut.TaskDefinition.TaskDefinitionArn,
	}, nil
}

func ecsTaskDefinitionFromService(svc *ecs.ECS, ctx context.Context, cluster string, service string, w io.Writer) (*ecs.TaskDefinition, error) {
	// 1. fetch running task
	fmt.Fprintf(w, "Fetching service definition\n")
	serviceOut, err := svc.DescribeServicesWithContext(ctx, &ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []*string{aws.String(service)},
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch current running task: %w", err)
	}

	if len(serviceOut.Services) == 0 {
		return nil, fmt.Errorf("No service definition found")
	}
	if len(serviceOut.Services) != 1 {
		// this should not happen because we defined only 1 service in input but stays for sanity check
		return nil, fmt.Errorf("Ambigious match, found more than 1 service")
	}

	var taskDefinition *string
	for _, s := range serviceOut.Services {
		taskDefinition = s.TaskDefinition
	}

	// 2. fetch task definition family
	fmt.Fprintf(w, "Fetching task definition\n")
	task, err := svc.DescribeTaskDefinitionWithContext(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: taskDefinition,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to describe task definition: %v", err)
	}

	return task.TaskDefinition, nil
}

type inputRegisterTaskDefinition struct {
	TaskDefinitionOrARN string
	DockerImage         string
}

type outputRegisterTaskDefinition struct {
	ARN string
}

// convention used in terraform modules
func idOneOff(prefix string, command string) string {
	if command == "" {
		return prefix
	} else {
		return fmt.Sprint(prefix, "-", command)
	}
}

func failuresToError(failures []*ecs.Failure) error {
	if len(failures) > 0 {
		failure := failures[0]
		detail := ""
		if failure.Detail != nil {
			detail = fmt.Sprint(" ", *failure.Detail, " ")
		}
		return fmt.Errorf(
			"%s%s (%s)",
			*failure.Reason,
			detail,
			*failure.Arn,
		)
	}
	return nil
}
