package efs

import (
	"io/fs"
)

func Files(f fs.FS) ([]string, error) {
	files := []string{}

	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	return files, err
}
