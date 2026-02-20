package device

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/style"
)

type deviceItem struct {
	usb          USBDevice
	mounting     bool
	spinnerFrame string
}

func (i deviceItem) cursor() string {
	if i.mounting {
		return i.spinnerFrame
	}
	return style.SelectedItemStyle.Render(">")
}
func (i deviceItem) namePath() string {
	return i.usb.Path + " | " + i.usb.Name
}

func (i deviceItem) mountedPath() string {
	if i.mounting {
		return " -> Mounting..."
	}
	if i.usb.Mountpoint == "" {
		return " -> Not mounted"
	}
	return " -> " + i.usb.Mountpoint
}

func (i deviceItem) FilterValue() string {
	return i.usb.Name
}

type deviceItemDelegate struct{}

func (d deviceItemDelegate) Height() int {
	return 1
}

func (d deviceItemDelegate) Spacing() int {
	return 0
}

func (d deviceItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d deviceItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var i, ok = item.(deviceItem)
	if !ok {
		return
	}

	fn := style.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return i.cursor() + style.SelectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprintf(w, fn(i.namePath()+i.mountedPath()))
}
