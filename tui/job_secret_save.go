package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type SaveSecretJob struct {
	OldKey    string
	OldSecret string
	NewKey    string
	NewSecret string
}

func (j SaveSecretJob) Run() tea.Msg {
	time.Sleep(1 * time.Second)

	model := newIndexSecretsModel(
		map[string]string{
			"AWS_ACCESS_KEY_ID":     "secret",
			"AWS_SECRET_ACCESS_KEY": "secret",
		},
	)

	return JobMessage{
		Err:       nil,
		NextModel: model,
	}
}
