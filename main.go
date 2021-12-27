package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Selleo/cli/selleo"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/urfave/cli/v2"
	"github.com/wzshiming/ctc"
)

type AwsEcsDeployInput struct {
	Region      *string
	Cluster     *string
	Service     *string
	DockerImage *string
}

type AwsSecretsManagerExportInput struct {
	Region *string
	ID     *string
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Prints CLI verison",
				Action: func(c *cli.Context) error {
					fmt.Fprintf(c.App.Writer, "%s\n", selleo.Version)
					return nil
				},
			},
			{
				Name:  "aws",
				Usage: "AWS cloud commands",
				Subcommands: []*cli.Command{
					{
						Name:  "secrets",
						Usage: "Secrets manager",
						Subcommands: []*cli.Command{
							{
								Name:  "export",
								Usage: "Export secrets for shell",
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
									&cli.StringFlag{Name: "id", Usage: "Secrets ID", Required: true},
								},
								Action: func(c *cli.Context) error {
									actionInput := AwsSecretsManagerExportInput{
										Region: aws.String(c.String("region")),
										ID:     aws.String(c.String("id")),
									}
									sess, err := session.NewSession(&aws.Config{Region: actionInput.Region})
									if err != nil {
										return fmt.Errorf("Failed to initiate new session: %w", err)
									}

									svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(*actionInput.Region))
									input := &secretsmanager.GetSecretValueInput{
										SecretId:     actionInput.ID,
										VersionStage: aws.String("AWSCURRENT"),
									}
									result, err := svc.GetSecretValue(input)
									if err != nil {
										return fmt.Errorf("Failed to get secret value: %w", err)
									}
									if result.SecretString != nil {
										var kv map[string]string
										err = json.NewDecoder(strings.NewReader(*result.SecretString)).Decode(&kv)
										if err != nil {
											return fmt.Errorf("Failed to decode secrets into key-value map: %w", err)
										}
										for k, v := range kv {
											escaped := strings.ReplaceAll(v, `"`, `\"`)
											fmt.Fprint(c.App.Writer, "export ", k, "=", `"`, escaped, `"`, "\n")
										}
									} else {
										return fmt.Errorf("Failed to fetch secret value (SecretString is empty): %w", err)
									}
									return nil
								},
							},
						},
					},
					{
						Name:  "ecs",
						Usage: "Elastic Container Service",
						Subcommands: []*cli.Command{
							{
								Name:  "deploy",
								Usage: "Deploy new image to service. This will replace all container tasks.",
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
									&cli.StringFlag{Name: "cluster", Usage: "ECS cluster ID", Required: true},
									&cli.StringFlag{Name: "service", Usage: "ECS service ID", Required: true},
									&cli.StringFlag{Name: "docker-image", Usage: "Docker image to replace task definition with", Required: true},
								},
								Action: func(c *cli.Context) error {
									actionInput := AwsEcsDeployInput{
										Region:      aws.String(c.String("region")),
										Cluster:     aws.String(c.String("cluster")),
										Service:     aws.String(c.String("service")),
										DockerImage: aws.String(c.String("docker-image")),
									}

									ses, err := session.NewSession(&aws.Config{Region: actionInput.Region})
									if err != nil {
										return fmt.Errorf("Failed to initiate new session: %w", err)
									}
									svc := ecs.New(ses)

									// 1. fetch running task
									serviceOut, err := svc.DescribeServicesWithContext(context.TODO(), &ecs.DescribeServicesInput{
										Cluster:  actionInput.Cluster,
										Services: []*string{actionInput.Service},
									})

									if err != nil {
										return fmt.Errorf("Failed to fetch current running task: %w", err)
									}

									if len(serviceOut.Services) == 0 {
										return fmt.Errorf("No service definition found")
									}
									if len(serviceOut.Services) != 1 {
										// this should not happen because we defined only 1 service in input but stays for sanity check
										return fmt.Errorf("Ambigious match, found more than 1 service")
									}

									var taskDefinition *string
									for _, s := range serviceOut.Services {
										taskDefinition = s.TaskDefinition
									}

									// 2. build new task definition
									taskDefinitionOut, err := svc.DescribeTaskDefinitionWithContext(context.TODO(), &ecs.DescribeTaskDefinitionInput{
										TaskDefinition: taskDefinition,
									})

									containerDefinitions := []*ecs.ContainerDefinition{}
									for _, container := range taskDefinitionOut.TaskDefinition.ContainerDefinitions {
										container.Image = actionInput.DockerImage
										containerDefinitions = append(containerDefinitions, container)
									}

									// 3. register new task revision
									registerOut, err := svc.RegisterTaskDefinitionWithContext(context.TODO(), &ecs.RegisterTaskDefinitionInput{
										ContainerDefinitions: containerDefinitions,
										Family:               taskDefinitionOut.TaskDefinition.Family,
									})
									if err != nil {
										return fmt.Errorf("Failed to register new task revision: %w", err)
									}
									arn := registerOut.TaskDefinition.TaskDefinitionArn

									// 4. update ecs service with new task arn
									updateOut, err := svc.UpdateServiceWithContext(context.TODO(), &ecs.UpdateServiceInput{
										Cluster:        actionInput.Cluster,
										Service:        actionInput.Service,
										TaskDefinition: arn,
									})
									if err != nil {
										return fmt.Errorf("Failed to update service with new task revision: %w", err)
									}

									fmt.Fprintf(c.App.Writer, "%sNew deployment for service `%s` created%s\n", ctc.ForegroundGreen, *updateOut.Service.ServiceName, ctc.Reset)

									return nil
								},
							},
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(app.ErrWriter, "%s%v%s\n", ctc.ForegroundRed, err, ctc.Reset)
		os.Exit(1)
	}
}
