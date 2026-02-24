package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4040"
	}
	verifyToken := os.Getenv("VERIFY_TOKEN")

	r := chi.NewRouter()

	// GET / - verification endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		mode := r.URL.Query().Get("hub.mode")
		challenge := r.URL.Query().Get("hub.challenge")
		token := r.URL.Query().Get("hub.verify_token")

		if mode == "subscribe" && token == verifyToken {
			log.Println("WEBHOOK VERIFIED")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(challenge)
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusForbidden)
	})

	// POST / - webhook receiver
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	log.Printf("\nListening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
