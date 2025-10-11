package client

import (
	"Iris/core/handler"
	"Iris/core/structs"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

func Start(args []string) {
	if err := register(args); err != nil {
		log.Println("[client] failed to register process")
		return
	}
	time.Sleep(2 * time.Second) // wait for registration
	connect(args[0])
}

func register(args []string) error {
	uri := url.URL{Scheme: "http", Host: "localhost:8888", Path: handler.RegisterPath}
	pwd, _ := os.Getwd()
	connectRequest := structs.ConnectionRequest{
		SessionName: args[0],
		Env:         os.Environ(),
		Pwd:         pwd,
		Command:     strings.Join(args[1:], " "),
	}

	j, _ := json.Marshal(connectRequest)
	response, err := http.Post(uri.String(), "application/json", bytes.NewReader(j))
	if err != nil {
		log.Println("[client] http post error:", err)
		return err
	}

	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Println("[client] http post response status error:", response.StatusCode)
		body := make([]byte, response.ContentLength)
		response.Body.Read(body)
	}

	return nil
}

func connect(sessionName string) {
	uri := url.URL{Scheme: "ws", Host: "localhost:8888", Path: handler.ConnectPath}
	header := http.Header{}
	header.Add(handler.SessionNameHeader, sessionName)
	conn, _, err := websocket.DefaultDialer.Dial(uri.String(), header)
	if err != nil {
		log.Println("[client] websocket dial error:", err)
		return
	}

	defer func() {
		conn.Close()
	}()

	sig := make(chan os.Signal, 3)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	end := make(chan struct{}, 3)
	scan := make(chan bool, 1)

	connClosed := false
	sentEnd := false

	stdin := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer func() {
			if !sentEnd {
				end <- struct{}{}
				sentEnd = true
			}
			if !connClosed {
				if err := conn.Close(); err != nil {
					log.Println("[client] websocket connection close error:", err)
				}
				connClosed = true
			}
			wg.Done()
		}()

		select {
		case <-sig:
			return
		case <-end:
			return
		}
	}()

	go func() {
		defer func() {
			if !sentEnd {
				end <- struct{}{}
				sentEnd = true
			}
			wg.Done()
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("[client] websocket read error:", err)
				return
			}

			var post structs.Post
			if err := json.Unmarshal(msg, &post); err != nil {
				log.Println("[client] received data parse error:", err)
				return
			}

			if post.Type == structs.ProcessPost {
				fmt.Print(string(post.Line))
			}
		}
	}()

	go func() {
		defer func() {
			if !sentEnd {
				end <- struct{}{}
			}
			wg.Done()
		}()
		for {
			select {
			case <-sig:
				return
			case <-end:
				return
			case scan <- stdin.Scan():
				if <-scan {
					input := stdin.Bytes()
					userPost := structs.Post{
						Type:        structs.UserPost,
						SessionName: sessionName,
						Line:        input,
					}
					j, _ := json.Marshal(userPost)
					if err := conn.WriteMessage(websocket.TextMessage, j); err != nil {
						log.Println("[client] websocket write message error:", err)
						return
					}
				} else {
					// received EOF or else
					return
				}
			}
		}
	}()

	wg.Wait()
}
