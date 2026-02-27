package app

import (
	tea "github.com/charmbracelet/bubbletea"
	dev "github.com/nightails/leafy/internal/device"
)

// initDevices finds and mount usb devices.
func initDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		devs, err := dev.FindUSBDevices()
		var mdevs []device
		if err != nil {
			return errMsg(err)
		}
		for _, d := range devs {
			if d.Mountpoint == "" {
				d, err = dev.MountDevice(d)
				if err != nil {
					return errMsg(err)
				}
				mdevs = append(mdevs, device{
					d.Name,
					d.Path,
					d.Mountpoint,
				})
			}
		}
		return mountedDevsMsg(mdevs)
	}
}

// uninitDevices unmounts all usb devices.
func uninitDevicesCmd(devs []device) tea.Cmd {
	return func() tea.Msg {
		for _, d := range devs {
			if err := dev.UnmountDevice(d); err != nil {
				return errMsg(err)
			}
		}
		return nil
	}
}

// findMedia searchs for supported media formats and return a list.
func findMediaCmd() tea.Cmd {
	return func() tea.Msg {
		// TODO: implements find media logic
		return nil
	}
}

// copyMedia copys given media to destination path.
func copyMedia(media []medium) tea.Cmd {
	return func() tea.Msg {
		// TODO: implements copying media logic
		return nil
	}
}
