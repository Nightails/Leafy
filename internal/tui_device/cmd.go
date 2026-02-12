package tui_device

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/device"
)

type errMsg error

type usbDevicesMsg []device.USBDevice

func scanUSBDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		devs, err := device.FindUSBDevices()
		if err != nil {
			return errMsg(err)
		}
		return usbDevicesMsg(devs)
	}
}

type finishedMsg struct{}
type quitNowMsg struct{}

// afterCmd returns a command that sends the given message after the given duration
func afterCmd(d time.Duration, msg tea.Msg) tea.Cmd {
	if d <= 0 {
		return func() tea.Msg { return msg }
	}
	return tea.Tick(d, func(time.Time) tea.Msg { return msg })
}
