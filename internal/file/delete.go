package file

import (
	"fmt"
	"os"
)

func Delete(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("not allowed to delete a directory: %s", path)
	}
	_ = os.Remove(path)
	return nil
}
