package tui_device

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/device"
	app "github.com/nightails/leafy/internal/tui_app"
)

type usbDevicesMsg []device.USBDevice

// scanUSBDevicesCmd returns a command that scans for USB devices and sends the results as a usbDevicesMsg
func scanUSBDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		devs, err := device.FindUSBDevices()
		if err != nil {
			return app.ErrMsg(err)
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
			return app.ErrMsg(err)
		}
		return mountUSBDeviceMsg(md)
	}
}
