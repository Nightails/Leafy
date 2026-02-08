package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/usb"
)

type errMsg error
type usbDevicesMsg []usb.BlockDevice
type finishScanMsg struct{}
type quitNowMsg struct{}

const (
	minScanDuration = 2 * time.Second
	quitDelay       = 1 * time.Second
)

type USBModel struct {
	spinner       spinner.Model
	devList       list.Model
	scanning      bool
	scanStartedAt time.Time
	quitting      bool
	err           error
}

func NewUSBModel() USBModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = spinnerStyle

	// the list can show only 4 items at a time, pagination is enabled
	l := list.New([]list.Item{}, usbItemDelegate{}, 0, 4)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	l.DisableQuitKeybindings()
	l.SetFilteringEnabled(false)

	return USBModel{
		spinner:  s,
		devList:  l,
		scanning: true,
	}
}

func scanUSBDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		devs, err := usb.FindUSBDevices()
		if err != nil {
			return errMsg(err)
		}
		return usbDevicesMsg(devs)
	}
}

func finishAfterCmd(d time.Duration) tea.Cmd {
	if d <= 0 {
		return func() tea.Msg { return finishScanMsg{} }
	}
	return tea.Tick(d, func(t time.Time) tea.Msg { return finishScanMsg{} })
}

func quitAfterCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg { return quitNowMsg{} })
}

func (m USBModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,      // start spinning the spinner
		scanUSBDevicesCmd(), // start scanning for USB devices immediately
	)
}

func (m USBModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.devList.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			m.scanning = false
			return m, quitAfterCmd(quitDelay)
		case "s":
			if m.quitting {
				return m, nil
			}
			m.err = nil
			m.scanning = true
			m.scanStartedAt = time.Now()
			return m, tea.Batch(
				m.spinner.Tick,      // restart the spinner
				scanUSBDevicesCmd(), // start scanning
			)
		default:
			if !m.scanning && !m.quitting && m.err == nil {
				var cmd tea.Cmd
				m.devList, cmd = m.devList.Update(msg)
				return m, cmd
			}
			return m, nil
		}
	case usbDevicesMsg:
		// if scanStartedAt was never set, assume an initial scan
		if m.scanStartedAt.IsZero() {
			m.scanStartedAt = time.Now()
		}

		// Convert devices to list items
		devs := msg
		items := make([]list.Item, 0, len(devs))
		for _, d := range devs {
			items = append(items, usbDeviceItem{dev: d})
		}
		m.devList.SetItems(items)

		remaining := minScanDuration - time.Since(m.scanStartedAt)
		return m, finishAfterCmd(remaining)
	case errMsg:
		// save error immediately, but optionally keep the spinner alive for the minScanDuration too
		m.err = msg
		remaining := minScanDuration - time.Since(m.scanStartedAt)
		return m, finishAfterCmd(remaining)
	case finishScanMsg:
		m.scanning = false // stop spinning the spinner
		return m, nil
	case quitNowMsg:
		return m, tea.Quit
	default:
		// spinner only alive while scanning
		if m.scanning && !m.quitting {
			// ensure scanStartedAt is set on the first update
			if m.scanStartedAt.IsZero() {
				m.scanStartedAt = time.Now()
			}
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	}
}

func (m USBModel) View() string {
	if m.quitting {
		return "\n" + textStyle.Render("Quitting...")
	}

	if m.err != nil && !m.scanning {
		return textStyle.Render(m.err.Error())
	}

	var b strings.Builder

	if m.scanning {
		b.WriteString(fmt.Sprintf("\n%s%s\n\n", m.spinner.View(), textStyle.Render("Scanning for USB devices...")))
	} else {
		b.WriteString("\n" + textStyle.Render("Scanning complete.") + "\n\n")
		if len(m.devList.Items()) == 0 {
			b.WriteString(textStyle.Render("No USB devices found.") + "\n\n")
			b.WriteString(helpBarView())
			return b.String()
		}

		b.WriteString(textStyle.Render("USB Devices:") + "\n\n")
		b.WriteString(m.devList.View())
		b.WriteString("\n\n" + helpBarView())
		return b.String()
	}

	b.WriteString("\n\n" + helpBarView())
	return b.String()
}
