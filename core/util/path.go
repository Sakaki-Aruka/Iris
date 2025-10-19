package util

import (
	"os"
	"path/filepath"
)

func GetConfigDir() (string, error) {
	c, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	c = filepath.Join(c, "iris")
	return c, nil
}

func GetProfileDir() (string, error) {
	base, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "profiles"), nil
}
