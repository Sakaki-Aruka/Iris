package handler

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	SessionNameHeader = "X-WS-SESSION-NAME"
)

const (
	ProcessPath  = "/process"
	RegisterPath = "/register"
	ConnectPath  = "/connect"
	ListPath     = "/list"
)

var cMu sync.Mutex
var pMu sync.Mutex

var connections = make(map[string]map[*websocket.Conn]bool)
var processes = make(map[string]*websocket.Conn)

func GetConnectionsMap() map[string]map[*websocket.Conn]bool {
	cMu.Lock()
	defer cMu.Unlock()
	return connections
}

func GetConnected(sessionName string) (map[*websocket.Conn]bool, error) {
	cMu.Lock()
	defer cMu.Unlock()

	m, exists := connections[sessionName]
	if !exists {
		return make(map[*websocket.Conn]bool), fmt.Errorf("connection not found")
	}
	return m, nil
}

func PutConnected(sessionName string, conn *websocket.Conn) {
	cMu.Lock()
	defer cMu.Unlock()

	if _, exists := connections[sessionName]; !exists {
		connections[sessionName] = make(map[*websocket.Conn]bool)
	}
	connections[sessionName][conn] = true
}

func DeleteConnected(sessionName string, conn *websocket.Conn) {
	cMu.Lock()
	defer cMu.Unlock()
	_, exists := connections[sessionName]
	if !exists {
		return
	}
	delete(connections[sessionName], conn)
}

func ContainsSession(sessionName string) bool {
	cMu.Lock()
	defer cMu.Unlock()
	_, exists := connections[sessionName]
	return exists
}

func DeleteSession(sessionName string) {
	cMu.Lock()
	defer cMu.Unlock()
	_, exists := connections[sessionName]
	if !exists {
		return
	}
	delete(connections, sessionName)
}

func DeleteProcessSocket(sessionName string) {
	pMu.Lock()
	defer pMu.Unlock()
	if _, exists := processes[sessionName]; !exists {
		return
	}
	delete(processes, sessionName)
}

func AddProcessSocket(sessionName string, conn *websocket.Conn) {
	pMu.Lock()
	defer pMu.Unlock()
	if _, exists := processes[sessionName]; exists {
		return
	}
	processes[sessionName] = conn
}

func GetProcessSocket(sessionName string) (*websocket.Conn, error) {
	pMu.Lock()
	defer pMu.Unlock()
	c, exists := processes[sessionName]
	if !exists {
		return nil, fmt.Errorf("'%v' not contained", sessionName)
	} else {
		return c, nil
	}
}

func ContainsProcessSocket(sessionName string) bool {
	pMu.Lock()
	defer pMu.Unlock()
	_, exists := processes[sessionName]
	return exists
}

func GetProcessMap() map[string]*websocket.Conn {
	pMu.Lock()
	defer pMu.Unlock()
	return processes
}
