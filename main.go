package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/app"
	"github.com/nightails/leafy/internal/device"
	"github.com/nightails/leafy/internal/media"
	"github.com/nightails/leafy/internal/transfer"
)

func main() {
	// Initialize tabs
	deviceTab := device.NewModel()
	mediaTab := media.NewModel()
	transferTab := transfer.NewModel()

	// Start the TUI app
	m := app.NewAppModel([]tea.Model{deviceTab, mediaTab, transferTab})
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Panicf("Program exited with error: %v", err)
	}
}
