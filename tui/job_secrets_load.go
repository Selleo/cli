package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type LoadSecretsJob struct {
	Region string
	Path   string
}

func (j LoadSecretsJob) Run() tea.Msg {
	time.Sleep(200 * time.Millisecond)

	model := newIndexSecretsModel(
		map[string]string{
			"AWS_ACCESS_KEY_ID":     "secret",
			"AWS_SECRET_ACCESS_KEY": "secret",
		},
	)

	return JobMessage{
		Err:   nil,
		NextModel: model,
	}
}
