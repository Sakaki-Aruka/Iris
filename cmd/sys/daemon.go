package sys

import (
	"Iris/daemon"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "sys",
	Short: "Iris daemon command",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start Iris daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := daemon.StartDaemon(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(startCmd)
}
