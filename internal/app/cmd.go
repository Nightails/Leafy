package app

import (
	tea "github.com/charmbracelet/bubbletea"
	dev "github.com/nightails/leafy/internal/device"
	"github.com/nightails/leafy/internal/file"
)

// initDevicesCmd finds and mounts usb devices.
func initDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		var mdevs []device
		devs, err := dev.FindUSBDevices()
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
		return devicesMsg(mdevs)
	}
}

// removeDevicesCmd unmounts all usb devices.
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

// findMediaCmd searches for supported file formats and returns a list.
func findMediaCmd(devices []device) tea.Cmd {
	return func() tea.Msg {
		if len(devices) == 0 {
			return nil
		}

		var paths []string
		for _, d := range devices {
			paths = append(paths, d.mountpoint)
		}

		var media []medium
		files, err := file.GetFiles(paths, formats)
		if err != nil {
			return errMsg(err)
		}
		if len(files) == 0 {
			return nil
		}
		for _, f := range files {
			media = append(media, medium{
				name: f.Name,
				src:  f.Path,
			})
		}
		return mediaMsg(media)
	}
}

// copyMedia copies given media to destination path.
func copyMedia(media []medium) tea.Cmd {
	return func() tea.Msg {
		// TODO: implements copying file logic
		return nil
	}
}
