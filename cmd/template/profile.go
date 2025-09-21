package template

import (
	"Iris/internal/profile"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var ProfileTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Show a template",
	Run: func(cmd *cobra.Command, args []string) {
		s := profile.Schedule{
			Timing: "0 0 * * *",
			Script: "echo Good night Iris",
		}

		a := profile.AutoRestartCfg{
			Max:           3,
			TriggerCodes:  []int{1, 2, 3},
			StartupScript: "echo Hello Iris again",
			Scheduled:     []profile.Schedule{s},
		}

		p := profile.Profile{
			Name:          "Template",
			KeepLogs:      false,
			StartupScript: "echo Hello Iris",
			AutoRestart:   a,
		}

		j, _ := json.MarshalIndent(p, "", "  ")
		fmt.Println(string(j))
	},
}
