package app

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/style"
)

type tabID int

const (
	tabDevice tabID = iota
	tabMedia
	tabTransfer
)

type Model struct {
	state  State
	active tabID
	tabs   []tea.Model
}

func NewAppModel(tabs []tea.Model) Model {
	return Model{
		state:  State{},
		active: tabDevice,
		tabs:   tabs,
	}
}

func (m Model) Init() tea.Cmd {
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
		// TODO: move tab handling to individual tab models
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
		if i := slices.Index(m.state.MountPoints, msg.MountPoint); i >= 0 {
			m.state.MountPoints = slices.Delete(m.state.MountPoints, i, i+1)
		}
		return m.broadcastState()
	case FileSelectedMsg:
		// TODO: refactor to use new MediaFile struct
		//m.state.MediaFiles = append(m.state.MediaFiles, MediaFile{Name: msg.Name, Src: msg.Path, Dest: ""})
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
	tTab := style.TabTextStyle.Render("Transfer")

	switch m.active {
	case tabDevice:
		dTab = style.TabSelectedStyle.Render("Device")
	case tabMedia:
		mTab = style.TabSelectedStyle.Render("Media")
	case tabTransfer:
		tTab = style.TabSelectedStyle.Render("Transfer")
	}

	b.WriteString(fmt.Sprintf(" %s | %s | %s\n", dTab, mTab, tTab))
	b.WriteString(m.tabs[m.active].View())

	return b.String()
}

func (m Model) switchTo(id tabID) (tea.Model, tea.Cmd) {
	m.active = id
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
