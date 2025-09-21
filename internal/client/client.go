package client

import (
	"Iris/daemon"
	"Iris/internal/profile"
	"Iris/ipc"
	"Iris/util"
	"fmt"
	"net"
	"net/rpc"
)

func newClient() (*rpc.Client, error) {
	socketPath, err := util.GetSocketPath()
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connecto to daemon: %w", err)
	}

	return rpc.NewClient(conn), nil
}

func call(method string, arg interface{}, reply interface{}) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()

	return client.Call(method, arg, reply)
}

func Create(p *profile.Profile) error {
	arg := &daemon.CreateArg{
		Profile: p,
	}
	reply := daemon.Reply{}
	err := call(ipc.ServiceName+".Create", arg, &reply)
	if reply.Err != nil && *reply.Err != "" {
		return fmt.Errorf(*reply.Err)
	} else {
		return err
	}
}

func Connect(sessionName string) error {
	arg := &daemon.SessionArg{
		SessionName: &sessionName,
	}
	reply := daemon.Reply{}
	err := call(ipc.ServiceName+".Connect", arg, &reply)
	if reply.Err != nil && *reply.Err != "" {
		return fmt.Errorf(*reply.Err)
	} else {
		return err
	}
}

func Detach(sessionName string) error {
	arg := &daemon.SessionArg{
		SessionName: &sessionName,
	}
	reply := daemon.Reply{}
	err := call(ipc.ServiceName+".Detach", arg, &reply)
	if reply.Err != nil && *reply.Err != "" {
		return fmt.Errorf(*reply.Err)
	} else {
		return err
	}
}

func Send(sessionName string, cmd string) error {
	arg := &daemon.SendArg{
		Name:    &sessionName,
		Command: &cmd,
	}
	reply := daemon.Reply{}
	err := call(ipc.ServiceName+".Send", arg, &reply)
	if reply.Err != nil && *reply.Err != "" {
		return fmt.Errorf(*reply.Err)
	} else {
		return err
	}
}
