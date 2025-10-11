package daemon

import (
	"Iris/core/handler"
	"net/http"
)

func Start() {
	http.HandleFunc(handler.ConnectPath, handler.Connect)
	http.HandleFunc(handler.RegisterPath, handler.Register)
	http.HandleFunc(handler.ProcessPath, handler.Process)
}
