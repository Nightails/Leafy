package app

import (
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type mediumItem struct {
	medium   medium
	selected bool
}

func (mi mediumItem) FilterValue() string {
	return mi.medium.name
}

type mediaItemDelegate struct{}

func (mi mediaItemDelegate) Height() int                             { return 1 }
func (mi mediaItemDelegate) Spacing() int                            { return 0 }
func (mi mediaItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (mi mediaItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var i, ok = item.(mediumItem)
	if !ok {
		return
	}
	if index == m.Index() {
		_, _ = io.WriteString(w, "> ")
	}
	if i.selected {
		_, _ = io.WriteString(w, "[x] ")
	} else {
		_, _ = io.WriteString(w, "[ ] ")
	}

	fileSize := formatFileSize(i.medium.total)
	_, _ = io.WriteString(w, i.medium.name+" ("+fileSize+")")
}
