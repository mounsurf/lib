package file

import (
	"errors"
	"fmt"
	"os"
)

func MkDir(path string) error {
	// check
	if s, err := os.Stat(path); err == nil {
		if s.IsDir() {
			return nil
		} else {
			return errors.New(fmt.Sprintf("File \"%s\" exists and is not dir", path))
		}
	} else {
		err := os.MkdirAll(path, 0766)
		return err
	}
}
