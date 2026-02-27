package app

import (
	tea "github.com/charmbracelet/bubbletea"
	dev "github.com/nightails/leafy/internal/device"
)

// initDevices finds and mounts usb devices.
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

// removeDevices unmounts all usb devices.
func removeDevicesCmd(devs []device) tea.Cmd {
	return func() tea.Msg {
		if len(devs) == 0 {
			return nil
		}

		for _, d := range devs {
			if _, err := dev.UnmountDevice(dev.USBDevice{
				Name:       d.name,
				Path:       d.path,
				Mountpoint: d.mountpoint,
			}); err != nil {
				return errMsg(err)
			}
		}
		return nil
	}
}

// findMedia searches for supported media formats and returns a list.
func findMediaCmd() tea.Cmd {
	return func() tea.Msg {
		// TODO: implements find media logic
		return nil
	}
}

// copyMedia copies given media to destination path.
func copyMedia(media []medium) tea.Cmd {
	return func() tea.Msg {
		// TODO: implements copying media logic
		return nil
	}
}
