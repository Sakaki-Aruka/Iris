package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Profile struct {
	Name          string         `json:"name"`
	KeepLogs      bool           `json:"keep_logs"`
	StartupScript string         `json:"startup_script"`
	AutoRestart   AutoRestartCfg `json:"auto_restart"`
}

type AutoRestartCfg struct {
	Max           uint       `json:"max"`
	TriggerCodes  []int      `json:"trigger_codes"`
	StartupScript string     `json:"startup_script"`
	Scheduled     []Schedule `json:"scheduled"`
}

type Schedule struct {
	Timing string `json:"timing"`
	Script string `json:"script"`
}

var Profiles = make(map[string]*Profile)
var ProfilesWithPath = make(map[string]string) // key: ProfileName, value: FilePath

func LoadProfile(filename string) (*Profile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func AddProfileWithFile(p Profile, path string) {
	Profiles[p.Name] = &p
	ProfilesWithPath[p.Name] = path
}

func GetProfileDir() (string, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir")
	}
	return filepath.Join(confDir, "iris/profile"), nil
}

func LoadExistProfiles() error {
	// for tool init process
	// Load all profiles
	confDir, err := GetProfileDir()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	entries, err := os.ReadDir(confDir)
	if err != nil {
		return fmt.Errorf("failed to get files from user config dir")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := filepath.Join(confDir, entry.Name())
		p, err := LoadProfile(filename)
		if err != nil {
			continue
		}
		AddProfileWithFile(*p, filename)
	}
	return nil
}

func DeleteProfile(filename string) error {
	return os.Remove(filename)
}
