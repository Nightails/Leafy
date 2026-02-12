package usb

import (
	"errors"
)

type MountedUSBPartition struct {
	Device     BlockDevice
	Mountpoint string
}

// FindUSBDevices returns a list of USB devices by running lsblk.
func FindUSBDevices() ([]BlockDevice, error) {
	blk, err := readLSBLK()
	if err != nil {
		return nil, err
	}

	devs := make([]BlockDevice, 0)
	for _, d := range blk.BlockDevices {
		if d.Type == "disk" && d.Tran == "usb" {
			devs = append(devs, d)
		}
	}
	return devs, nil
}

// MountUSBPartition mounts all mountable partitions of a USB device and creates a mountpoint for each.
func MountUSBPartition(bd BlockDevice) ([]MountedUSBPartition, error) {
	if !HasPartitions(bd) {
		return nil, errors.New("device has no mountable partition")
	}

	mounted := make([]MountedUSBPartition, 0)

	for _, c := range bd.Children {
		if IsMounted(c) {
			continue // already mounted, skip
		}
		if c.Type != "part" {
			continue // skip non-partition devices
		}

		mountPoint, err := mountDevice(c.Path)
		if err != nil {
			return nil, err
		}

		mounted = append(mounted, MountedUSBPartition{c, mountPoint})
	}

	return mounted, nil
}

// UnmountUSBPartition unmounts a mounted USB partition and removes its mountpoint.
func UnmountUSBPartition(m MountedUSBPartition) error {
	if err := unmountDevice(m.Device.Path); err != nil {
		return err
	}
	return nil
}

func HasPartitions(bd BlockDevice) bool {
	return bd.Children != nil && len(bd.Children) > 0
}

func IsMounted(bd BlockDevice) bool {
	if bd.Mountpoints == nil || len(bd.Mountpoints) == 0 {
		return false
	}
	return true
}
