package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type Job interface {
	Run() tea.Msg
}

type JobMessage struct {
	NextModel tea.Model // job should trigger new model
	Err       error     // job returned an error
}

type modelJob struct {
	text    string
	spinner spinner.Model
	err     error
	job     Job
}

func newJobModel(text string, job Job) modelJob {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlue))

	return modelJob{
		text:    text,
		spinner: s,
		job:     job,
	}
}

func (m modelJob) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.job.Run,
	)
}

func (m modelJob) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}

	case JobMessage:
		m.err = msg.Err
		if msg.NextModel != nil {
			return msg.NextModel, msg.NextModel.Init()
		}
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m modelJob) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n %s %s\n\n", m.spinner.View(), m.text)

	return str
}
