package command

import (
	"Iris/core/handler"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

var ListCommand = &cobra.Command{
	Use:   "list",
	Short: "displays all sessions",
	Run: func(cmd *cobra.Command, _ []string) {
		log.Println("=== Iris sessions ===")
		uri := url.URL{Scheme: "http", Host: "localhost:8888", Path: handler.ListPath}
		response, err := http.Get(uri.String())
		if err != nil {
			log.Println("http get error:", err)
			return
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("response read error:", err)
			return
		}

		var result map[string]int
		if err := json.Unmarshal(body, &result); err != nil {
			log.Println("response parse error:", err)
			return
		}

		for k, v := range result {
			log.Println(fmt.Sprintf("SessionName: %v, Connected Clients: %v", k, v))
		}
	},
}
