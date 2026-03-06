package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/app"
)

var version = "dev"

func main() {
	m := app.New(version)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Panicf("Program exited with error: %v", err)
	}
}
