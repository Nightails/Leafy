package device

type USBDevice struct {
	Name       string // the Model the device
	Path       string // the partition path of the device, assumed to be the first and only partition
	Label      string // the label of the partition, assumed to be the first and only partition
	Mountpoint string // the mountpoint of the partition, assumed to be the first and only partition
}

// FindUSBDevices returns a list of USB devices by running lsblk.
func FindUSBDevices() ([]USBDevice, error) {
	blk, err := readLSBLK()
	if err != nil {
		return nil, err
	}

	devs := make([]USBDevice, 0)
	for _, bd := range blk.BlockDevices {
		// filter out non-USB disks
		if bd.Type == "disk" && bd.Tran == "usb" {
			// filter out non-partitioned disks
			if hasPartition(bd) {
				mp := ""
				if len(bd.Children[0].Mountpoints) > 0 {
					mp = bd.Children[0].Mountpoints[0]
				}
				devs = append(devs, USBDevice{
					bd.Model,
					bd.Children[0].Path,
					bd.Children[0].Label,
					mp,
				})
			}
		}
	}

	return devs, nil
}

func hasPartition(bd BlockDevice) bool {
	return bd.Children != nil && len(bd.Children) > 0
}

// MountDevice mounts the given USB device using udisksctl and returns the updated device.
func MountDevice(d USBDevice) (USBDevice, error) {
	// already mounted, skip
	if d.Mountpoint != "" {
		return d, nil
	}

	mp, err := mountUdisks(d.Path)
	if err != nil {
		return d, err
	}
	d.Mountpoint = mp
	return d, nil
}

// UnmountDevice unmounts the given USB device using udisksctl and returns the updated device.
func UnmountDevice(d USBDevice) (USBDevice, error) {
	// not mounted, skip
	if d.Mountpoint == "" {
		return d, nil
	}

	if err := unmountUdisks(d.Path); err != nil {
		return d, err
	}
	d.Mountpoint = ""
	return d, nil
}

// PowerOffDevice powers off the given USB device using udisksctl.
func PowerOffDevice(d USBDevice) error {
	if d.Mountpoint != "" {
		if err := unmountUdisks(d.Path); err != nil {
			return err
		}
	}
	return powerOffUdisks(d.Path)
}
