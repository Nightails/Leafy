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
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return Model{
		state: state{},
		list:  l,
		err:   nil,
	}
}

func (m Model) Init() tea.Cmd {
	return initDevicesCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		default:
			if m.err == nil {
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}
			return m, nil
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
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
	case devicesMsg:
		if len(msg) == 0 {
			return m, nil
		}
		m.state.devices = msg
		return m, findMediaCmd(msg)
	case mediaMsg:
		if len(msg) == 0 {
			return m, nil
		}
		items := make([]list.Item, 0, len(msg))
		for _, i := range msg {
			items = append(items, mediumItem{i, false})
		}
		m.list.SetItems(items)
		m.list.SetHeight(len(items))
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

	b.WriteString(m.list.View())
	return b.String()
}
