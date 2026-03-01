// Package app implements the main loop of the Leafy app.
// This includes user interactions and prompts.
package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	state state
	list  list.Model
	err   error
}

func New() Model {
	l := list.New([]list.Item{}, mediaItemDelegate{}, 100, 10)
	l.SetShowTitle(false)
	l.SetShowStatusBar(true)
	l.SetShowPagination(true)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return Model{
		state: state{},
		list:  l,
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
		// TODO: handle msg for transferring media
	}
	return m, nil
}

func (m Model) View() string {
	// TODO: 1.display found file files
	// TODO: 2.display prompt for destination path
	// TODO: 3.display transferring progress
	// TODO: 4.display help bar
	var b strings.Builder
	if m.err != nil {
		b.WriteString(m.err.Error())
		return b.String()
	}

	if len(m.state.media) == 0 {
		b.WriteString("No file files found")
		return b.String()
	}

	for _, m := range m.state.media {
		b.WriteString("\n" + m.name)
	}
	return b.String()
}
