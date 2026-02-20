package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/style"
)

type State struct {
	MountPoints     []string
	FilesToTransfer []string
}

type tabID int

const (
	tabDevice tabID = iota
	tabMedia
)

type Model struct {
	state  State
	active tabID
	tabs   []tea.Model
	inited map[tabID]bool
}

func NewAppModel(tabs []tea.Model) Model {
	return Model{
		state:  State{},
		active: tabDevice,
		tabs:   tabs,
		inited: map[tabID]bool{},
	}
}

func (m Model) Init() tea.Cmd {
	m.inited[m.active] = true
	return m.tabs[m.active].Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "enter", "return":
			i := min(int(m.active)+1, len(m.tabs)-1)
			return m.switchTo(tabID(i))
		case "backspace":
			i := max(0, int(m.active)-1)
			return m.switchTo(tabID(i))
		}
	case DeviceMountedMsg:
		m.state.MountPoints = append(m.state.MountPoints, msg.MountPoint)
		return m.broadcastState()
	case DeviceUnmountedMsg:
		for i, mp := range m.state.MountPoints {
			if mp == msg.MountPoint {
				m.state.MountPoints = append(m.state.MountPoints[:i], m.state.MountPoints[i+1:]...)
				break
			}
		}
		return m.broadcastState()
	case FileSelectedMsg:
		m.state.FilesToTransfer = append(m.state.FilesToTransfer, msg.Path)
		return m.broadcastState()
	}

	var cmd tea.Cmd
	m.tabs[m.active], cmd = m.tabs[m.active].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var b strings.Builder

	dTab := style.TabTextStyle.Render("Device")
	mTab := style.TabTextStyle.Render("Media")

	switch m.active {
	case tabDevice:
		dTab = style.TabSelectedStyle.Render("Device")
	case tabMedia:
		mTab = style.TabSelectedStyle.Render("Media")
	}

	b.WriteString(fmt.Sprintf(" %s | %s\n", dTab, mTab))
	b.WriteString(m.tabs[m.active].View())

	return b.String()
}

func (m Model) switchTo(id tabID) (tea.Model, tea.Cmd) {
	m.active = id
	if m.inited[id] {
		return m, nil
	}
	m.inited[id] = true
	return m, m.tabs[id].Init()
}

func (m Model) broadcastState() (tea.Model, tea.Cmd) {
	msg := StateMsg{m.state}

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
