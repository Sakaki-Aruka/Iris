package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "iris",
	Short: "iris is a simple minecraft server session manager.",
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("iris command execute error:", err)
		os.Exit(10)
	}
}

func init() {
	cmd.AddCommand(ListCommand)
	cmd.AddCommand(DaemonCommand)
	cmd.AddCommand(SessionCommand)
}
