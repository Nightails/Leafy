package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/app_v1"
)

var version = "v0.1.1-dev"

func main() {
	m := app_v1.New(version)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Panicf("Program exited with error: %v", err)
	}
}
