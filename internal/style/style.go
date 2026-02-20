package style

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).PaddingLeft(1)

	TabSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("84"))
	TabTextStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	ItemTextStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).PaddingLeft(1)
	ItemStyle                = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).PaddingLeft(3)
	ItemDescriptionStyle     = ItemStyle.PaddingLeft(5)
	SelectedItemStyle        = ItemStyle.Foreground(lipgloss.Color("84")).PaddingLeft(1)
	SelectedDescriptionStyle = ItemStyle.Foreground(lipgloss.Color("78")).PaddingLeft(5)

	SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("84")).PaddingLeft(1)
	HelpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).PaddingLeft(1).PaddingBottom(1)
)

const (
	LoadDelay = 1 * time.Second
	QuitDelay = 1 * time.Second
)
