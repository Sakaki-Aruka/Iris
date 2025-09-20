package profile

import (
	"Iris/internal/profile"
	"fmt"
	"path/filepath"

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
		p, err := profile.LoadProfile(filename)
		if err != nil {
			return fmt.Errorf("failed to load profile defined file")
		}
		profile.AddProfileWithFile(*p, filename)
		return nil
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete [file]",
	Short: "Delete a profile JSON file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]
		p, pExists := profile.Profiles[profileName]
		if pExists {
			delete(profile.Profiles, p.Name)
		}
		path, fExists := profile.ProfilesWithPath[profileName]
		if fExists {
			delete(profile.ProfilesWithPath, p.Name)
			if err := profile.DeleteProfile(path); err != nil {
				return err
			}
		}
		return nil
	},
}

var profileUpdateCmd = &cobra.Command{
	Use:   "update [file]",
	Short: "Reload/Update a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := profile.GetProfileDir()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		path := filepath.Join(configDir, args[0])
		p, err := profile.LoadProfile(path)
		if err != nil {
			return err
		}

		fmt.Printf("Profile loaded: %+v\n", p)
		profile.AddProfileWithFile(*p, path)
		return nil
	},
}

func init() {
	Cmd.AddCommand(profileCreateCmd)
	Cmd.AddCommand(profileDeleteCmd)
	Cmd.AddCommand(profileUpdateCmd)
}
