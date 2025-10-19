package handler

import (
	"Iris/core/structs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

func Connect(w http.ResponseWriter, r *http.Request) {
	sessionName := r.Header.Get(SessionNameHeader)
	if !ContainsProcessSocket(sessionName) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("'%v' not exists session name", sessionName)))
		return
	}

	conn, _ := upgrader.Upgrade(w, r, nil)
	PutConnected(sessionName, conn)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			DeleteConnected(sessionName, conn)
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("[connect] websocket error:", err)
				return
			}

			var post structs.Post
			if err := json.Unmarshal(msg, &post); err != nil {
				log.Println("[connect] post parse error:", err)
				continue
			}

			if post.Type == structs.UserPost {
				targetConn, err := GetProcessSocket(post.SessionName)
				if err == nil {
					if err := targetConn.WriteMessage(websocket.TextMessage, msg); err != nil {
						log.Println("[connect] user input write error:", err)
						continue
					}
				}
			}
		}
	}()

	wg.Wait()
}
