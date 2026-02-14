package tui_app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/tui_device"
	"github.com/nightails/leafy/internal/tui_media"
)

type tabID int

const (
	tabDevice tabID = iota
	tabMedia
)

type AppModel struct {
	active tabID
	tabs   []tea.Model
	inited map[tabID]bool
}

func NewAppModel() AppModel {
	dev := tui_device.NewDeviceModel()
	media := tui_media.NewMediaModel()
	return AppModel{
		active: tabDevice,
		tabs:   []tea.Model{dev, media},
		inited: map[tabID]bool{},
	}
}

func (m AppModel) Init() tea.Cmd {
	m.inited[m.active] = true
	return m.tabs[m.active].Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "tab":
			return m.switchTo(tabID((int(m.active) + 1) % len(m.tabs)))
		case "shift+tab":
			return m.switchTo(tabID((int(m.active) - 1 + len(m.tabs)) % len(m.tabs)))
		}
	}

	if _, ok := msg.(tea.WindowSizeMsg); ok {
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
