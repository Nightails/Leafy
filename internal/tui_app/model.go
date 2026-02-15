package tui_app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type tabID int

const (
	tabDevice tabID = iota
	tabMedia
)

type AppModel struct {
	state  AppState
	active tabID
	tabs   []tea.Model
	inited map[tabID]bool
}

func NewAppModel(tabs []tea.Model) AppModel {
	return AppModel{
		state:  AppState{},
		active: tabDevice,
		tabs:   tabs,
		inited: map[tabID]bool{},
	}
}

func (m AppModel) Init() tea.Cmd {
	m.inited[m.active] = true
	return m.tabs[m.active].Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmds []tea.Cmd
		for i := range m.tabs {
			var cmd tea.Cmd
			m.tabs[i], cmd = m.tabs[i].Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			return m.switchTo(tabID((int(m.active) + 1) % len(m.tabs)))
		case "shift+tab":
			return m.switchTo(tabID((int(m.active) - 1 + len(m.tabs)) % len(m.tabs)))
		}
	case DeviceMountedMsg:
		m.state.MountPoints = append(m.state.MountPoints, msg.MountPoint)
		return m.broadcastState()
	case FileSelectedMsg:
		m.state.FilesToTransfer = append(m.state.FilesToTransfer, msg.Path)
		return m.broadcastState()
	}

	var cmd tea.Cmd
	m.tabs[m.active], cmd = m.tabs[m.active].Update(msg)
	return m, cmd
}

func (m AppModel) View() string {
	return m.tabs[m.active].View()
}

func (m AppModel) switchTo(id tabID) (tea.Model, tea.Cmd) {
	m.active = id
	if m.inited[id] {
		return m, nil
	}
	m.inited[id] = true
	return m, m.tabs[id].Init()
}

func (m AppModel) broadcastState() (tea.Model, tea.Cmd) {
	msg := AppStateMsg{m.state}

	var cmds []tea.Cmd
	for i := range m.tabs {
		var cmd tea.Cmd
		m.tabs[i], cmd = m.tabs[i].Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}
