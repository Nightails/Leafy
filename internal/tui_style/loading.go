package tui_style

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func NewLineSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = SpinnerStyle
	return s
}

// MsgAfter delays the given message by the given duration
func MsgAfter(d time.Duration, msg tea.Msg) tea.Msg {
	if d <= 0 {
		return func() tea.Msg { return msg }
	}
	return tea.Tick(d, func(time.Time) tea.Msg { return msg })
}

type MinDuration struct {
	Min       time.Duration
	startedAt time.Time
}

func (m *MinDuration) StartNow() {
	m.startedAt = time.Now()
}

func (m *MinDuration) EnsureStarted() {
	if m.startedAt.IsZero() {
		m.StartNow()
	}
}

func (m *MinDuration) Remaining() time.Duration {
	m.EnsureStarted()
	return m.Min - time.Since(m.startedAt)
}
