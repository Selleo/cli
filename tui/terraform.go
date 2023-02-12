package tui

import (
	"context"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

// #e8e8e8 	#00add0 	#ee5c1a 	---
// Gray 	#5b6266 	#84bf40 	#8a1dcf 	#edd81a
// Dark 	#293333 	#ee1a1a 	#663c12 	---
	// titleStyle = lipgloss.NewStyle().
	// 		Foreground(lipgloss.Color("#FFFDF5")).
	// 		Background(lipgloss.Color("#25A065")).
	// 		Padding(0, 1)
	//
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type Terraform struct {
}

func (t Terraform) Run(ctx context.Context) error {
	items := []list.Item{
		awsRegionItem{"us-east-1", "US (N. Virginia)", "🇺🇸"},
		awsRegionItem{"us-east-2", "US (Ohio)", "🇺🇸"},
		awsRegionItem{"us-west-1", "US (N. California)", "🇺🇸"},
		awsRegionItem{"us-west-2", "US (Oregon)", "🇺🇸"},
		awsRegionItem{"ca-central-1", "Canada (Central)", "🇨🇦"},

		awsRegionItem{"eu-central-1", "Europe (Frankfurt)", "🇩🇪"},
		awsRegionItem{"eu-central-2", "Europe (Zurich)", "🇨🇭"},
		awsRegionItem{"eu-west-1", "Europe (Ireland)", "🇮🇪"},
		awsRegionItem{"eu-west-2", "Europe (London)", "🇬🇧"},
		awsRegionItem{"eu-west-3", "Europe (Paris)", "🇫🇷"},
		awsRegionItem{"eu-south-1", "Europe (Milan)", "🇮🇹"},
		awsRegionItem{"eu-south-2", "Europe (Spain)", "🇪🇸"},
		awsRegionItem{"eu-north-1", "Europe (Stockholm)", "🇸🇪"},

		awsRegionItem{"ap-south-1", "Asia Pacific (Mumbai)", "🇮🇳"},
		awsRegionItem{"ap-south-2", "Asia Pacific (Hyderabad)", "🇮🇳"},
		awsRegionItem{"ap-southeast-1", "Asia Pacific (Singapore)", "🇸🇬"},
		awsRegionItem{"ap-southeast-2", "Asia Pacific (Sydney)", "🇦🇺"},
		awsRegionItem{"ap-southeast-3", "Asia Pacific (Jakarta)", "🇮🇩"},
		awsRegionItem{"ap-southeast-5", "Asia Pacific (Melbourne)", "🇦🇺"},
		awsRegionItem{"ap-northeast-1", "Asia Pacific (Tokio)", "🇯🇵"},
		awsRegionItem{"ap-northeast-2", "Asia Pacific (Seoul)", "🇰🇷"},
		awsRegionItem{"ap-northeast-3", "Asia Pacific (Osaka)", "🇯🇵"},

		awsRegionItem{"me-central-1", "Middle East (United Arab Emirates)", "🇦🇪"},
		awsRegionItem{"me-south-1", "Middle East (Bahrain)", "🇧🇭"},
		awsRegionItem{"sa-east-1", "South America (São Paulo)", "🇧🇷"},
		awsRegionItem{"af-south-1", "Africa (Cape Town)", "🇿🇦"},
		awsRegionItem{"cn-northwest-1", "Mainland China (Ningxia)", "🇨🇳"},
		awsRegionItem{"cn-north-1", "Mainland China (Beijing)", "🇨🇳"},
		awsRegionItem{"us-gov-east-1", "AWS GovCloud (US-East)", "🇺🇸"},
		awsRegionItem{"us-gov-west-1", "AWS GovCloud (US-West)", "🇺🇸"},
	}

	const defaultWidth = 20

	l := list.New(items, awsRegionItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Choose AWS Region for deployment."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	_, err := tea.NewProgram(m).Run()
	return err
}

const listHeight = 20

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#00add0"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type awsRegionItem struct {
	ID string
	Name string
	Flag string
}

func (i awsRegionItem) FilterValue() string { return i.ID }

type awsRegionItemDelegate struct{}

func (d awsRegionItemDelegate) Height() int                               { return 1 }
func (d awsRegionItemDelegate) Spacing() int                              { return 0 }
func (d awsRegionItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d awsRegionItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(awsRegionItem)
	if !ok {
		panic("must be awsRegionItem")
	}

	str := fmt.Sprintf("%2s %15s   %s", i.Flag, i.ID, i.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(awsRegionItem)
			if ok {
				m.choice = string(i.ID)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	return "\n" + m.list.View()
}
