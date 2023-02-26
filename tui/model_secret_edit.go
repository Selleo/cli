package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type modelSecretEdit struct {
	title       string
	secretInput textarea.Model
	keyInput    textinput.Model
	focused     int
	oldKey      string
	oldValue    string
	inputsCount int
}

func newEditSecretModel(title string, key string, value string) modelSecretEdit {
	keyInput := textinput.New()
	keyInput.Placeholder = "Paste your key here"
	keyInput.Focus()
	keyInput.Width = 100
	keyInput.PromptStyle = lipgloss.NewStyle().MarginLeft(2)
	keyInput.Prompt = ""
	keyInput.SetValue(key)

	secretInput := textarea.New()
	secretInput.Placeholder = "Paste your secret here..."
	secretInput.CharLimit = 4096
	secretInput.Prompt = ""
	secretInput.SetHeight(20)
	secretInput.SetValue(value)

	return modelSecretEdit{
		title:       title,
		inputsCount: 2,
		secretInput: secretInput,
		keyInput:    keyInput,
		focused:     0,
		oldKey:      key,
		oldValue:    value,
	}
}

func (m modelSecretEdit) Init() tea.Cmd {
	return textarea.Blink
}

func (m modelSecretEdit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.secretInput.SetWidth(msg.Width - h)
		m.secretInput.SetHeight(msg.Height - v - 2)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			job := &SaveSecretJob{
				OldKey:    m.oldKey,
				OldSecret: m.oldValue,
				NewKey:    m.keyInput.Value(),
				NewSecret: m.secretInput.Value(),
			}
			model := newJobModel(fmt.Sprintf("Saving secret %s", job.NewKey), job)
			return model, model.Init()

		case tea.KeyEsc:
			if m.secretInput.Focused() {
				m.secretInput.Blur()
			}
			if m.keyInput.Focused() {
				m.keyInput.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}

		// default:
		// 	if !m.secretInput.Focused() {
		// 		cmd = m.secretInput.Focus()
		// 		cmds = append(cmds, cmd)
		// 	}
		m.keyInput.Blur()
		m.secretInput.Blur()

		if m.focused == 0 {
			m.keyInput.Focus()
		}
		if m.focused == 1 {
			m.secretInput.Focus()
		}
	}

	m.keyInput, cmd = m.keyInput.Update(msg)
	cmds = append(cmds, cmd)
	m.secretInput, cmd = m.secretInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m modelSecretEdit) View() string {
	return fmt.Sprintf(
		"  %s:\n\n%s\n%s\n\n%s\n%s\n",
		m.title,
		labelStyle.Width(30).Render("Key"),
		m.keyInput.View(),
		labelStyle.Width(30).Render("Secret"),
		m.secretInput.View(),
	)
}

func (m *modelSecretEdit) nextInput() {
	m.focused = (m.focused + 1) % m.inputsCount
}

func (m *modelSecretEdit) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = m.inputsCount - 1
	}
}
