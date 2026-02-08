package main

import (
	"fmt"
	"leafy/internal/usb"
)

func main() {
	devs, err := usb.FindUSBDevices()
	if err != nil {
		panic(err)
	}

	if len(devs) == 0 {
		fmt.Println("No USB devices found")
		return
	}

	for _, d := range devs {
		if len(d.Children) == 0 {
			fmt.Printf("Device %s has no partition\n", d.Name)
			continue
		}
		for _, p := range d.Children {
			fmt.Printf("Device %s has partition %s\n", d.Name, p.Label)
			if p.Mountpoints != nil && len(p.Mountpoints) > 0 {
				for _, m := range p.Mountpoints {
					if m != "" {
						fmt.Printf("Mounted on %s\n", m)
					}
				}
			}
		}
	}
}
