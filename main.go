package main

import (
	"leafy/internal/tui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := tui.NewUSBModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Panicf("Program exited with error: %v", err)
	}
}
