package tui

import (
	"fmt"
	"leafy/internal/usb"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type USBModel struct {
	Devices []usb.BlockDevice
	Error   string
}

func (m USBModel) Init() tea.Cmd {
	return nil
}

func (m USBModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m USBModel) View() string {
	var b strings.Builder
	b.WriteString("Welcome to Leafy!\nA simple TUI app for importing media to your computer.\n\n")

	if m.Error != "" {
		b.WriteString(fmt.Sprintf("Error: %s\n", m.Error))
		return b.String()
	}

	if len(m.Devices) == 0 {
		b.WriteString("No USB devices found\n")
		return b.String()
	}

	b.WriteString("USB Devices:\n")
	for _, d := range m.Devices {
		if len(d.Children) > 0 {
			for _, c := range d.Children {
				var name string
				if c.Name == "" {
					name = c.Label
				} else {
					name = c.Name
				}
				b.WriteString(fmt.Sprintf(" - %s (%s)\n", name, c.Path))
			}
		} else {
			var name string
			if d.Name == "" {
				name = d.Label
			} else {
				name = d.Name
			}
			b.WriteString(fmt.Sprintf(" - %s (%s)\n", name, d.Path))
		}
	}

	return b.String()
}
