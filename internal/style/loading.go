package style

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
)

func NewLineSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = SpinnerStyle
	return s
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
