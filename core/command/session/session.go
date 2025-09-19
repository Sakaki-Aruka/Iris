package session

import (
	"Iris/core/command/profile"
	"os/exec"
)

type Session struct {
	ID      string
	Name    string
	profile *profile.Profile
	Cmd     *exec.Cmd
	LogFile string
	Options SessionOptions
}

type SessionOptions struct {
	LogFile       string
	KeepLogs      *bool
	StartupScript *profile.ScriptRef
	AutoRestart   *profile.AutoRestartCfg
	SaveSession   bool
}
