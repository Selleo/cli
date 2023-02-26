package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type awsSecretItem struct {
	Key   string
	Value string
}

func newIndexSecretsModel(data map[string]string) modelSecretsIndex {
	items := []list.Item{}
	for k, v := range data {
		items = append(items, awsSecretItem{k, v})
	}

	const defaultWidth = 20
	const listHeight = 20

	l := list.New(items, awsSecretItemDelegate{}, defaultWidth, listHeight)
	l.SetFilteringEnabled(true)
	l.Title = "Choose secret to edit or add new one:"
	l.SetShowStatusBar(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			// key.NewBinding(
			// 	key.WithKeys("a"),
			// 	key.WithHelp("a", "add secret"),
			// ),
		}
	}

	m := modelSecretsIndex{
		list:    l,
		secrets: data,
	}
	return m
}

func (i awsSecretItem) FilterValue() string { return i.Key }

type awsSecretItemDelegate struct{}

func (d awsSecretItemDelegate) Height() int                               { return 1 }
func (d awsSecretItemDelegate) Spacing() int                              { return 0 }
func (d awsSecretItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d awsSecretItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(awsSecretItem)
	if !ok {
		panic("must be awsSecretItem")
	}

	str := i.Key

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type modelSecretsIndex struct {
	list    list.Model
	secrets map[string]string
}

func (m modelSecretsIndex) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m modelSecretsIndex) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch keypress := msg.String(); keypress {
		// case "a":
		// 	edit := newEditSecretModel("New secret", "", "")
		// 	return edit, edit.Init()

		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(awsSecretItem)
			if !ok {
				panic("must be awsSecretItem")
			}
			secretKey := string(i.Key)
			secretValue := m.secrets[secretKey]
			edit := newEditSecretModel("Edit secret", secretKey, secretValue)
			return edit, edit.Init()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m modelSecretsIndex) View() string {
	return m.list.View()
}
