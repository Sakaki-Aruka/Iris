package session

import (
	"Iris/internal/profile"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type Session struct {
	Name         string
	profile      *profile.Profile
	Cmd          *exec.Cmd
	LogFile      string
	Options      SessionOptions
	StdinWriter  io.WriteCloser
	StdoutReader io.ReadCloser
	StderrReader io.ReadCloser
	detachCh     chan bool
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
	cmd := exec.Command(p.StartupScript)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	s := &Session{
		Name:    p.Name,
		profile: &p,
		Cmd:     cmd,
		LogFile: "",
		Options: SessionOptions{
			LogFile:       "",
			KeepLogs:      &p.KeepLogs,
			StartupScript: &p.StartupScript,
			AutoRestart:   &p.AutoRestart,
			SaveSession:   false,
		},
		StdinWriter:  stdinPipe,
		StdoutReader: stdoutPipe,
		StderrReader: stderrPipe,
		detachCh:     make(chan bool),
	}

	if err := s.Cmd.Start(); err != nil {
		return nil, err
	}

	m.Sessions[p.Name] = s
	return s, nil
}

func (s *Session) Connect() error {
	go io.Copy(os.Stdout, s.StdoutReader)
	go io.Copy(os.Stderr, s.StderrReader)

	go func() {
		buffer := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buffer)
			if err != nil {
				if err == io.EOF {
					log.Println("Ctrl+D detected, detaching...")
					s.detachCh <- true
					return
				}
				log.Printf("Error reading from Stdin: %v", err)
				return
			}

			if n > 0 {
				if _, err := s.StdinWriter.Write(buffer); err != nil {
					log.Printf("Error writing to process Stdin: %v", err)
					return
				}
			}
		}
	}()

	<-s.detachCh
	s.Detach()

	fmt.Printf("Connected to session %s (PID=%d)", s.Name, s.Cmd.Process.Pid)
	return nil
}

func (s *Session) Detach() {
	log.Println("Detaching session...")
	if s.StdinWriter != nil {
		s.StdinWriter.Close()
	}

	if s.StdoutReader != nil {
		s.StdoutReader.Close()
	}

	if s.StderrReader != nil {
		s.StderrReader.Close()
	}

	log.Println("Session detached.")
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
