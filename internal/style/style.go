package style

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	HeaderStyle = lipgloss.NewStyle().
			Border(lipgloss.ASCIIBorder()).
			Foreground(lipgloss.Color("255")).
			Bold(true)
	TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	HelpTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

const (
	LoadDelay = 1 * time.Second
	QuitDelay = 1 * time.Second
)
