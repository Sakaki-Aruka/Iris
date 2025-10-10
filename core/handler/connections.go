package handler

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var cMu sync.Mutex
var pMu sync.Mutex

var connections = make(map[string]map[*websocket.Conn]bool)
var processes = make(map[string]*websocket.Conn)

func GetConnected(sessionName string) (map[*websocket.Conn]bool, error) {
	cMu.Lock()
	defer cMu.Unlock()

	m, exists := connections[sessionName]
	if !exists {
		return make(map[*websocket.Conn]bool), fmt.Errorf("connection not found")
	}
	return m, nil
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

func DeleteSession(sessionName string) {
	cMu.Lock()
	defer cMu.Unlock()
	_, exists := connections[sessionName]
	if !exists {
		return
	}
	delete(connections, sessionName)
}
