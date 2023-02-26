package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

type Program struct {
}

func (p *Program) Run(ctx context.Context, region string, path string) error {
	job := &LoadSecretsJob{
		Region: region,
		Path:   path,
	}
	model := newJobModel("Loading secrets", job)
	prog := tea.NewProgram(model)
	_, err := prog.Run()
	return err
}
