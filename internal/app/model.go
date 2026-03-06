// Package app implements the main loop of the Leafy app.
// This includes user interactions and prompts.
package app

import (
	"path/filepath"
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
	l.DisableQuitKeybindings()

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
		case "ctrl+c":
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
			switch m.currStep {
			case media:
				for _, i := range m.mediaList.Items() {
					mi := i.(mediumItem)
					if mi.selected {
						m.state.media = append(m.state.media, mi.medium)
					}
				}
				if m.state.media == nil {
					return m, nil
				}
				m.currStep++
				m.destInput.Focus()
				return m, nil
			case destination:
				dest := m.destInput.Value()
				if dest == "" {
					return m, nil
				}
				// TODO: apply naming scheme with increment
				for i := range m.state.media {
					m.state.media[i].dest = filepath.Join(dest, m.state.media[i].name, m.state.media[i].format)
				}
				m.currStep++
				m.destInput.Blur()
				// TODO: start copy tasks
				return m, nil
			case copying:
				// Do nothing while copying
				// Only interaction is quitting/cancelling
				return m, nil
			}
			return m, nil
		}
	case errMsg:
		m.err = msg
		return m, nil
	case tea.WindowSizeMsg:
		m.mediaList.SetWidth(msg.Width)
		return m, nil
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
		b.WriteString("Found media:\n")
		b.WriteString(m.mediaList.View())
	}
	if m.currStep >= destination {
		b.WriteString("\n\n")
		b.WriteString("Destination Path:\n")
		b.WriteString(m.destInput.View())
	}
	if m.currStep >= copying {
		b.WriteString("\n\n")
		b.WriteString("Copying Progress:\n")
		// TODO: print the copying progress
	}
	// TODO: print the help bar
	return b.String()
}
