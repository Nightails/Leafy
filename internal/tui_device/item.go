package tui_device

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/device"
	style "github.com/nightails/leafy/internal/tui_style"
)

type deviceItem struct {
	usb device.USBDevice
}

func (i deviceItem) Title() string {
	return i.usb.Path + " | " + i.usb.Name
}

func (i deviceItem) Description() string {
	return "↳" + i.usb.Mountpoint
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
	i, ok := item.(deviceItem)
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

	if i.usb.Mountpoint == "" {
		_, _ = fmt.Fprintf(w, "%s\n", fnTitle(i.Title()))
	} else {
		_, _ = fmt.Fprintf(w, "%s\n%s\n", fnTitle(i.Title()), fnDescription(i.Description()))
	}
}
