package session

import (
	"Iris/internal/profile"
	"fmt"
	"io"
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
	StartupScript *string
	AutoRestart   *profile.AutoRestartCfg
	SaveSession   bool
}

type SessionManager struct {
	Sessions map[string]*Session
}

var Manager = &SessionManager{
	Sessions: make(map[string]*Session),
}

func (m *SessionManager) Create(p profile.Profile) (*Session, error) {
	//
}

func (m *SessionManager) Connect(name string) error {
	s, exists := m.Sessions[name]
	if !exists {
		return fmt.Errorf("session not found: %s", name)
	}

	fmt.Printf("Connected to session %s (PID=%d)", name, s.Cmd.Process.Pid)
	return nil
}

func (m *SessionManager) Send(name, command string) error {
	s, exists := m.Sessions[name]
	if !exists {
		return fmt.Errorf("session not found: %s", name)
	}

	fmt.Printf("Sending command to %s: %s", name, command)

	stdin, err := s.Cmd.StdinPipe()
	if err != nil {
		return err
	}

	if _, writeErr := io.WriteString(stdin, command); writeErr != nil {
		if closeErr := stdin.Close(); closeErr != nil {
			return closeErr
		}
		return writeErr
	}

	if closeErr := stdin.Close(); closeErr != nil {
		return closeErr
	}
	return nil
}
