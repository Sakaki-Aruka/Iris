package handler

import (
	"Iris/core/structs"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Process(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("websocket error:", err)
				return
			}

			var post structs.Post
			if err := json.Unmarshal(msg, &post); err != nil {
				log.Println("post-json parse error:", err)
				continue
			}

			if post.Type == structs.Process {
				c, err := GetConnected(post.SessionName)
				if err != nil {
					continue
				}
				for connected := range c {
					if err := connected.WriteMessage(websocket.TextMessage, msg); err != nil {
						DeleteConnected(post.SessionName, connected)
						connected.Close()
					}
				}
			} else if post.Type == structs.User {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Println("received write error:", err)
				}
			}
		}
	}()

	wg.Wait()
}
