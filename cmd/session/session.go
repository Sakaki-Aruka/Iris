package session

import (
	"Iris/internal/profile"
	"Iris/internal/session"
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "session",
	Short: "Manage sessions",
}

var ( // for create
	name                     string
	startupScript            string
	profileName              string
	saveSession              bool
	autoRestart              bool
	autoRestartMax           uint
	autoRestartTriggerCode   []int
	autoRestartStartupScript string
	scheduleTiming           string
	scheduleScript           string
)

var (
	sessionName string // for connect | send | restart
	command     string // for send
)

var sessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sessions",
	Run: func(cmd *cobra.Command, args []string) {
		for name, s := range session.Manager.Sessions {
			fmt.Printf("- %s (running=%v)\n", name, !s.Cmd.ProcessState.Exited())
		}
	},
}

var sessionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new session",
	RunE: func(cmd *cobra.Command, args []string) error {
		var p profile.Profile
		prof, exists := profile.Profiles[profileName]
		if exists {
			p = *prof
		} else {
			schedule := profile.Schedule{
				Timing: scheduleTiming,
				Script: scheduleScript,
			}
			autoRestartCfg := profile.AutoRestartCfg{
				Max:           autoRestartMax,
				TriggerCodes:  autoRestartTriggerCode,
				StartupScript: autoRestartStartupScript,
				Scheduled:     []profile.Schedule{schedule},
			}
			p = profile.Profile{
				Name:          profileName,
				KeepLogs:      false, // TODO: impl
				StartupScript: startupScript,
				AutoRestart:   autoRestartCfg,
			}
		}

		s, err := session.Manager.Create(p)
		if err != nil {
			return err
		}

		if err := s.Connect(); err != nil {
			return err
		}
		return nil
	},
}

var sessionConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a session",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := cmd.Flags().GetString("session")
		if err != nil {
			return err
		}

		s, exists := session.Manager.Sessions[name]
		if !exists {
			return fmt.Errorf("no session are there what named '%s'", name)
		}

		return s.Connect()
	},
}

var sessionSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send command to a session",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := cmd.Flags().GetString("session")
		if err != nil {
			return err
		}

		s, exists := session.Manager.Sessions[name]
		if !exists {
			return fmt.Errorf("no session are there what named '%s'", name)
		}

		return s.Connect()
	},
}

var sessionRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a closed session",
	RunE: func(cmd *cobra.Command, args []string) error {
		//
		// TODO: impl here
		return nil
	},
}

func init() {
	Cmd.AddCommand(sessionListCmd)

	Cmd.AddCommand(sessionCreateCmd)
	sessionCreateCmd.Flags().StringVar(&name, "name", "", "name of a session")
	sessionCreateCmd.Flags().StringVar(&startupScript, "startup-script", "", "scripts what runs on startup")
	sessionCreateCmd.Flags().StringVar(&profileName, "profile", "", "name of will applied profile")
	sessionCreateCmd.Flags().BoolVar(&saveSession, "save-session", false, "save session or not")
	sessionCreateCmd.Flags().BoolVar(&autoRestart, "auto-restart", false, "auto restart session or not")
	sessionCreateCmd.Flags().UintVar(&autoRestartMax, "auto-restart-max", 0, "max restart times")
	sessionCreateCmd.Flags().IntSliceVar(&autoRestartTriggerCode, "auto-restart-trigger-code", []int{}, "trigger exit code array")
	sessionCreateCmd.Flags().StringVar(&autoRestartStartupScript, "auto-restart-startup-script", "", "auto restart startup script")
	sessionCreateCmd.Flags().StringVar(&scheduleTiming, "scheduled-script-timing", "", "timing to run script")
	sessionCreateCmd.Flags().StringVar(&scheduleScript, "scheduled-script", "", "script to run scheduled timing")
	if err := sessionCreateCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	Cmd.AddCommand(sessionConnectCmd)
	Cmd.AddCommand(sessionSendCmd)
	sessionSendCmd.Flags().StringVar(&command, "command", "", "command to run on session")
	if err := sessionSendCmd.MarkFlagRequired("command"); err != nil {
		panic(err)
	}
	Cmd.AddCommand(sessionRestartCmd)
	for _, c := range []*cobra.Command{sessionConnectCmd, sessionSendCmd, sessionRestartCmd} {
		c.Flags().StringVar(&sessionName, "session", "", "target session name")
		if err := c.MarkFlagRequired("session"); err != nil {
			panic(err)
		}
	}
}
