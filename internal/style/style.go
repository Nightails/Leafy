package style

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	HelpTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

const (
	LoadDelay = 1 * time.Second
	QuitDelay = 1 * time.Second
)
