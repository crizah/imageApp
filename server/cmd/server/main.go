package main

import (
	"log"
	"net/http"
	"os"

	"fmt"
	"server/internal/server"

	"github.com/joho/godotenv"
)

// type MsgRequest struct {
// 	MessageID string `json:"messageID"`
// }

func main() {
	godotenv.Load()

	s, err := server.InitialiseServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	x := os.Getenv("WITH_INGRESS")

	mux := http.NewServeMux()

	mux.HandleFunc(fmt.Sprintf("%s/signup", x), s.SignUpHandler)        // works
	mux.HandleFunc(fmt.Sprintf("%s/login", x), s.LoginHandler)          // works
	mux.HandleFunc(fmt.Sprintf("%s/auth/check", x), s.CheckAuthHandler) // works
	mux.HandleFunc(fmt.Sprintf("%s/verify", x), s.VerificationHandler)  // works

	mux.HandleFunc(fmt.Sprintf("%s/usernames", x), s.AuthMiddleware(s.UserHandler))      // works
	mux.HandleFunc(fmt.Sprintf("%s/upload", x), s.AuthMiddleware(s.UploadHandler))       // works
	mux.HandleFunc(fmt.Sprintf("%s/notifs", x), s.AuthMiddleware(s.NotificationHandler)) // works
	mux.HandleFunc(fmt.Sprintf("%s/files", x), s.AuthMiddleware(s.FileHandler))          // works
	mux.HandleFunc(fmt.Sprintf("%s/logout", x), s.AuthMiddleware(s.LogoutHandler))       // DOES NOT WORK but i kinda dont care actually

	fmt.Println("server starting on :8082")
	http.ListenAndServe(":8082", mux)
}
