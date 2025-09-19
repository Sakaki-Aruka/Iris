package profile

import (
	"encoding/json"
	"os"
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
	Timing  string `json:"timing"`
	Profile string `json:"profile"`
}

var Profiles []Profile = []Profile{}

func NewDefaultProfile(name string) *Profile {
	return &Profile{
		Name:          name,
		KeepLogs:      false,
		StartupScript: "echo Hello Iris",
		AutoRestart: AutoRestartCfg{
			Max:          0,
			TriggerCodes: []int{},
		},
	}
}

func SaveProfile(filename string, p *Profile) error {
	data, err := json.MarshalIndent(p, "", "")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

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

func DeleteProfile(filename string) error {
	return os.Remove(filename)
}
