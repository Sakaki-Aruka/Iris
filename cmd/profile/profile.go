package profile

import (
	"Iris/internal/profile"
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles",
}

var profileCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Create a new profile JSON file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		//
		//p := profile.NewDefaultProfile("default")
		//return profile.SaveProfile(filename, p)
		// TODO: impl here
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete [file]",
	Short: "Delete a profile JSON file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return profile.DeleteProfile(args[0])
	},
}

var profileUpdateCmd = &cobra.Command{
	Use:   "update [file]",
	Short: "Reload/Update a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := profile.LoadProfile(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Profile loaded: %+v\n", p)
		// TODO: impl (when a new conflicts with old)
		return nil
	},
}

func init() {
	Cmd.AddCommand(profileCreateCmd)
	Cmd.AddCommand(profileDeleteCmd)
	Cmd.AddCommand(profileUpdateCmd)
}
