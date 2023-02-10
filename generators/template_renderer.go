package generators

import (
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	"github.com/Selleo/cli/efs"
)

type TemplateRenderer struct {
	tpls map[string]*template.Template
}

func New(f embed.FS) *TemplateRenderer {
	tpls := map[string]*template.Template{}

	files, err := efs.Files(f)
	if err != nil {
		panic(fmt.Errorf("Failed to init templates: %v", err))
	}
	for _, file := range files {
		t, err := template.New(filepath.Base(file)).Delims("{{{", "}}}").ParseFS(f, file)
		if err != nil {
			panic(err)
		}
		tpls[file] = t
	}
	return &TemplateRenderer{
		tpls: tpls,
	}
}

func (r *TemplateRenderer) Render(w io.Writer, name string, data any) error {
	if tpl, ok := r.tpls[name]; ok {
		return tpl.Execute(w, data)
	}
	return fmt.Errorf("Template %s not found", name)
}
