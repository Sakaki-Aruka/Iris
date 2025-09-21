package daemon

import (
	"Iris/internal/profile"
	"Iris/internal/session"
	"Iris/ipc"
	"Iris/util"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type CreateArg struct {
	Profile *profile.Profile
}

type SessionArg struct {
	SessionName *string //*session.Session
}

type Reply struct {
	Err *string
}

type SendArg struct {
	Name    *string
	Command *string
}

type IrisService struct{}

func (i *IrisService) Create(arg *CreateArg, reply *Reply) error {
	if reply == nil {
		return fmt.Errorf("received pointer 'reply' points nil")
	}
	err := session.Manager.Create(*arg.Profile)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	*reply = Reply{
		Err: &errStr,
	}
	return err
}

func (i *IrisService) Connect(arg *SessionArg, reply *Reply) error {
	if reply == nil {
		return fmt.Errorf("received pointer 'reply' points nil")
	}

	err := session.Connect(*arg.SessionName)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	*reply = Reply{
		Err: &errStr,
	}

	return err
}

func (i *IrisService) Detach(arg *SessionArg, reply *Reply) error {
	if reply == nil {
		return fmt.Errorf("received pointer 'reply' points nil")
	}

	err := session.Detach(*arg.SessionName)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	*reply = Reply{
		Err: &errStr,
	}
	return err
}

func (i *IrisService) Send(arg *SendArg, reply *Reply) error {
	if reply == nil {
		return fmt.Errorf("received pointer 'reply' points nil")
	}
	err := session.Manager.Send(*arg.Name, *arg.Command)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	*reply = Reply{
		Err: &errStr,
	}
	return err
}

func StartDaemon() error {
	sock, err := util.GetSocketPath()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(sock); err != nil {
		return err
	}

	pid, err := util.GetDaemonPidPath()
	if err != nil {
		return err
	} else if _, err := os.Stat(pid); err == nil {
		fmt.Println("Iris pid file has already exists. The daemon is running or try to read an old pid file.")
		return err
	}

	if pidWriteErr := os.WriteFile(pid, []byte(strconv.Itoa(os.Getpid())), 0770); pidWriteErr != nil {
		return pidWriteErr
	}
	defer func() {
		if err := os.Remove(pid); err != nil {
			fmt.Println("failed to remove pid file. \n" + err.Error())
		}
	}()

	listener, err := net.Listen("unix", sock)
	if err != nil {
		fmt.Println("Iris daemon boot error")
		return err
	}
	defer listener.Close()

	service := IrisService{}
	if err := rpc.RegisterName(ipc.ServiceName, &service); err != nil {
		return err
	}
	defer listener.Close()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-signalCh
		listener.Close()
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			continue
		}
		go rpc.ServeConn(conn)
	}
	return nil
}
