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
									input := &awscmd.InputSecretsAll{
										Region: c.String("region"),
										ID:     c.String("id"),
									}
									out, err := awscmd.SecretsAll(context.TODO(), input)
									if err != nil {
										return err
									}

									shellcmd.KeyValueToExports(c.App.Writer, out.Secrets)

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
									input := &awscmd.InputEcsDeploy{
										Region:      c.String("region"),
										Cluster:     c.String("cluster"),
										Service:     c.String("service"),
										DockerImage: c.String("docker-image"),
									}
									out, err := awscmd.EcsDeploy(context.TODO(), input)
									if err != nil {
										return err
									}

									fmt.Fprintf(
										c.App.Writer,
										"%sNew deployment for service `%s` created%s\n",
										ctc.ForegroundGreen,
										out.Service,
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
