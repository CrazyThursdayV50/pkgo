package file

import (
	"io"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ReadFileToString(path string) (string, error) {
	data, err := ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
