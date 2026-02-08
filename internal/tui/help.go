package tui

import "strings"

var controls = []string{
	"[↑/k] up",
	"[↓/j] down",
	"[enter] mount",
	"[tab] continue",
	"[s] scan",
	"[ctrl+c] quit",
}

func helpBarView() string {
	help := strings.Join(controls, " | ")
	return helpStyle.Render(help)
}
