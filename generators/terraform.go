package generators

import (
	"fmt"
	"os"
	"path/filepath"
)

type Terraform struct {
	TerraformCloudOrganization string
	TerraformCloudWorkspace    string
	Region                     string
	Namespace                  string
	Stage                      string
	Name                       string
	Subdomain                  string
	Domain                     string
	IAMCI                      string
	IAMApp                     string
	ECSInstanceType            string
	ECSMinSize                 int
	ECSServiceMinMemory        int
	ECSServiceCpu              int
	ECSServiceMaxMemory        int
	ECSServicePort             int
	ECSOneOffs                 []string
	LBName                     string
	DBIdentifier               string
	DBName                     string
	DBUser                     string
	DBMultiAZ                  bool
	DBApplyImmediately         bool
	BucketName                 string
}

func (tf *Terraform) Render(t *TemplateRenderer) error {
	err := os.MkdirAll(filepath.Join("terraform", tf.Stage), 0755)
	if err != nil {
		return err
	}
	return errors(
		tf.write(t, "acm.tf"),
		tf.write(t, "ci.tf"),
		tf.write(t, "ecr.tf"),
		tf.write(t, "ecs.tf"),
		tf.write(t, "iam.tf"),
		tf.write(t, "load_balancer.tf"),
		tf.write(t, "main.tf"),
		tf.write(t, "rds.tf"),
		tf.write(t, "route53.tf"),
		tf.write(t, "s3.tf"),
		tf.write(t, "variables.tf"),
		tf.write(t, "versions.tf"),
		tf.write(t, "vpc.tf"),
	)
}

func (tf *Terraform) write(t *TemplateRenderer, name string) error {
	path := filepath.Join("terraform", tf.Stage, name)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_ = f.Chmod(0755)
	defer f.Close()
	return t.Render(f, fmt.Sprint("templates/terraform/", name), tf)
}
