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
						Name:  "secrets",
						Usage: "Secrets manager",
						Subcommands: []*cli.Command{
							{
								Name:  "kv",
								Usage: "Key-value secrets",
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
									{
										Name:  "set",
										Usage: "Sets key value secret",
										Flags: []cli.Flag{
											&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
											&cli.StringFlag{Name: "id", Usage: "Secrets ID", Required: true},
										},
										ArgsUsage: "KEY VALUE",
										Action: func(c *cli.Context) error {
											if c.NArg() != 2 {
												return fmt.Errorf("Invalid number of arguments. Missing KEY and VALUE.")
											}
											k, v := c.Args().Get(0), c.Args().Get(1)

											input := &awscmd.InputSecretsPut{
												Region:     c.String("region"),
												ID:         c.String("id"),
												NewSecrets: map[string]string{k: v},
											}
											_, err := awscmd.SecretsPut(context.TODO(), input)
											if err != nil {
												return err
											}

											return nil
										},
									},
								},
							},
							{
								Name:  "bin",
								Usage: "binary secrets",
								Subcommands: []*cli.Command{
									{
										Name:  "get",
										Usage: "Get binary secret",
										Flags: []cli.Flag{
											&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
											&cli.StringFlag{Name: "id", Usage: "Secrets ID", Required: true},
										},
										Action: func(c *cli.Context) error {
											input := &awscmd.InputSecretsBinaryGet{
												Region: c.String("region"),
												ID:     c.String("id"),
											}
											out, err := awscmd.SecretsBinaryGet(context.TODO(), input)
											if err != nil {
												return err
											}

											c.App.Writer.Write(out.Bytes)

											return nil
										},
									},
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
										Region:      c.String("region"),
										Cluster:     c.String("cluster"),
										Service:     c.String("service"),
										DeploymentID: out.PrimaryDeploymentID,
									}
									_, err = awscmd.EcsDeployWait(context.TODO(), waitInput, c.App.Writer)
									if err != nil {
										return err
									}

									fmt.Fprintf(
										c.App.Writer,
										"%sDeployment for service `%s` finished%s\n",
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
