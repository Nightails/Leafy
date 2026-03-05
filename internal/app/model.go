// Package app implements the main loop of the Leafy app.
// This includes user interactions and prompts.
package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	media = iota
	destination
	copying
)

type Model struct {
	state     state
	currStep  step
	mediaList list.Model
	destInput textinput.Model
	err       error
}

func New() Model {
	l := list.New([]list.Item{}, mediaItemDelegate{}, 0, 1)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	ti := textinput.New()
	ti.Placeholder = "Destination Path"

	return Model{
		state:     state{},
		currStep:  media,
		mediaList: l,
		destInput: ti,
		err:       nil,
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
				switch m.currStep {
				case media:
					m.mediaList, cmd = m.mediaList.Update(msg)
					return m, cmd
				case destination:
					m.destInput, cmd = m.destInput.Update(msg)
					return m, cmd
				case copying:
					return m, nil
				}
			}
			return m, nil
		case "q", "ctrl+c":
			// TODO: check for running tasks
			return m, tea.Batch(
				removeDevicesCmd(m.state.devices),
				tea.Quit,
			)
		case " ":
			index := m.mediaList.Index()
			item := m.mediaList.SelectedItem()
			selectItem := item.(mediumItem)
			selectItem.selected = !selectItem.selected
			m.mediaList.SetItem(index, selectItem)
			return m, nil
		case "enter", "return":
			return m, nil
		}
	case errMsg:
		m.err = msg
		return m, nil
	case tea.WindowSizeMsg:
		m.mediaList.SetWidth(msg.Width)
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
		m.mediaList.SetItems(items)
		m.mediaList.SetHeight(len(items))
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

	if m.currStep >= media {
		b.WriteString(m.mediaList.View())
	}
	if m.currStep >= destination {
		b.WriteString("\n")
		b.WriteString(m.destInput.View())
	}
	if m.currStep >= copying {
		b.WriteString("\n")
		// TODO: print the copying progress
	}
	// TODO: print the help bar
	return b.String()
}
