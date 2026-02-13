package tui_media

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/tui_style"
)

type MediaModel struct {
	mediaList list.Model
	timer     tui_style.MinDuration
	spinner   spinner.Model // loading spinner
	err       error
}

func NewMediaModel() MediaModel {
	s := tui_style.NewLineSpinner()

	l := list.New([]list.Item{}, list.DefaultDelegate{}, 0, 20)
	l.Title = "Media:"
	l.SetShowTitle(true)
	l.Styles.Title = tui_style.ItemTitleStyle
	l.SetShowPagination(true)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	return MediaModel{
		mediaList: l,
		spinner:   s,
		err:       nil,
	}
}

func (m MediaModel) Init() tea.Cmd {
	return nil
}

func (m MediaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	default:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		default:
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	}
}

func (m MediaModel) View() string {
	return "Media View"
}
