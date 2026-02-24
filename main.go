package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4040"
	}
	verifyToken := os.Getenv("VERIFY_TOKEN")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mode := r.URL.Query().Get("hub.mode")
			challenge := r.URL.Query().Get("hub.challenge")
			token := r.URL.Query().Get("hub.verify_token")

			if mode == "subscribe" && token == verifyToken {
				log.Println("WEBHOOK VERIFIED")
		    w.Header().Set("Content-Type", "application/json")
    		w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(challenge)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}

		case http.MethodPost:
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			log.Printf("\n\nWebhook received %s\n", timestamp)

			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				log.Println("Error decoding JSON:", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			prettyBody, _ := json.MarshalIndent(body, "", "  ")
			log.Println(string(prettyBody))
			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Printf("\nListening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

