package util

import (
	"Iris/ipc"
	"os"
	"path/filepath"
)

func GetConfDir() (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	d = filepath.Join(d, "iris")
	return d, nil
}

func GetProfileDir() (string, error) {
	d, err := GetConfDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "profile"), nil
}

func GetSystemDir() (string, error) {
	d, err := GetConfDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "sys"), nil
}

func GetSocketPath() (string, error) {
	d, err := GetSystemDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, ipc.SocketName), nil
}

func GetDaemonPidPath() (string, error) {
	d, err := GetSystemDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, ipc.PidFileName), nil
}
