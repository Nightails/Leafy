package device

import (
	"encoding/json"
	"os/exec"
)

type LSBLK struct {
	BlockDevices []BlockDevice `json:"blockdevices"`
}
type BlockDevice struct {
	Name        string        `json:"name"`
	Path        string        `json:"path"`
	Label       string        `json:"label"`
	Tran        string        `json:"tran"`
	Type        string        `json:"type"`
	Model       string        `json:"model"`
	Mountpoints []string      `json:"mountpoints"`
	Children    []BlockDevice `json:"children"`
}

func readLSBLK() (LSBLK, error) {
	args := []string{
		"-J",
		"--tree",
		"-f",
		"-o", "NAME,PATH,LABEL,TRAN,TYPE,MODEL,MOUNTPOINTS",
	}
	cmd := exec.Command("lsblk", args...)
	raw, err := cmd.Output()
	if err != nil {
		return LSBLK{}, err
	}

	var blk LSBLK
	if err := json.Unmarshal(raw, &blk); err != nil {
		return LSBLK{}, err
	}
	return blk, nil
}
