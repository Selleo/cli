package generators

import "os"

type Docker struct {
	CmdServer string
	OneOffs   map[string]string
}

func (d *Docker) Render(t *TemplateRenderer) error {
	f, err := os.Create("entrypoint.sh")
	if err != nil {
		return err
	}
	_ = f.Chmod(0755)
	defer f.Close()

	err = t.Render(f, "templates/docker/entrypoint.sh", d)
	return err
}
