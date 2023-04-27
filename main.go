package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/Selleo/cli/awscmd"
	"github.com/Selleo/cli/generators"
	"github.com/Selleo/cli/selleo"
	"github.com/Selleo/cli/shellcmd"
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
			// {
			// 	Name:  "ui",
			// 	Usage: "Start UI",
			// 	Action: func(c *cli.Context) error {
			// 		return web.UI(c.Context, embededUI)
			// 	},
			// },
			{
				Name:  "gen",
				Usage: "Generate files from templates",
				Subcommands: []*cli.Command{
					{
						Name:  "terraform",
						Usage: "Terraform related generators",
						Subcommands: []*cli.Command{
							{
								Name:  "app",
								Usage: "Generate single app envrionemnt",
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "tf-cloud-org", Usage: "Terraform Cloud organization", Required: true},
									&cli.StringFlag{Name: "tf-cloud-workspace", Usage: "Terraform Cloud workspace", Required: true},
									&cli.StringFlag{Name: "region", Usage: "AWS Region for resources", Required: true},
									&cli.StringFlag{Name: "stage", Usage: "App stage", Required: true},
									&cli.StringFlag{Name: "namespace", Usage: "App namespace", Required: true},
									&cli.StringFlag{Name: "name", Usage: "App ID", Required: true},
									&cli.StringFlag{Name: "domain", Usage: "Domain", Required: true},
									&cli.StringFlag{Name: "subdomain", Usage: "Subdomain", Required: true},
								},
								Action: func(c *cli.Context) error {
									gen := generators.Terraform{
										TerraformCloudOrganization: c.String("tf-cloud-org"),
										TerraformCloudWorkspace:    c.String("tf-cloud-workspace"),
										Region:                     c.String("region"),
										Namespace:                  c.String("namespace"),
										Stage:                      c.String("stage"),
										Name:                       c.String("name"),
										Subdomain:                  c.String("subdomain"),
										Domain:                     c.String("domain"),
										IAMCI:                      fmt.Sprintf("ci-%s-%s", c.String("namespace"), c.String("stage")),
										IAMApp:                     fmt.Sprintf("%s-%s-%s", c.String("namespace"), c.String("stage"), c.String("name")),
										ECSInstanceType:            "t3.medium",
										ECSMinSize:                 1,
										ECSServiceMinMemory:        256,
										ECSServiceMaxMemory:        1024,
										ECSServiceCpu:              1024,
										ECSServicePort:             3000,
										ECSOneOffs:                 []string{},
										LBName:                     fmt.Sprintf("%s-%s-%s", c.String("namespace"), c.String("stage"), c.String("name")),
										DBName:                     c.String("name"),
										DBIdentifier:               c.String("name"),
										DBUser:                     "app",
										DBMultiAZ:                  false,
										DBApplyImmediately:         false,
										BucketName:                 fmt.Sprintf("%s-%s-%s-storage", c.String("namespace"), c.String("stage"), c.String("name")),
									}
									return gen.Render(generators.New(embededTemplates))
								},
							},
						},
					},
					{
						Name:  "docker",
						Usage: "Docker related generators",
						Subcommands: []*cli.Command{
							{
								Name:  "rails",
								Usage: "Generate Dockerfile with entrypoint",
								Action: func(c *cli.Context) error {
									gen := generators.Docker{
										CmdServer: "rails server",
										OneOffs: map[string]string{
											"migrate": "rails db:migrate",
										},
									}
									return gen.Render(generators.New(embededTemplates))
								},
							},
						},
					},
					{
						Name:  "github",
						Usage: "GitHub workflows",
						Subcommands: []*cli.Command{
							{
								Name:  "backend",
								Usage: "Generate GitHub actions for backend",
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "workdir", Usage: "Working directory (source root)", Required: true},
									&cli.StringFlag{Name: "domain", Usage: "Application domain", Required: true},
									&cli.StringFlag{Name: "subdomain", Usage: "Application subdomain", Required: true},
									&cli.StringFlag{Name: "region", Usage: "AWS Region for S3", Required: true},
									&cli.StringFlag{Name: "ecs-cluster", Usage: "ECS cluster name", Required: true},
									&cli.StringFlag{Name: "ecs-service", Usage: "ECS service name", Required: true},
									&cli.StringFlag{Name: "stage", Usage: "Application environment", Required: true},
									/// optional
									&cli.BoolFlag{Name: "tag-release", Usage: "Trigger CI on tag release", Required: false, Value: false},
									&cli.StringSliceFlag{Name: "one-off", Usage: "One-off commands (multiple use of flag allowed)", Required: false},
								},
								Action: func(c *cli.Context) error {
									tpls := generators.New(embededTemplates)
									gen := generators.GitHub{
										CITagTrigger: c.Bool("tag-release"),
										CIBranch:     "main",
										CIWorkingDir: c.String("workdir"),
										Stage:        c.String("stage"),
										Domain:       c.String("domain"),
										Subdomain:    c.String("subdomain"),
										Region:       c.String("region"),
										ECSCluster:   c.String("ecs-cluster"),
										ECSService:   c.String("ecs-service"),
										ECSOneOffs:   c.StringSlice("one-off"),
									}
									if err := gen.RenderBackend(tpls); err != nil {
										return err
									}
									return nil
								},
							},
							{
								Name:  "frontend",
								Usage: "Generate GitHub actions for CDN",
								Flags: []cli.Flag{
									&cli.StringFlag{Name: "workdir", Usage: "Working directory (source root)", Required: true},
									&cli.StringFlag{Name: "domain", Usage: "Application domain", Required: true},
									&cli.StringFlag{Name: "region", Usage: "AWS Region for S3", Required: true},
									&cli.StringFlag{Name: "app_id", Usage: "App ID specified in Terraform", Required: true},
									&cli.StringFlag{Name: "stage", Usage: "Application environment", Required: true},
									/// optional
									&cli.BoolFlag{Name: "tag-release", Usage: "Trigger CI on tag release", Required: false, Value: false},
								},
								Action: func(c *cli.Context) error {
									tpls := generators.New(embededTemplates)
									gen := generators.GitHub{
										CITagTrigger: c.Bool("tag-release"),
										CIBranch:     "main",
										CIWorkingDir: c.String("workdir"),
										Stage:        c.String("stage"),
										Domain:       c.String("domain"),
										Region:       c.String("region"),
										AppID:        c.String("app_id"),
									}
									if err := gen.RenderFrontend(tpls); err != nil {
										return err
									}

									return nil
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
						Name:  "dev",
						Usage: "Start a service with SSM secrets",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "region", Usage: "AWS region", Required: true},
							&cli.StringFlag{Name: "path", Usage: "SSM Path", Required: true},
						},
						Action: func(c *cli.Context) error {
							input := &awscmd.InputSSMGetParameters{
								Region: c.String("region"),
								Path:   c.String("path"),
							}
							fmt.Fprintf(
								c.App.Writer,
								"%sFetching secrets %s%s\n",
								ctc.ForegroundYellow,
								fmt.Sprint(input.Path, "/*"),
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
