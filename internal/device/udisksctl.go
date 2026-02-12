package device

import (
	"fmt"
	"os/exec"
	"regexp"
)

var (
	// matches: Mounted /dev/sda1 at /media/USER/LABLE
	reMounted = regexp.MustCompile(`^Mounted\s+(/dev/\S+)\s+at\s+(.+?)\.?\s*$`)
	// matches: Unmounted /dev/sda1 at /media/USER/LABLE
	reUnmounted = regexp.MustCompile(`^Unmounted\s+(/dev/\S+)\s+at\s+(.+?)\.?\s*$`)
)

// mountUdisks mounts the given device using udisksctl and returns the mountpoint.
func mountUdisks(device string) (string, error) {
	args := []string{"mount", "-b", device}
	cmd := exec.Command("udisksctl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, string(out))
	}
	m := reMounted.FindStringSubmatch(string(out))
	return m[1], nil
}

// unmountUdisks unmounts the given device using udisksctl.
func unmountUdisks(device string) error {
	args := []string{"unmount", "-b", device}
	cmd := exec.Command("udisksctl", args...)
	_, err := cmd.CombinedOutput()
	return err
}
