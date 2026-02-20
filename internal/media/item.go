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
	srcPath      string
	destPath     string
	transferring bool
	spinnerFrame string
}

func (i mediaItem) Title() string { return i.srcPath }

func (i mediaItem) Description() string {
	if i.transferring {
		return i.spinnerFrame + style.ItemTextStyle.Render("Transferring...")
	}
	if i.destPath == "" {
		return ""
	}
	return "↳" + i.destPath
}

func (i mediaItem) FilterValue() string { return i.srcPath }

type mediaItemDelegate struct{}

func (d mediaItemDelegate) Height() int { return 2 }

func (d mediaItemDelegate) Spacing() int { return 0 }

func (d mediaItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d mediaItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var i, ok = item.(mediaItem)
	if !ok {
		return
	}

	fnTitle := style.ItemStyle.Render
	if index == m.Index() {
		fnTitle = func(s ...string) string {
			return style.SelectedItemStyle.Render(fmt.Sprintf("> %s", strings.Join(s, " ")))
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
