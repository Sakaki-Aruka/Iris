package command

import (
	"Iris/core/structs"
	"Iris/core/util"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func GetProfileFromConfigDir() ([]structs.Profile, error) {
	var profiles []structs.Profile
	dir, err := util.GetProfileDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.Type().IsDir() {
			continue
		}

		d, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}

		var profile structs.Profile
		if err := json.Unmarshal(d, &profile); err != nil {
			log.Println("[profile] failed to parse profile file. :", filepath.Join(dir, entry.Name()))
			continue
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}
