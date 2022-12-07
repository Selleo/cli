package awscmd

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
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
	Service             string
	PrimaryDeploymentID string
}

func EcsDeploy(ctx context.Context, input *InputEcsDeploy) (*OuputEcsDeploy, error) {
	sess, err := NewSession(input.Region)
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)

	// 1. fetch running task
	serviceOut, err := svc.DescribeServicesWithContext(ctx, &ecs.DescribeServicesInput{
		Cluster:  aws.String(input.Cluster),
		Services: []*string{aws.String(input.Service)},
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

	// 2. build new task definition
	taskDefinitionOut, err := svc.DescribeTaskDefinitionWithContext(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: taskDefinition,
	})

	containerDefinitions := []*ecs.ContainerDefinition{}
	for _, container := range taskDefinitionOut.TaskDefinition.ContainerDefinitions {
		container.Image = aws.String(input.DockerImage)
		containerDefinitions = append(containerDefinitions, container)
	}

	// 3. register new task revision
	registerOut, err := svc.RegisterTaskDefinitionWithContext(ctx, &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: containerDefinitions,
		Family:               taskDefinitionOut.TaskDefinition.Family,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to register new task revision: %w", err)
	}
	arn := registerOut.TaskDefinition.TaskDefinitionArn

	// 4. update ecs service with new task arn
	updateOut, err := svc.UpdateServiceWithContext(ctx, &ecs.UpdateServiceInput{
		Cluster:        aws.String(input.Cluster),
		Service:        aws.String(input.Service),
		TaskDefinition: arn,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to update service with new task revision: %w", err)
	}
	var primaryDeployment *ecs.Deployment = nil
	for _, deployment := range updateOut.Service.Deployments {
		if *deployment.Status == "PRIMARY" {
			primaryDeployment = deployment
			break
		}
	}

	if primaryDeployment == nil {
		return &OuputEcsDeploy{
			Service: *updateOut.Service.ServiceName,
		}, fmt.Errorf("Service %s deployed but couldn't fetch primary deployment status", *updateOut.Service.ServiceName)
	}

	return &OuputEcsDeploy{
		Service:             *updateOut.Service.ServiceName,
		PrimaryDeploymentID: *primaryDeployment.Id,
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

		var primaryDeployment *ecs.Deployment = nil
		for _, deployment := range serviceOut.Services[0].Deployments {
			if *deployment.Id == input.DeploymentID {
				primaryDeployment = deployment
				break
			}
		}

		if primaryDeployment == nil {
			return nil, fmt.Errorf("Failed to monitor service deployment")
		}

		fmt.Fprintf(w, "%s%d%s Running | %s%d%s Pending | %s%d%s Desired (Attempt %d, retrying in 10s)\n",
			ctc.ForegroundGreen, *primaryDeployment.RunningCount, ctc.Reset,
			ctc.ForegroundYellow, *primaryDeployment.PendingCount, ctc.Reset,
			ctc.ForegroundRed, *primaryDeployment.DesiredCount, ctc.Reset,
			attempt,
		)

		if *primaryDeployment.RolloutState != "IN_PROGRESS" {
			completed := (*primaryDeployment.RolloutState == "COMPLETED")
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
