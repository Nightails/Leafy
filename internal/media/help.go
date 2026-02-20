package media

import (
	"strings"

	"github.com/nightails/leafy/internal/style"
)

var controls = []string{
	"[↑/k] up",
	"[↓/j] down",
	"[space] select/de-select",
	"[backspace] return",
	"[enter] continue",
	"[s] re-scan",
	"[q] quit",
}

func helpBarView() string {
	var b strings.Builder
	b.WriteString("Please select media files, then continue.\n")
	b.WriteString(strings.Join(controls, " | "))
	return style.HelpStyle.Render(b.String())
}
