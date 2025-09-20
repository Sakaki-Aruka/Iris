package cmd

import (
	"Iris/cmd/profile"
	"Iris/cmd/session"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "iris",
	Short: "Iris - Minecraft server session manager in Go",
	Long:  "Iris is a CLI tool to manager Minecraft server sessions with profiles.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(profile.Cmd)
	rootCmd.AddCommand(session.Cmd)
}
