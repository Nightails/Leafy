package usb

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const baseMountPoint = "/mnt/leafy/"

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
	if bd.Children == nil || len(bd.Children) == 0 {
		return nil, errors.New("device has no mountable partition")
	}

	mounted := make([]MountedUSBPartition, 0)

	for _, c := range bd.Children {
		if c.Mountpoints != nil && len(c.Mountpoints) > 0 {
			continue // already mounted, skip
		}
		if c.Type != "part" {
			continue // skip non-partition devices
		}

		src := c.Path
		dst := baseMountPoint + strings.ToLower(c.Label)

		if err := os.MkdirAll(dst, 0700); err != nil {
			return nil, err
		}

		cmd := exec.Command("mount", src, dst)
		if out, err := cmd.CombinedOutput(); err != nil {
			return nil, fmt.Errorf("%v: %s", err, string(out))
		}

		mounted = append(mounted, MountedUSBPartition{c, dst})
	}

	return mounted, nil
}

// UnmountUSBPartition unmounts a mounted USB partition and removes its mountpoint.
func UnmountUSBPartition(m MountedUSBPartition) error {
	cmd := exec.Command("umount", m.Mountpoint)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%v: %s", err, string(out))
	}
	if err := os.RemoveAll(m.Mountpoint); err != nil {
		return fmt.Errorf("failed to remove mountpoint: %v", err)
	}
	return nil
}
