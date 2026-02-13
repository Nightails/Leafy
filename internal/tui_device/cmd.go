package tui_device

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/device"
)

type errMsg error

type usbDevicesMsg []device.USBDevice

// scanUSBDevicesCmd returns a command that scans for USB devices and sends the results as a usbDevicesMsg
func scanUSBDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		devs, err := device.FindUSBDevices()
		if err != nil {
			return errMsg(err)
		}
		return usbDevicesMsg(devs)
	}
}

type mountUSBDeviceMsg device.USBDevice

// mountUSBDeviceCmd returns a command that mounts the given USB device and sends the updated device as a mountUSBDeviceMsg
func mountUSBDeviceCmd(d device.USBDevice) tea.Cmd {
	return func() tea.Msg {
		md, err := device.MountDevice(d)
		if err != nil {
			return errMsg(err)
		}
		return mountUSBDeviceMsg(md)
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
