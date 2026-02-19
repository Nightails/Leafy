package media

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	app "github.com/nightails/leafy/internal/app"
	style "github.com/nightails/leafy/internal/style"
)

type state int

const (
	idle state = iota
	scan
	transfer
	quit
)

const quitDelay = 1 * time.Second

type Model struct {
	state       state
	mountPoints []string
	mediaList   list.Model
	timer       style.MinDuration
	spinner     spinner.Model // loading spinner
	err         error
}

func NewModel() Model {
	s := style.NewLineSpinner()

	l := list.New([]list.Item{}, mediaItemDelegate{}, 0, 10)
	l.SetShowTitle(false)
	l.SetShowPagination(true)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	// start timer for the first scan
	t := style.MinDuration{Min: quitDelay}
	t.StartNow()

	return Model{
		state:     scan, // the first scan on init
		mediaList: l,
		timer:     t,
		spinner:   s,
		err:       nil,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,              // start scanning spinner
		scanMediaCmd(m.mountPoints), // start scanning for media files
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	default:
		if (m.state == scan || m.state == transfer) && m.err == nil {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
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
			if m.state == idle && m.err == nil {
				var cmd tea.Cmd
				m.mediaList, cmd = m.mediaList.Update(msg)
				return m, cmd
			}
			return m, nil
		case "ctrl+c":
			m.state = quit
			return m, app.AfterCmd(quitDelay, app.QuitNowMsg{})
		}
	// handle app messages
	case app.StateMsg:
		m.mountPoints = msg.State.MountPoints
		return m, nil
	case mediaMsg:
		var media []list.Item
		for _, m := range msg {
			media = append(media, mediaItem{
				srcPath:      m,
				destPath:     "",
				transferring: false,
				spinnerFrame: "",
			})
		}
		m.mediaList.SetItems(media)
		return m, app.AfterCmd(m.timer.Remaining(), app.FinishedMsg{})
	case app.FinishedMsg:
		m.state = idle
		return m, nil
	case app.QuitNowMsg:
		return m, tea.Quit
	}
}

func (m Model) View() string {
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
