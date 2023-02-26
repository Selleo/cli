package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	colorBlue   = "#00add0"
	colorOrange = "#ee5c1a"
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
	labelStyle = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("#00add0"))

	titleStyle        = lipgloss.NewStyle().MarginLeft(2).MarginTop(1)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#00add0"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)
