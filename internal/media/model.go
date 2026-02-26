package media

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/app"
	"github.com/nightails/leafy/internal/style"
)

type state int

const (
	idle state = iota
	scan
	quit
)

type Model struct {
	state       state
	mountPoints []string
	mediaList   list.Model
	destInput   textinput.Model
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
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	// start timer for the first scan
	t := style.MinDuration{Min: style.QuitDelay}
	t.StartNow()

	i := textinput.New()
	i.Placeholder = "~/Videos/"

	return Model{
		state:     scan, // the first scan on init
		mediaList: l,
		destInput: i,
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
		if (m.state == scan) && m.err == nil {
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
		m.destInput.Width = msg.Width
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
		case "q": // q, quit
			m.state = quit
			return m, app.AfterCmd(style.QuitDelay, app.QuitNowMsg{})
		case "s": // s, scan again
			if m.state == quit {
				return m, nil
			}
			m.err = nil
			m.state = scan
			m.timer.StartNow()
			return m, tea.Batch(
				m.spinner.Tick, // restart the spinner
				scanMediaCmd(m.mountPoints),
			)
		case " ": // space bar, add the selected file to the transfer queue
			i := m.mediaList.SelectedItem().(mediaItem)
			i.selected = !i.selected
			m.mediaList.SetItem(m.mediaList.Index(), i)
			return m, nil // TODO: add the selected file to the transfer queue
		case "enter", "return": // enter, start transferring selected files
			// TODO: start transferring selected files
			return m, nil
		}
	// handle app messages
	case app.StateMsg:
		m.mountPoints = msg.State.MountPoints
		return m, nil
	case mediaMsg:
		var media []list.Item
		for _, m := range msg {
			s := strings.Split(m, "/")
			media = append(media, mediaItem{
				name:     s[len(s)-1],
				srcPath:  m,
				destPath: "",
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
	b.WriteString("\n" + "Destination: " + m.destInput.View())
	b.WriteString("\n\n" + helpBarView())
	return b.String()
}
