package tui_device

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/device"
	"github.com/nightails/leafy/internal/tui_style"
)

type state int

const (
	idle state = iota
	scan
	mount
	quit
)

const (
	loadDelay = 2 * time.Second
	quitDelay = 1 * time.Second
)

type DeviceModel struct {
	mountingIndex int
	deviceList    list.Model
	timer         tui_style.MinDuration
	spinner       spinner.Model // loading spinner
	state         state
	err           error
}

func NewDeviceModel() DeviceModel {
	s := tui_style.NewLineSpinner()

	// the list can only show 4 items at a time, pagination is enabled
	l := list.New([]list.Item{}, deviceItemDelegate{}, 0, 8)
	l.Title = "Devices:"
	l.SetShowTitle(true)
	l.Styles.Title = tui_style.ItemTitleStyle
	l.SetShowPagination(true)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	// start timer for the first scan
	t := tui_style.MinDuration{Min: loadDelay}
	t.StartNow()

	return DeviceModel{
		mountingIndex: -1,
		deviceList:    l,
		timer:         t,
		spinner:       s,
		state:         scan,
		err:           nil,
	}
}

func (m DeviceModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,      // start scanning spinner
		scanUSBDevicesCmd(), // start scanning for USB devices immediately
	)
}

func (m DeviceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case errMsg:
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
		case "ctrl+c":
			m.state = quit
			return m, afterCmd(quitDelay, quitNowMsg{})
		case "s":
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
		case "enter", "return":
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
			return m, tea.Batch(
				m.spinner.Tick,       // restart the spinner
				mountUSBDeviceCmd(d), // start mounting
			)
		}
	case usbDevicesMsg:
		items := make([]list.Item, 0, len(msg))
		for _, d := range msg {
			items = append(items, deviceItem{
				usb:          d,
				mounting:     false,
				spinnerFrame: "",
			})
		}
		m.deviceList.SetItems(items)
		return m, afterCmd(m.timer.Remaining(), finishedMsg{})
	case mountUSBDeviceMsg:
		m.deviceList.SetItem(m.mountingIndex, deviceItem{
			usb:          device.USBDevice(msg),
			mounting:     true,
			spinnerFrame: m.spinner.View(),
		})
		m.err = nil
		return m, afterCmd(m.timer.Remaining(), finishedMsg{})
	case finishedMsg:
		if m.mountingIndex >= 0 && m.mountingIndex < len(m.deviceList.Items()) {
			if it, ok := m.deviceList.Items()[m.mountingIndex].(deviceItem); ok {
				it.mounting = false
				_ = m.deviceList.SetItem(m.mountingIndex, it)
			}
		}
		m.mountingIndex = -1
		m.state = idle
		return m, nil
	case quitNowMsg:
		return m, tea.Quit
	}
}

func (m DeviceModel) View() string {
	if m.state == quit {
		return "\n" + tui_style.TextStyle.Render("Quitting leafy...")
	}
	if m.err != nil {
		return tui_style.TextStyle.Render(m.err.Error())
	}

	var b strings.Builder

	if m.state == scan {
		b.WriteString(fmt.Sprintf("\n%s%s\n\n", m.spinner.View(), tui_style.TextStyle.Render("Scanning for USB devices...")))
		b.WriteString("\n\n" + helpBarView())
		return b.String()
	}

	if len(m.deviceList.Items()) == 0 {
		b.WriteString("\n" + tui_style.TextStyle.Render("No USB devices found."))
		b.WriteString("\n\n" + helpBarView())
		return b.String()
	}

	b.WriteString("\n" + m.deviceList.View())
	b.WriteString("\n" + helpBarView())
	return b.String()
}
