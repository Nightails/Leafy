package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/Nightails/Leafy/internal/usb"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type usbDeviceItem struct {
	dev usb.BlockDevice
}

func (i usbDeviceItem) Title() string {
	if i.dev.Label != "" {
		return fmt.Sprintf("%s (%s)", i.dev.Name, i.dev.Label)
	}
	return i.dev.Name
}

func (i usbDeviceItem) Description() string {
	return fmt.Sprintf("%s | %s | %s", i.dev.Path, i.dev.Tran, i.dev.Type)
}

func (i usbDeviceItem) FilterValue() string {
	return strings.Join([]string{i.dev.Name, i.dev.Path, i.dev.Label, i.dev.Tran, i.dev.Type}, " ")
}

type usbItemDelegate struct{}

func (d usbItemDelegate) Height() int                             { return 1 }
func (d usbItemDelegate) Spacing() int                            { return 0 }
func (d usbItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d usbItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(usbDeviceItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s: %s", i.Title(), i.Description())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprintf(w, fn(str))
}
