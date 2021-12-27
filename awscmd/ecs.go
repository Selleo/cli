package awscmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type InputEcsDeploy struct {
	Region      string
	Cluster     string
	Service     string
	DockerImage string
}

type OuputEcsDeploy struct {
	Service string
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

	return &OuputEcsDeploy{
		Service: *updateOut.Service.ServiceName,
	}, nil
}
