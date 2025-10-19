package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func List(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println("[list] request body close error:", err)
		}
	}()

	if r.Method != http.MethodGet {
		log.Println(fmt.Sprintf("[list] method '%v' is not allowed. (source: %v)", r.Method, r.RemoteAddr))
		return
	}

	// return -> map[string]int
	// Key: SessionName, Value: Connected Client Amount

	var result = make(map[string]int)
	for k, _ := range GetProcessMap() {
		connected, err := GetConnected(k)
		if err != nil {
			result[k] = 0
		} else {
			result[k] = len(connected)
		}
	}
	
	response, _ := json.Marshal(result)
	if _, err := w.Write(response); err != nil {
		log.Println("[list] response write error:", err)
	}
}
