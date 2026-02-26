package transfer

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
	// TODO: mediaFiles that contain src and dest path
	transferCh chan string
}

func NewModel() Model {
	return Model{
		transferCh: make(chan string),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "Transferring"
}
