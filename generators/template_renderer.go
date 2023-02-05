package generators

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"text/template"
)

type TemplateRenderer struct {
	tpls map[string]*template.Template
}

func New(efs embed.FS) *TemplateRenderer {
	tpls := map[string]*template.Template{}

	files := []string{}
	err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("Failed to init templates: %v", err))
	}
	for _, file := range files {
		t, err := template.New(filepath.Base(file)).Delims("{{{", "}}}").ParseFS(efs, file)
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
