package generators

import (
	"fmt"
	"os"
)

type GitHub struct {
	CITagTrigger bool
	CIBranch     string
	CIWorkingDir string
	Stage        string
	Domain       string
	Subdomain    string
	Region       string
	AppID        string
	ECSCluster   string
	ECSService   string
	ECSOneOffs   []string
}

func (g *GitHub) RenderFrontend(t *TemplateRenderer) error {
	err := os.MkdirAll(".github/workflows", 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(fmt.Sprintf(".github/workflows/deploy-%s-frontend.yml", g.Stage))
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Render(f, "templates/github/workflows/deploy-frontend.yml", g)
	return err
}

func (g *GitHub) RenderBackend(t *TemplateRenderer) error {
	err := os.MkdirAll(".github/workflows", 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(fmt.Sprintf(".github/workflows/deploy-%s-backend.yml", g.Stage))
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Render(f, "templates/github/workflows/deploy-backend.yml", g)
	return err
}
