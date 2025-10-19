package command

import (
	"Iris/core/daemon"

	"github.com/spf13/cobra"
)

var DaemonCommand = &cobra.Command{
	Use:   "daemon",
	Short: "daemon command",
}

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "daemon start",
	Run: func(_ *cobra.Command, _ []string) {
		daemon.Start()
	},
}

func init() {
	DaemonCommand.AddCommand(startCommand)
}
