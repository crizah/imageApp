package main

import (
	"log"
	"net/http"

	"fmt"
	"imageApp_2/server"
)

type MsgRequest struct {
	MessageID string `json:"messageID"`
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {

	s, err := server.InitialiseServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	http.HandleFunc("/upload", withCORS(s.UploadHandler))
	http.HandleFunc("/messages", withCORS(s.NotificationHandler))
	http.HandleFunc("/files", withCORS(s.FileHandler))
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
