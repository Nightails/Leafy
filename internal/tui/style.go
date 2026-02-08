package tui

import "github.com/charmbracelet/lipgloss"

var (
	textStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("255")).
			PaddingLeft(1)
	itemStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("244")).
			PaddingLeft(3)
	selectedStyle = itemStyle.
			Foreground(lipgloss.Color("84")).
			PaddingLeft(1)
	spinnerStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("84")).
			PaddingLeft(1)
	helpStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("244")).
			PaddingLeft(1).
			PaddingBottom(1)
)
