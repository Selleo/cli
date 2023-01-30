package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Selleo/cli/awscmd"
	"github.com/Selleo/cli/selleo"
	"github.com/Selleo/cli/shellcmd"
	"github.com/urfave/cli/v2"
	"github.com/wzshiming/ctc"
)

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
						Name:  "export",
						Usage: "Export SSM",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
							&cli.StringFlag{Name: "path", Usage: "SSM Path", Required: true},
						},
						Action: func(c *cli.Context) error {
							input := &awscmd.InputSSMGetParameters{
								Region: c.String("region"),
								Path:   c.String("path"),
							}
							out, err := awscmd.SSMGetParameters(context.TODO(), input)
							if err != nil {
								return err
							}
							shellcmd.KeyValueToExports(c.App.Writer, out.Parameters)

							return nil
						},
					},
					{
						Name:  "configure",
						Usage: "Configure AWS profile",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
							&cli.StringFlag{Name: "profile", Usage: "AWS profile", Required: true},
							&cli.StringFlag{Name: "key", Usage: "Access key ID", Required: true},
							&cli.StringFlag{Name: "secret", Usage: "Secret access key", Required: true},
						},
						Action: func(c *cli.Context) error {
							input := &awscmd.InputConfigure{
								Region:    c.String("region"),
								Profile:   c.String("profile"),
								AccessKey: c.String("key"),
								SecretKey: c.String("secret"),
							}
							_, err := awscmd.Configure(context.TODO(), input)
							if err != nil {
								return err
							}

							return nil
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
									&cli.StringSliceFlag{Name: "one-off", Usage: "One-off commands (multiple use of flag allowed)", Required: false},
								},
								Action: func(c *cli.Context) error {
									input := &awscmd.InputEcsDeploy{
										Region:      c.String("region"),
										Cluster:     c.String("cluster"),
										Service:     c.String("service"),
										DockerImage: c.String("docker-image"),
										OneOffs:     c.StringSlice("one-off"),
									}
									out, err := awscmd.EcsDeploy(context.TODO(), input, c.App.Writer)
									if out != nil {
										fmt.Fprintf(
											c.App.Writer,
											"%sNew deployment for service `%s` created%s\n",
											ctc.ForegroundYellow,
											out.Service,
											ctc.Reset,
										)
									}
									if err != nil {
										return err
									}

									waitInput := &awscmd.InputEcsDeployWait{
										Region:       c.String("region"),
										Cluster:      c.String("cluster"),
										Service:      c.String("service"),
										DeploymentID: out.MonitoredDeploymentID,
									}
									_, err = awscmd.EcsDeployWait(context.TODO(), waitInput, c.App.Writer)
									if err != nil {
										return err
									}

									fmt.Fprintf(
										c.App.Writer,
										"%sDeployment for service `%s` reached stable state%s\n",
										ctc.ForegroundGreen,
										out.Service,
										ctc.Reset,
									)

									return nil
								},
							},
							{
								Name:  "run",
								Usage: "Starts a new one-off task",
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
									&cli.StringFlag{Name: "cluster", Usage: "ECS cluster ID", Required: true},
									&cli.StringFlag{Name: "service", Usage: "ECS service ID", Required: true},
									&cli.StringFlag{Name: "one-off", Usage: "One-off command to run", Required: true},
								},
								Action: func(c *cli.Context) error {
									runTaskInput := &awscmd.InputEcsRunTask{
										Region:        c.String("region"),
										Cluster:       c.String("cluster"),
										Service:       c.String("service"),
										OneOffCommand: c.String("one-off"),
									}
									out, err := awscmd.EcsRunTask(context.TODO(), runTaskInput, c.App.Writer)
									if err != nil {
										return err
									}
									fmt.Fprintf(
										c.App.Writer,
										"%sNew task started: `%s`%s\n",
										ctc.ForegroundGreen,
										out.ID,
										ctc.Reset,
									)
									waitInput := &awscmd.InputEcsTaskWait{
										Region:  c.String("region"),
										Cluster: c.String("cluster"),
										ARN:     out.ARN,
									}
									_, err = awscmd.EcsTaskWait(context.TODO(), waitInput, c.App.Writer)
									if err != nil {
										return err
									}

									fmt.Fprintf(
										c.App.Writer,
										"%sTask `%s` has stopped%s\n",
										ctc.ForegroundGreen,
										out.ID,
										ctc.Reset,
									)

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
