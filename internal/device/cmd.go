package device

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/app"
)

type usbDevicesMsg []USBDevice

// scanUSBDevicesCmd returns a command that scans for USB devices and sends the results as a usbDevicesMsg
func scanUSBDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		devs, err := FindUSBDevices()
		if err != nil {
			return app.ErrMsg(err)
		}
		return usbDevicesMsg(devs)
	}
}

type mountUSBDeviceMsg USBDevice

// mountUSBDeviceCmd returns a command that mounts the given USB device and sends the updated device as a mountUSBDeviceMsg
func mountUSBDeviceCmd(d USBDevice) tea.Cmd {
	return func() tea.Msg {
		md, err := MountDevice(d)
		if err != nil {
			return app.ErrMsg(err)
		}
		return mountUSBDeviceMsg(md)
	}
}
