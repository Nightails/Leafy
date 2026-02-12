package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/tui_device"
)

func main() {
	m := tui_device.NewDeviceModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Panicf("Program exited with error: %v", err)
	}
}
