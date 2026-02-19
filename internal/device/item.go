package device

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	style "github.com/nightails/leafy/internal/style"
)

type deviceItem struct {
	usb          USBDevice
	mounting     bool
	spinnerFrame string
}

func (i deviceItem) Title() string {
	return i.usb.Path + " | " + i.usb.Name
}

func (i deviceItem) Description() string {
	if i.mounting {
		return i.spinnerFrame + style.ItemTextStyle.Render("Mounting...")
	}
	if i.usb.Mountpoint == "" {
		return "Not mounted"
	}
	return "↳" + i.usb.Mountpoint
}

func (i deviceItem) FilterValue() string {
	return i.usb.Name
}

type deviceItemDelegate struct{}

func (d deviceItemDelegate) Height() int {
	return 2
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

	fnTitle := style.ItemStyle.Render
	if index == m.Index() {
		fnTitle = func(s ...string) string {
			return style.SelectedTitleStyle.Render(fmt.Sprintf("> %s", strings.Join(s, " ")))
		}
	}

	fnDescription := style.ItemDescriptionStyle.Render
	if index == m.Index() {
		fnDescription = func(s ...string) string {
			return style.SelectedDescriptionStyle.Render(strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprintf(w, "%s\n%s\n", fnTitle(i.Title()), fnDescription(i.Description()))
}
