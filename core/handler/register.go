package handler

import (
	"Iris/core/structs"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

func Register(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println("[register] request body close error:", err)
		}
	}()

	if r.Method != http.MethodPost {
		log.Println(fmt.Sprintf("[register] method '%v' is not allowed. (source: %v)", r.Method, r.RemoteAddr))
		return
	}

	buffer := new(bytes.Buffer)
	io.Copy(buffer, r.Body)
	body := buffer.Bytes()

	var request structs.ConnectionRequest
	if err := json.Unmarshal(body, &request); err != nil {
		log.Println("[register] create-request json parse error:", err)
		return
	}

	response, _ := json.Marshal(request)
	if _, err := w.Write(response); err != nil {
		log.Println("[register] response write error:", err)
	}
}

func startProcess(conReq structs.ConnectionRequest) error {
	commands := strings.Split(conReq.Command, " ")
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Env = conReq.Env
	cmd.Dir = conReq.Pwd
	buffer := make([]byte, 4096)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("[register] stdout get error:", err)
		return err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println("[register] stdin get error:", err)
		return err
	}

	end := make(chan struct{})

	endpoint := url.URL{Scheme: "ws", Host: "localhost:8888", Path: ProcessPath}
	header := http.Header{}
	header.Set(SessionNameHeader, conReq.SessionName)
	conn, _, err := websocket.DefaultDialer.Dial(endpoint.String(), header)
	if err != nil {
		log.Println("[register] dial error:", err)
		return err
	}

	go func() {
		defer func() {
			if cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
				cmd.Process.Kill()
			}
			cmd.Wait()
			conn.Close()

			cs, err := GetConnected(conReq.SessionName)
			if err == nil {
				for connected := range cs {
					connected.Close()
				}
				DeleteSession(conReq.SessionName)
				DeleteProcessSocket(conReq.SessionName)
			}
		}()

		AddProcessSocket(conReq.SessionName, conn)
		for {
			select {
			case <-end:
				return
			default:
				n, err := stdout.Read(buffer)
				if err != nil {
					log.Println("[register] stdout error:", err)
					return
				}

				if n > 0 {
					data := buffer[:n]
					post := structs.Post{
						Type:        structs.ProcessPost,
						SessionName: conReq.SessionName,
						Line:        data,
					}
					j, _ := json.Marshal(post)
					conn.WriteMessage(websocket.TextMessage, j)
				}
			}
		}
	}()

	go func() {
		defer func() {
			end <- struct{}{}
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("[register] process received message error:", err)
				return
			}

			var userPost structs.Post
			if err := json.Unmarshal(msg, &userPost); err != nil {
				log.Println("[register] process received-message parse error:", err)
				return
			}

			if userPost.Type == structs.UserPost {
				in := userPost.Line
				if !strings.HasSuffix(string(in), "\n") {
					in = append(in, "\n"...)
				}

				_, err := stdin.Write(in)
				if err != nil {
					log.Println("[register] process write message error:", err)
					return
				}
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Println("[register] process start error:", err)
		DeleteSession(conReq.SessionName)
		return err
	}

	return nil
}
