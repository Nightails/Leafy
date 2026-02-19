package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type StateMsg struct {
	State State
}
type DeviceMountedMsg struct {
	MountPoint string
}
type FileSelectedMsg struct {
	Path string
}

type ErrMsg error
type FinishedMsg struct{}
type QuitNowMsg struct{}

// AfterCmd returns a command that sends the given message after the given duration
func AfterCmd(d time.Duration, msg tea.Msg) tea.Cmd {
	if d <= 0 {
		return func() tea.Msg { return msg }
	}
	return tea.Tick(d, func(time.Time) tea.Msg { return msg })
}
