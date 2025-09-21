package sys

import (
	"Iris/daemon"
	"Iris/util"
	"os"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "sys",
	Short: "Iris daemon command",
}

var startCmd = &cobra.Command{
	Short: "start Iris daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := util.GetSocketPath()
		if err != nil {
			return err
		}
		if _, err := os.Stat(s); err != nil {
			return err
		}
		if err := daemon.StartDaemon(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(startCmd)
}
