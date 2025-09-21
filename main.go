package main

import (
	"Iris/cmd"
	"Iris/internal/profile"
	"Iris/util"
	"fmt"
	"os"
)

func main() {

	if err := CreateConfigDir(); err != nil {
		os.Exit(1)
	}

	if err := CreateProfileDir(); err != nil {
		os.Exit(1)
	}

	if err := profile.LoadExistProfiles(); err != nil {
		fmt.Println("failed to load profiles")
		os.Exit(1)
	}

	cmd.Execute()
}

func CreateConfigDir() error {
	confDir, err := util.GetConfDir()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if _, err := os.Stat(confDir); err != nil { // exists check
		// '~/.config/iris' not exists.
		mkdErr := os.Mkdir(confDir, 0770)
		if mkdErr != nil {
			fmt.Printf("failed to create Iris config dir. (%s)\n", confDir)
			return mkdErr
		}
	}
	return nil
}

func CreateProfileDir() error {
	confDir, err := util.GetProfileDir()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if _, err := os.Stat(confDir); err != nil { // exists check
		mkdErr := os.Mkdir(confDir, 0770)
		if mkdErr != nil {
			fmt.Printf("failed to create Iris profile config dir. (%s)\n", confDir)
			return mkdErr
		}
	}
	return nil
}
