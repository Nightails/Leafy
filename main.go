package main

import tea "github.com/charmbracelet/bubbletea"

func main() {
	m := model{}
	p := tea.NewProgram(m)
	_, _ = p.Run()
}

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "Hello World!"
}
