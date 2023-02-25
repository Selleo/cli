package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type editSecretModel struct {
	secretInput textarea.Model
	keyInput    textinput.Model
	focused     int
	key         string
	value       string
	inputsCount int
}

func newEditSecretModel(key string, value string) editSecretModel {
	ta := textarea.New()
	ta.Placeholder = "Paste your secret here..."
	ta.CharLimit = 4096
	// ta.SetWidth(100)
	ta.Prompt = ""
	ta.SetHeight(20)

	keyInput := textinput.New()
	keyInput.Placeholder = "Paste your key here"
	keyInput.Focus()
	keyInput.Width = 100
	keyInput.PromptStyle = lipgloss.NewStyle().MarginLeft(2)
	keyInput.Prompt = ""
	// keyInput.Validate = ccnValidator

	return editSecretModel{
		inputsCount: 2,
		secretInput: ta,
		keyInput:    keyInput,
		focused:     0,
		key:         key,
		value:       value,
	}
}

func Background(color lipgloss.Color) {
	panic("unimplemented")
}

func (m editSecretModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m editSecretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// if m.focused == 0 {
			// 	m.nextInput()
			// }
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

func (m editSecretModel) View() string {
	return fmt.Sprintf(
		"Edit secret:\n\n%s\n%s\n\n%s\n%s\n",
		labelStyle.Width(30).Render("Key"),
		m.keyInput.View(),
		labelStyle.Width(30).Render("Secret"),
		m.secretInput.View(),
	)
}

func (m *editSecretModel) nextInput() {
	m.focused = (m.focused + 1) % m.inputsCount
}

func (m *editSecretModel) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = m.inputsCount - 1
	}
}
