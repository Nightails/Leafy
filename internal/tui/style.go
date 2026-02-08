package tui

import "github.com/charmbracelet/lipgloss"

var (
	itemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).PaddingLeft(3)
	selectedStyle = itemStyle.Foreground(lipgloss.Color("84")).PaddingLeft(1)
	spinnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("84"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).PaddingLeft(1).PaddingBottom(1)
)
