package tui_media

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	app "github.com/nightails/leafy/internal/tui_app"
	style "github.com/nightails/leafy/internal/tui_style"
)

type state int

const (
	idle state = iota
	scan
	transfer
	quit
)

type MediaModel struct {
	state     state
	mediaList list.Model
	timer     style.MinDuration
	spinner   spinner.Model // loading spinner
	err       error
}

func NewMediaModel() MediaModel {
	s := style.NewLineSpinner()

	l := list.New([]list.Item{}, list.DefaultDelegate{}, 0, 10)
	l.SetShowTitle(false)
	l.SetShowPagination(true)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	// start timer for the first scan
	t := style.MinDuration{Min: 1000}
	t.StartNow()

	return MediaModel{
		state:     idle, // TODO: switch to scan for first scan on init
		mediaList: l,
		timer:     t,
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
	case app.ErrMsg:
		m.err = msg
		return m, nil
	case tea.WindowSizeMsg:
		m.mediaList.SetWidth(msg.Width)
		return m, nil
	// handle key presses
	case tea.KeyMsg:
		switch msg.String() {
		default:
			return m, nil
		case "ctrl+c":
			return m, app.AfterCmd(m.timer.Remaining(), app.QuitNowMsg{})
		}
	}
}

func (m MediaModel) View() string {
	if m.state == quit {
		return "\n" + style.TextStyle.Render("Quitting leafy...")
	}
	if m.err != nil {
		return style.TextStyle.Render(m.err.Error())
	}

	var b strings.Builder
	b.WriteString("\n" + m.mediaList.View())
	return b.String()
}
