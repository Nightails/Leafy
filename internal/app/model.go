// Package app implements the main loop of the Leafy app.
// This includes user interactions and prompts.
package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	state state
	err   error
}

func New() Model {
	return Model{
		state: state{},
	}
}

func (m Model) Init() tea.Cmd {
	return initDevicesCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// TODO: check for running tasks
			return m, tea.Batch(
				removeDevicesCmd(m.state.devices),
				tea.Quit,
			)
		}
	case errMsg:
		m.err = msg
		return m, nil
	case devicesMsg:
		if len(msg) == 0 {
			return m, nil
		}
		m.state.devices = msg
		return m, findMediaCmd(msg)
	case mediaMsg:
		m.state.media = msg
		return m, nil
	}
	// TODO: 2.handle msgs for mounting devices/scanning media files/transfering media files
	return m, nil
}

func (m Model) View() string {
	// TODO: 1.display found media files
	// TODO: 2.display prompt for destination path
	// TODO: 3.display transferring progress
	// TODO: 4.display help bar
	var b strings.Builder
	if m.err != nil {
		b.WriteString(m.err.Error())
		return b.String()
	}

	if len(m.state.media) == 0 {
		b.WriteString("No media files found")
		return b.String()
	}

	for _, m := range m.state.media {
		b.WriteString("\n" + m.src)
	}
	return b.String()
}
