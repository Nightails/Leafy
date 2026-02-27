// Package app implements the main loop of the Leafy app.
// This includes user interactions and prompts.
package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	state state
}

func New() Model {
	return Model{
		state: state{},
	}
}

func (m Model) Init() tea.Cmd {
	// TODO: 1.scanning/mouting usb devices
	// TODO: 2.scanning for supported media files
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: 1.handle user input/navigation
	// TODO: 2.handle msgs for mounting devices/scanning media files/transfering media files
	return m, nil
}

func (m Model) View() string {
	// TODO: 1.display found media files
	// TODO: 2.display prompt for destination path
	// TODO: 3.display transfering progress
	// TODO: 4.display help bar
	return "leafy"
}
