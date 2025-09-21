package template

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "template",
	Short: "Shows setting templates",
}

func init() {
	Cmd.AddCommand(ProfileTemplateCmd)
}
