package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Nightails/Leafy/internal/usb"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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
	devices       []usb.BlockDevice
	spinner       spinner.Model
	scanning      bool
	scanStartedAt time.Time
	quitting      bool
	err           error
}

func NewUSBModel() USBModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = spinnerStyle
	return USBModel{
		spinner:  s,
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
			m.devices = nil
			m.scanning = true
			m.scanStartedAt = time.Now()
			return m, tea.Batch(
				m.spinner.Tick,      // restart the spinner
				scanUSBDevicesCmd(), // start scanning
			)
		default:
			return m, nil
		}
	case usbDevicesMsg:
		// if scanStartedAt was never set, assume an initial scan
		if m.scanStartedAt.IsZero() {
			m.scanStartedAt = time.Now()
		}
		// save results immediately, but keep "scanning" until minScanDuration has passed
		m.devices = []usb.BlockDevice(msg)
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
		return "Quitting..."
	}

	if m.err != nil && !m.scanning {
		return m.err.Error()
	}

	var b strings.Builder

	if m.scanning {
		b.WriteString(fmt.Sprintf("%s Scanning USB devices...\n\n", m.spinner.View()))
	} else {
		b.WriteString("Scanning complete. Press 's' to scan again.\n\n")
		if len(m.devices) == 0 {
			b.WriteString("No USB devices found.\n")
			return b.String()
		}

		for _, dev := range m.devices {
			b.WriteString(fmt.Sprintf("%s %s\n", dev.Name, dev.Path))
		}
	}

	b.WriteString("\n\n" + helpBar)

	return b.String()
}
