package main

import (
	"fmt"
	"leafy/internal/tui"
	"leafy/internal/usb"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := tui.USBModel{}
	devs, err := usb.FindUSBDevices()
	if err != nil {
		m.Error = fmt.Sprintf("Failed to find USB devices: %v", err)
	}
	m.Devices = devs

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Panicf("Program exited with error: %v", err)
	}
}
