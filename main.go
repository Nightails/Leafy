package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/tui_app"
	"github.com/nightails/leafy/internal/tui_device"
	"github.com/nightails/leafy/internal/tui_media"
)

func main() {
	// Initialize tabs
	var tabs []tea.Model
	deviceTab := tui_device.NewDeviceModel()
	mediaTab := tui_media.NewMediaModel()
	tabs = append(tabs, deviceTab, mediaTab)

	m := tui_app.NewAppModel(tabs)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Panicf("Program exited with error: %v", err)
	}
}
