package daemon

import (
	"Iris/core/handler"
	"log"
	"net/http"
)

func Start() {
	http.HandleFunc(handler.ConnectPath, handler.Connect)
	http.HandleFunc(handler.RegisterPath, handler.Register)
	http.HandleFunc(handler.ProcessPath, handler.Process)
	http.HandleFunc(handler.ListPath, handler.List)
	if err := http.ListenAndServe("localhost:8888", nil); err != nil {
		log.Println("[daemon] listen and serve error:", err)
	}
}
