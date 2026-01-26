package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Read verification token from environment (same as your Node.js example)
	verifyToken := os.Getenv("VERIFY_TOKEN")
	if verifyToken == "" {
		log.Fatal("VERIFY_TOKEN environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGet(w, r, verifyToken)

		case http.MethodPost:
			handlePost(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleGet(w http.ResponseWriter, r *http.Request, verifyToken string) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode == "subscribe" && token == verifyToken {
		log.Println("WEBHOOK VERIFIED")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(challenge))
		return
	}

	w.WriteHeader(http.StatusForbidden)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("\n\nWebhook received %s\n", timestamp)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Pretty print JSON
	var prettyJSON interface{}
	if err := json.Unmarshal(body, &prettyJSON); err != nil {
		fmt.Println("Received non-JSON body:")
		fmt.Println(string(body))
	} else {
		pretty, _ := json.MarshalIndent(prettyJSON, "", "  ")
		fmt.Println(string(pretty))
	}

	// You can also log headers if needed
	// log.Printf("Headers: %+v\n", r.Header)

	// Always respond with 200 OK quickly (important for WhatsApp/Meta webhooks)
	w.WriteHeader(http.StatusOK)
}
