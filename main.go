package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Selleo/cli/awscmd"
	"github.com/Selleo/cli/cryptographic"
	"github.com/Selleo/cli/fmtx"
	"github.com/Selleo/cli/generators"
	"github.com/Selleo/cli/random"
	"github.com/Selleo/cli/selleo"
	"github.com/Selleo/cli/shellcmd"
	"github.com/Selleo/cli/templates"
	"github.com/urfave/cli/v2"
	"github.com/wzshiming/ctc"
)

var (
	//go:embed templates
	embededTemplates embed.FS
	// ####go:embed packages/secrets-ui/dist
	// embededUI embed.FS
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
				Name:  "text",
				Usage: "Text operations",
				Subcommands: []*cli.Command{
					{
						Name:  "diff",
						Usage: "Show text differences",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "inline", Usage: "Show diff inline", Required: false, Value: false},
						},
						HelpName: "diff [--inline] <file1> <file2>",
						Action: func(c *cli.Context) error {
							args := c.Args().Slice()
							if len(args) != 2 {
								return fmt.Errorf("diff requires two path arguments")
							}
							a, err := os.ReadFile(args[0])
							if err != nil {
								return err
							}
							b, err := os.ReadFile(args[1])
							if err != nil {
								return err
							}

							fmtx.ContentDiff(c.App.Writer, string(a), string(b), c.Bool("inline"))

							return nil
						},
					},
					{
						Name:  "lorem",
						Usage: "Print Lorem Ipsum",
						Action: func(c *cli.Context) error {
							fmtx.Lorem(c.App.Writer)
							return nil
						},
					},
				},
			},
			{
				Name:  "adr",
				Usage: "Architecture Decision Records",
				Subcommands: []*cli.Command{
					{
						Name: "new",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "title", Usage: "Title of the ADR", Required: true},
						},
						Action: func(c *cli.Context) error {
							gen := generators.ADR{Title: c.String("title")}
							return gen.Render(generators.New(embededTemplates))
						},
					},
				},
			},
			{
				Name:  "rand",
				Usage: "Ranndom generators",
				Subcommands: []*cli.Command{
					{
						Name:  "uuid4",
						Usage: "Generates UUID v4",
						Action: func(c *cli.Context) error {
							fmt.Fprintf(
								c.App.Writer,
								"%s\n",
								random.UUID4(),
							)
							return nil
						},
					},
					{
						Name:  "bytes",
						Usage: "Random bytes",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "size", Usage: "Size of sequence", Required: false, Value: 64},
							&cli.StringFlag{Name: "format", Usage: "Output type (hex|base64|base32|raw)", Required: false, Value: "hex"},
						},
						Action: func(c *cli.Context) error {
							out, err := random.Bytes(c.Int("size"), c.String("format"))
							if err != nil {
								return err
							}
							fmt.Fprintf(
								c.App.Writer,
								"%s\n",
								out,
							)
							return nil
						},
					},
				},
			},
			{
				Name:  "crypto",
				Usage: "Cryptographic functions",
				Subcommands: []*cli.Command{
					{
						Name:  "hmac",
						Usage: "Keyed-Hash Message Authentication Code",
						Subcommands: []*cli.Command{
							{
								Name:  "sha256",
								Usage: "HMAC SHA256 signature",
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "key", Usage: "Secret key to sign a message", Required: true},
									&cli.StringFlag{Name: "message", Usage: "Message", Required: true},
									&cli.StringFlag{Name: "format", Usage: "Output type (hex|base64|base32|raw)", Required: false, Value: "hex"},
								},
								Action: func(c *cli.Context) error {
									out, err := fmtx.OutputFormat(
										cryptographic.HMACWithSHA256(
											[]byte(c.String("key")),
											[]byte(c.String("message")),
										),
										c.String("format"),
									)
									if err != nil {
										return err
									}
									fmt.Fprintf(
										c.App.Writer,
										"%s\n",
										out,
									)
									return nil
								},
							},
						},
					},
				},
			},
			{
				Name:  "gen",
				Usage: "Generate files from templates",
				Subcommands: []*cli.Command{
					{
						Name:  "terraform",
						Usage: "Terraform related generators",
						Subcommands: []*cli.Command{
							{
								Name:  "modules",
								Usage: "Generate latest Terraform modules",
								Flags: []cli.Flag{
									&cli.BoolFlag{Name: "force", Usage: "Overwrite existing files", Required: false, Value: false, Aliases: []string{"f"}},
									&cli.StringFlag{Name: "dir", Usage: "Output directory", Required: false, Aliases: []string{"d"}, Value: "terraform"},
								},
								Action: func(c *cli.Context) error {
									downloader, err := templates.NewDownloader()
									if err != nil {
										return err
									}
									return downloader.Download(c.Context, &templates.DownloadInput{
										RepoOwner:   "Selleo",
										RepoName:    "terraform-aws-modules",
										FolderPath:  "modules",
										Destination: c.String("dir"),
										Force:       c.Bool("force"),
									})
								},
							},
						},
					},
				},
			},
			{
				Name:  "aws",
				Usage: "AWS cloud commands",
				Subcommands: []*cli.Command{
					{
						Name:  "ses",
						Usage: "Simple Email Service",
						Subcommands: []*cli.Command{
							{
								Name:  "password",
								Usage: "Generate SMTP password",
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
									&cli.StringFlag{Name: "secret", Usage: "Secret Access Key", Required: true},
								},
								Action: func(c *cli.Context) error {
									password := awscmd.SESPasswordFromAccessKey(c.String("region"), c.String("secret"))
									fmtx.FGreenln(c.App.Writer, password)
									return nil
								},
							},
						},
					},
					{
						Name:  "dev",
						Usage: "Start a service with SSM secrets",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
							&cli.StringSliceFlag{Name: "path", Usage: "Store Parameters Path (multiple use of flag allowed)", Required: true},
						},
						Action: func(c *cli.Context) error {
							input := &awscmd.InputSSMGetParameters{
								Region: c.String("region"),
								Paths:  c.StringSlice("path"),
							}
							paths := strings.Builder{}
							for _, p := range input.Paths {
								paths.WriteString(p)
								paths.WriteString("/*\n")
							}
							fmt.Fprintf(
								c.App.Writer,
								"%sFetching secrets %s%s\n",
								ctc.ForegroundYellow,
								paths.String(),
								ctc.Reset,
							)
							out, err := awscmd.SSMGetParameters(context.TODO(), input)
							if err != nil {
								return err
							}
							return shellcmd.Pipe(context.TODO(), c.App.Writer, out.Parameters, c.Args().Slice())
						},
					},
					{
						Name:  "export",
						Usage: "Export SSM",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
							&cli.StringSliceFlag{Name: "path", Usage: "Store Parameters Path (multiple use of flag allowed)", Required: true},
						},
						Action: func(c *cli.Context) error {
							input := &awscmd.InputSSMGetParameters{
								Region: c.String("region"),
								Paths:  c.StringSlice("path"),
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
						Name:  "ec2",
						Usage: "EC2 instances",
						Subcommands: []*cli.Command{
							{
								Name:  "ls",
								Usage: "List all EC2 instances in region.",
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
								},
								Action: func(c *cli.Context) error {
									input := &awscmd.InputEc2List{
										Region: c.String("region"),
									}
									out, err := awscmd.Ec2List(context.Background(), input, c.App.Writer)
									if out != nil {
										for _, instance := range out.Instances {
											fmt.Fprintf(
												c.App.Writer,
												"%s %s %s %s %s%15s%s %s %s %s\n",
												instance.ID,
												instance.Name,
												instance.State,
												instance.IPv4private,
												ctc.ForegroundGreen, instance.IPv4, ctc.Reset,
												instance.Type,
												instance.Zone,
												instance.LaunchTime.Format("2006-01-02"),
											)
										}
									}
									if err != nil {
										return err
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
									&cli.StringFlag{Name: "timeout", Usage: `Timeout (time units are "ns", "us" (or "µs"), "ms", "s", "m", "h")`, Required: false, Value: "10m"},
									&cli.StringSliceFlag{Name: "one-off", Usage: "One-off commands (multiple use of flag allowed)", Required: false},
								},
								Action: func(c *cli.Context) error {
									fmt.Fprintf(
										c.App.Writer,
										"%sStarting deployment [timeout=%s]%s\n",
										ctc.ForegroundYellow,
										c.String("timeout"),
										ctc.Reset,
									)
									timeout, err := time.ParseDuration(c.String("timeout"))
									if err != nil {
										return err
									}
									ctx, cancel := context.WithTimeout(c.Context, timeout)
									defer cancel()

									input := &awscmd.InputEcsDeploy{
										Region:      c.String("region"),
										Cluster:     c.String("cluster"),
										Service:     c.String("service"),
										DockerImage: c.String("docker-image"),
										OneOffs:     c.StringSlice("one-off"),
									}
									out, err := awscmd.EcsDeploy(ctx, input, c.App.Writer)
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
									_, err = awscmd.EcsDeployWait(ctx, waitInput, c.App.Writer)
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
									&cli.StringFlag{Name: "timeout", Usage: `Timeout (time units are "ns", "us" (or "µs"), "ms", "s", "m", "h")`, Required: false, Value: "10m"},
								},
								Action: func(c *cli.Context) error {
									timeout, err := time.ParseDuration(c.String("timeout"))
									if err != nil {
										return err
									}
									ctx, cancel := context.WithTimeout(c.Context, timeout)
									defer cancel()

									runTaskInput := &awscmd.InputEcsRunTask{
										Region:        c.String("region"),
										Cluster:       c.String("cluster"),
										Service:       c.String("service"),
										OneOffCommand: c.String("one-off"),
									}
									out, err := awscmd.EcsRunTask(ctx, runTaskInput, c.App.Writer)
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
									_, err = awscmd.EcsTaskWait(ctx, waitInput, c.App.Writer)
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
							{
								Name:  "restart",
								Usage: "Restarts service",
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
									&cli.StringFlag{Name: "cluster", Usage: "ECS cluster ID", Required: true},
									&cli.StringFlag{Name: "service", Usage: "ECS service ID", Required: true},
									&cli.StringFlag{Name: "timeout", Usage: `Timeout (time units are "ns", "us" (or "µs"), "ms", "s", "m", "h")`, Required: false, Value: "10m"},
								},
								Action: func(c *cli.Context) error {
									timeout, err := time.ParseDuration(c.String("timeout"))
									if err != nil {
										return err
									}
									ctx, cancel := context.WithTimeout(c.Context, timeout)
									defer cancel()

									restartTaskInput := &awscmd.InputEcsRestartTask{
										Region:  c.String("region"),
										Cluster: c.String("cluster"),
										Service: c.String("service"),
									}
									out, err := awscmd.EcsRestartTask(ctx, restartTaskInput, c.App.Writer)
									if err != nil {
										return err
									}
									waitInput := &awscmd.InputEcsDeployWait{
										Region:       c.String("region"),
										Cluster:      c.String("cluster"),
										Service:      c.String("service"),
										DeploymentID: out.MonitoredDeploymentID,
									}
									_, err = awscmd.EcsDeployWait(ctx, waitInput, c.App.Writer)
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
