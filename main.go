package main

import (
	"Iris/cmd"
	"Iris/internal/profile"
	"Iris/util"
	"fmt"
	"os"
)

func main() {

	if err := check(util.GetConfDir()); err != nil {
		os.Exit(1)
	}

	if err := check(util.GetProfileDir()); err != nil {
		os.Exit(1)
	}

	if err := check(util.GetSystemDir()); err != nil {
		os.Exit(1)
	}

	if err := profile.LoadExistProfiles(); err != nil {
		fmt.Println("failed to load profiles")
		os.Exit(1)
	}

	cmd.Execute()
}

func check(d string, err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	if _, err := os.Stat(d); err != nil {
		if mkdErr := os.Mkdir(d, 0770); mkdErr != nil {
			fmt.Printf("failed to create %s\n", d)
			return mkdErr
		}
	}
	return nil
}
