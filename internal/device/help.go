package device

import (
	"strings"

	style "github.com/nightails/leafy/internal/style"
)

var controls = []string{
	"[↑/k] up",
	"[↓/j] down",
	"[space] mount/unmount",
	"[backspace] return",
	"[enter] continue",
	"[s] re-scan",
	"[q] quit",
}

func helpBarView() string {
	var b strings.Builder
	b.WriteString("Please select and mount usb devices, then continue.\n")
	b.WriteString(strings.Join(controls, " | "))
	return style.HelpStyle.Render(b.String())
}
