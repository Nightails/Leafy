package tui_style

import "github.com/charmbracelet/lipgloss"

var (
	TextStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("255")).
			PaddingLeft(1)
	ItemTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	ItemStyle      = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("244")).
			PaddingLeft(3)
	ItemDescriptionStyle = ItemStyle.PaddingLeft(5)
	SelectedTitleStyle   = ItemStyle.
				Foreground(lipgloss.Color("84")).
				PaddingLeft(1)
	SelectedDescriptionStyle = ItemStyle.
					Foreground(lipgloss.Color("78")).
					PaddingLeft(5)
	SpinnerStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("84")).
			PaddingLeft(1)
	HelpStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("244")).
			PaddingLeft(1).
			PaddingBottom(1)
)
