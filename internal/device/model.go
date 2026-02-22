package device

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/app"
	"github.com/nightails/leafy/internal/style"
)

type state int

const (
	idle state = iota
	scan
	mount
	quit
)

type Model struct {
	mountingIndex int
	deviceList    list.Model
	timer         style.MinDuration
	spinner       spinner.Model // loading spinner
	state         state
	err           error
}

func NewModel() Model {
	s := style.NewLineSpinner()

	// the list can only show 4 items at a time, pagination is enabled
	l := list.New([]list.Item{}, deviceItemDelegate{}, 0, 10)
	l.SetShowTitle(false)
	l.SetShowPagination(true)
	l.SetShowStatusBar(true)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	// start timer for the first scan
	t := style.MinDuration{Min: style.LoadDelay}
	t.StartNow()

	return Model{
		mountingIndex: -1,
		deviceList:    l,
		timer:         t,
		spinner:       s,
		state:         scan,
		err:           nil,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,      // start scanning spinner
		scanUSBDevicesCmd(), // start scanning for USB devices immediately
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	default:
		// keep the spinner running while scanning or mounting
		if (m.state == scan || m.state == mount) && m.err == nil {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)

			if m.state == mount && m.mountingIndex >= 0 && m.mountingIndex < len(m.deviceList.Items()) {
				if it, ok := m.deviceList.Items()[m.mountingIndex].(deviceItem); ok {
					it.spinnerFrame = m.spinner.View()
					it.mounting = true
					_ = m.deviceList.SetItem(m.mountingIndex, it)
				}
			}

			return m, cmd
		}
		return m, nil
	case app.ErrMsg:
		m.err = msg
		return m, nil
	case tea.WindowSizeMsg:
		m.deviceList.SetWidth(msg.Width)
		return m, nil
	// handle key presses
	case tea.KeyMsg:
		switch msg.String() {
		default:
			if m.state == idle && m.err == nil {
				var cmd tea.Cmd
				m.deviceList, cmd = m.deviceList.Update(msg)
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
				m.spinner.Tick,      // restart the spinner
				scanUSBDevicesCmd(), // start scanning
			)
		case " ": // space bar, mount/unmount select device
			if m.state == quit {
				return m, nil
			}
			m.err = nil
			m.state = mount
			m.timer.StartNow()

			d := m.deviceList.SelectedItem().(deviceItem).usb
			m.mountingIndex = m.deviceList.Index()
			if m.mountingIndex >= 0 && m.mountingIndex < len(m.deviceList.Items()) {
				it := m.deviceList.Items()[m.mountingIndex].(deviceItem)
				it.mounting = true
				it.spinnerFrame = m.spinner.View()
				_ = m.deviceList.SetItem(m.mountingIndex, it)
			}

			var cmds []tea.Cmd
			cmds = append(cmds, m.spinner.Tick) // restart the spinner
			if d.Mountpoint == "" {
				cmds = append(cmds, mountUSBDeviceCmd(d))
			} else {
				cmds = append(cmds, unmountUSBDeviceCmd(d))
			}
			return m, tea.Batch(cmds...)
		}
	case usbDevicesMsg:
		var cmds []tea.Cmd
		items := make([]list.Item, 0, len(msg))
		for _, d := range msg {
			items = append(items, deviceItem{
				usb:          d,
				mounting:     false,
				spinnerFrame: "",
			})
			// already mounted device
			if d.Mountpoint != "" {
				cmds = append(cmds, func() tea.Msg { return app.DeviceMountedMsg{MountPoint: d.Mountpoint} })
			}
		}
		m.deviceList.SetItems(items)
		cmds = append(cmds, app.AfterCmd(m.timer.Remaining(), app.FinishedMsg{}))
		return m, tea.Batch(cmds...)
	case mountUSBDeviceMsg:
		m.deviceList.SetItem(m.mountingIndex, deviceItem{
			usb:          USBDevice(msg),
			mounting:     true,
			spinnerFrame: m.spinner.View(),
		})
		m.err = nil
		return m, tea.Batch(
			func() tea.Msg { return app.DeviceMountedMsg{MountPoint: msg.Mountpoint} },
			app.AfterCmd(m.timer.Remaining(), app.FinishedMsg{}),
		)
	case unmountUSBDeviceMsg:
		m.deviceList.SetItem(m.mountingIndex, deviceItem{
			usb:          USBDevice(msg),
			mounting:     true,
			spinnerFrame: m.spinner.View(),
		})
		m.err = nil
		return m, tea.Batch(
			func() tea.Msg { return app.DeviceUnmountedMsg{MountPoint: msg.Mountpoint} },
			app.AfterCmd(m.timer.Remaining(), app.FinishedMsg{}),
		)
	case app.FinishedMsg:
		if m.mountingIndex >= 0 && m.mountingIndex < len(m.deviceList.Items()) {
			if it, ok := m.deviceList.Items()[m.mountingIndex].(deviceItem); ok {
				it.mounting = false
				_ = m.deviceList.SetItem(m.mountingIndex, it)
			}
		}
		m.mountingIndex = -1
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

	if m.state == scan {
		b.WriteString(fmt.Sprintf("\n%s%s\n\n", m.spinner.View(), style.TextStyle.Render("Scanning for USB devices...")))
		b.WriteString("\n\n" + helpBarView())
		return b.String()
	}

	if len(m.deviceList.Items()) == 0 {
		b.WriteString("\n" + style.TextStyle.Render("No USB devices found."))
		b.WriteString("\n\n" + helpBarView())
		return b.String()
	}

	b.WriteString("\n" + m.deviceList.View())
	b.WriteString("\n" + helpBarView())
	return b.String()
}
