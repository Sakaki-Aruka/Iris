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
)

type CreateArg struct {
	Profile *profile.Profile
}

type CreateReply struct {
	Session *session.Session
	Err     *string
}

type SessionArg struct {
	Session *session.Session
}

type Reply struct {
	Err *string
}

type SendArg struct {
	Name    *string
	Command *string
}

type IrisService struct{}

func (i *IrisService) Create(arg *CreateArg, reply *CreateReply) error {
	if reply == nil {
		return fmt.Errorf("received pointer 'reply' points nil")
	}
	s, err := session.Manager.Create(*arg.Profile)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	*reply = CreateReply{
		Session: s,
		Err:     &errStr,
	}
	return err
}

func (i *IrisService) Connect(arg *SessionArg, reply *Reply) error {
	if reply == nil {
		return fmt.Errorf("received pointer 'reply' points nil")
	}
	err := arg.Session.Connect()
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	*reply = Reply{
		Err: &errStr,
	}

	return err
}

func (i *IrisService) Detach(arg *SessionArg, reply *struct{}) error {
	arg.Session.Detach()
	return nil
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

	if err := os.Remove(sock); err != nil {
		return err
	}

	listener, err := net.Listen("unix", sock)
	if err != nil {
		fmt.Println("Iris daemon boot error")
		return err
	}
	defer listener.Close()

	service := IrisService{}
	if err := rpc.RegisterName(ipc.ServiceName, service); err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return err
			}
			continue
		}
		go rpc.ServeConn(conn)
	}
}
