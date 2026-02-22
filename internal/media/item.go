package media

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	style "github.com/nightails/leafy/internal/style"
)

type mediaItem struct {
	name     string
	srcPath  string
	destPath string
	selected bool
}

func (i mediaItem) selectedBox() string {
	if !i.selected {
		return " "
	}
	return style.TextHighlightStyle.Render("✓")
}

func (i mediaItem) fileName() string { return i.name }

func (i mediaItem) FilterValue() string { return i.srcPath }

type mediaItemDelegate struct{}

func (d mediaItemDelegate) Height() int { return 1 }

func (d mediaItemDelegate) Spacing() int { return 0 }

func (d mediaItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d mediaItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var i, ok = item.(mediaItem)
	if !ok {
		return
	}

	cursor := style.SelectedItemStyle.Render(">")

	fn := func(s ...string) string {
		o := style.ItemStyle.Render("[")
		c := style.ItemTextStyle.Render("] ")
		return o + i.selectedBox() + c + style.ItemTextStyle.Render(strings.Join(s, " "))
	}
	if index == m.Index() {
		fn = func(s ...string) string {
			o := style.TextHighlightStyle.Render(" [")
			c := style.TextHighlightStyle.Render("]")
			return cursor + o + i.selectedBox() + c + style.SelectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprintf(w, fn(i.fileName()))
}
