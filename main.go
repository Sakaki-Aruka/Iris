package main

import (
	"Iris/cmd"
	"Iris/internal/profile"
	"fmt"
	"os"
)

func main() {
	cmd.Execute()
	if err := profile.LoadExistProfiles(); err != nil {
		fmt.Println("failed to load profiles")
		os.Exit(1)
	}
}
