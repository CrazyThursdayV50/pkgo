package file

import (
	"os"
)

func WriteFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	n, err := file.WriteAt(data, 0)
	if err != nil {
		return err
	}

	err = file.Truncate(int64(n))
	if err != nil {
		return err
	}

	return nil
}
