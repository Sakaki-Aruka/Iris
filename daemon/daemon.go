package daemon

import (
	"Iris/internal/profile"
	isession "Iris/internal/session"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Input struct {
	SessionName string `json:"session_name"`
	Action      Action `json:"action"`
	Line        string `json:"line"`
}

type Action int

const (
	SEND Action = iota
	CONNECT
	DETACH
	CREATE
)

type Response struct {
	Type int
	Line string
}

const (
	Ok  = 0
	Err = 1
)

func (r Response) GetLine() (string, error) {
	if r.Type != Ok {
		return "", fmt.Errorf("response type is not 'Ok'. (%v, %v)", r.Type, r.Line)
	}
	return r.Line, nil
}

var clients = make(map[string][]*websocket.Conn)

func session(writer http.ResponseWriter, reader *http.Request) {
	conn, err := upgrader.Upgrade(writer, reader, nil)
	if err != nil {
		fmt.Println("failed to upgrade connection to websocket from http")
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType != websocket.TextMessage {
			continue
		}

		var i Input
		if err := json.Unmarshal(message, &i); err != nil {
			continue
		}

		switch i.Action {
		case SEND:
			{
				if err := isession.Send(i.SessionName, i.Line, conn); err != nil {
					r := Response{Type: Err, Line: err.Error()}
					d, _ := json.Marshal(r)
					if err := conn.WriteMessage(websocket.TextMessage, d); err != nil {
						break
					}
				}
			}
		case CONNECT:
			{
				//TODO: impl
				if err := isession.Connect(i.SessionName, conn); err != nil {
					r := Response{Type: Err, Line: err.Error()}
					d, _ := json.Marshal(r)
					if err := conn.WriteMessage(websocket.TextMessage, d); err != nil {
						break
					}
				}
			}
		case DETACH:
			{
				//TODO: impl
				if err := isession.Detach(i.SessionName, conn); err != nil {
					r := Response{Type: Err, Line: err.Error()}
					d, _ := json.Marshal(r)
					if err := conn.WriteMessage(websocket.TextMessage, d); err != nil {
						break
					}
				}
			}

		case CREATE:
			{
				//TODO: impl
				var p profile.Profile
				if err := json.Unmarshal([]byte(i.Line), &p); err != nil {
					r := Response{Type: Err, Line: err.Error()}
					d, _ := json.Marshal(r)
					if err := conn.WriteMessage(websocket.TextMessage, d); err != nil {
						break
					}
				}
				if err := isession.Create(p, conn); err != nil {
					r := Response{Type: Err, Line: err.Error()}
					d, _ := json.Marshal(r)
					if err := conn.WriteMessage(websocket.TextMessage, d); err != nil {
						break
					}
				}
			}
		}
	}
}

func StartDaemon() error {
	//TODO: impl
	return nil
}
