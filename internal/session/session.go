package session

import (
	"Iris/internal/profile"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type emptyWriter struct{}

func (w emptyWriter) Write(p []byte) (int, error) {
	return 0, nil
}

type Session struct {
	Name          string
	Profile       *profile.Profile
	Cmd           *exec.Cmd
	StdOutChannel chan string
	StdErrChannel chan string
	StdInChannel  chan string
	wg            sync.WaitGroup
	writeMtx      sync.RWMutex
	connections   map[*websocket.Conn]bool
	LogFile       string
	Options       SessionOptions
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

func Create(p profile.Profile, conn *websocket.Conn) error {
	args := strings.Fields(p.StartupScript)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	s := &Session{
		Name:          p.Name,
		Profile:       &p,
		Cmd:           cmd,
		StdOutChannel: make(chan string),
		StdErrChannel: make(chan string),
		StdInChannel:  make(chan string),
		writeMtx:      sync.RWMutex{},
		connections:   map[*websocket.Conn]bool{conn: true},
		LogFile:       "",
		Options: SessionOptions{
			LogFile:       "",
			KeepLogs:      &p.KeepLogs,
			StartupScript: &p.StartupScript,
			AutoRestart:   &p.AutoRestart,
			SaveSession:   false,
		},
	}

	Manager.Sessions[p.Name] = s

	s.wg.Add(3)
	go func() {
		defer func() {
			stdoutPipe.Close()
			s.wg.Done()
		}()
		stdoutScanner := bufio.NewScanner(stdoutPipe)
		for stdoutScanner.Scan() {
			s.StdOutChannel <- stdoutScanner.Text()
		}
		close(s.StdOutChannel)
	}()

	go func() {
		defer func() {
			stderrPipe.Close()
			s.wg.Done()
		}()
		stderrScanner := bufio.NewScanner(stderrPipe)
		for stderrScanner.Scan() {
			s.StdErrChannel <- stderrScanner.Text()
		}
		close(s.StdErrChannel)
	}()

	go func() {

		defer func() {
			stdinPipe.Close()
			s.wg.Done()
		}()
		//stdinScanner := bufio.NewWriter(stdinPipe)
		for mT, reader, err := conn.NextReader(); err != nil; {
			if mT != websocket.TextMessage {
				continue
			}
			var d []byte
			if _, err := reader.Read(d); err != nil {
				fmt.Println("websocket reader error " + err.Error())
				return
			}
			s.StdInChannel <- string(d)
		}
		close(s.StdInChannel)
	}()

	go func() {
		for {
			select {
			case o, ok := <-s.StdOutChannel:
				if !ok {
					s.StdOutChannel = nil
					continue
				}
				s.writeMtx.RLock()
				if s.connections != nil {
					// send line to connected sockets
					for c := range s.connections {
						if err := c.WriteMessage(websocket.TextMessage, []byte("[Out]"+o)); err != nil {
							fmt.Println(err)
						}
					}
				} else if s.connections == nil {
					if _, err := io.Discard.Write([]byte(o)); err != nil {
						fmt.Println(err)
					}
				}
				s.writeMtx.RUnlock()
			case e, ok := <-s.StdErrChannel:
				if !ok {
					s.StdErrChannel = nil
					continue
				}
				s.writeMtx.RLock()
				if s.connections != nil {
					for c := range s.connections {
						if err := c.WriteMessage(websocket.TextMessage, []byte("[Err]"+e)); err != nil {
							fmt.Println(err)
						}
					}
				} else if s.connections == nil {
					if _, err := io.Discard.Write([]byte(e)); err != nil {
						fmt.Println(err)
					}
				}
				s.writeMtx.RUnlock()

			case i, ok := <-s.StdInChannel:
				if !ok {
					s.StdInChannel = nil
					continue
				}
				s.writeMtx.RLock()
				if _, err := s.Cmd.Stdin.Read([]byte(i)); err != nil {
					fmt.Println(err)
				}
				s.writeMtx.RUnlock()
			}
			if s.StdOutChannel == nil && s.StdErrChannel == nil && s.StdInChannel == nil {
				delete(Manager.Sessions, s.Name)
				break
			}
		}
	}()

	if err := s.Cmd.Start(); err != nil {
		return err
	}
	return nil
}

func Connect(name string, conn *websocket.Conn) error {
	s, exists := Manager.Sessions[name]
	if !exists {
		return fmt.Errorf("no such session are there")
	}

	s.writeMtx.RLock()
	defer s.writeMtx.RUnlock()

	if s.connections == nil {
		s.connections = map[*websocket.Conn]bool{conn: true}
	} else if s.connections != nil {
		s.connections[conn] = true
	}
	return nil
}

func Detach(name string, conn *websocket.Conn) error {
	s, exists := Manager.Sessions[name]
	if !exists {
		return fmt.Errorf("no such session are there")
	}
	log.Println("Detaching session...")
	delete(s.connections, conn)
	if len(s.connections) == 0 {
		s.connections = nil
	}
	log.Println("Session detached.")
	return nil
}

func Send(name, command string, conn *websocket.Conn) error {
	s, exists := Manager.Sessions[name]
	if !exists {
		return fmt.Errorf("session not found: %s", name)
	}

	fmt.Printf("Sending command to %s: %s", name, command)

	if _, err := s.Cmd.Stdin.Read([]byte(command)); err != nil {
		return err
	}

	//if err := Connect(name, conn); err != nil {
	//	return err
	//}
	//
	//if err := conn.WriteMessage(websocket.TextMessage, []byte(command)); err != nil {
	//	return err
	//}
	//
	//if err := Detach(name, conn); err != nil {
	//	return err
	//}

	return nil
}
