package usb

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
