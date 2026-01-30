package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// var allowedOrigin = map[string]bool{
// 	os.Getenv("CLIENT_IP"):  true,
// 	"http://localhost:3000": true,
// 	"http://localhost:5173": true,
// }

func EnableCors(w http.ResponseWriter, r *http.Request, origin string) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

}
func (s *Server) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	EnableCors(w, r, origin)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&req)
	res, err := s.CallCognito(req.Username, req.Password, req.Email)
	if err != nil {

		http.Error(w, "error "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Signup successful! Please check your email to verify your account.",
		"userSub": res.UserSub, // cognito user id
	})
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	EnableCors(w, r, origin)

	_, err := r.Cookie("access_token")
	if err == nil {
		http.Error(w, "Unauthorized - no token", http.StatusUnauthorized)
		return

	}
	sec := false
	if os.Getenv("SECURE") == "true" {
		sec = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "id_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   sec,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   sec,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   sec,
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Logged out successfully",
	})
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	EnableCors(w, r, origin)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"` // can be email or username
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// authenticate with cognito
	authResult, err := s.AuthenticateCognito(req.Username, req.Password)
	if err != nil {
		http.Error(w, "error authenticating with cognito"+err.Error(), http.StatusInternalServerError)
		return
	}

	// get tokens
	idToken := *authResult.AuthenticationResult.IdToken
	accessToken := *authResult.AuthenticationResult.AccessToken
	refreshToken := *authResult.AuthenticationResult.RefreshToken
	sec := false
	if os.Getenv("SECURE") == "true" {
		sec = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "id_token",
		Value:    idToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   sec,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   sec,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   sec,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   2592000, // 30 days
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Login successful",
		"username": req.Username,
	})
}

func (s *Server) NotificationHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	EnableCors(w, r, origin)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var receiver struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// query messages db to get the unread messages and return count to send back to frontend

	count, err := s.QueryForCount(receiver.Username)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	msgs, err := s.QueryForMsgs(receiver.Username)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"msgs":  msgs,
		"count": count,
	})

}
func (s *Server) UserHandler(w http.ResponseWriter, r *http.Request) {
	// this is a get request
	origin := r.Header.Get("Origin")
	EnableCors(w, r, origin)

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}
	// query the Users table to get all usernames
	usernames, err := s.QueryForUsernames()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"usernames": usernames,
	})

}

func (s *Server) FileHandler(w http.ResponseWriter, r *http.Request) {

	origin := r.Header.Get("Origin")
	EnableCors(w, r, origin)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type Results struct {
		MessageID   string `json:"messageID"`
		Receiver    string `json:"receiver"`
		EncryptedDK string `json:"encryptedDK"`
		S3Key       string `json:"s3Key"`
		FileName    string `json:"fileName"`
	}

	var receiver struct {
		Msgs []Results `json:"msgs"`
	}

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var base64EncodedFiles []string

	for _, msg := range receiver.Msgs {
		b1, err := s.Decryption(msg.S3Key, msg.Receiver, msg.EncryptedDK)
		if err != nil {
			http.Error(w, "Decryption failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// encode the decrypted bytes as base64
		encoded := base64.StdEncoding.EncodeToString(b1)
		base64EncodedFiles = append(base64EncodedFiles, encoded)
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]string{
		"files": base64EncodedFiles,
	})
}

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {

	origin := r.Header.Get("Origin")
	EnableCors(w, r, origin)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	recipient := r.FormValue("recipient")
	sender := r.FormValue("sender")
	msgID := r.FormValue("msgID")

	fmt.Printf("Received file: %s for recipient: %s by sender: %s\n", header.Filename, recipient, sender)

	filepath := "./uploads/" + header.Filename
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(dst, file)
	dst.Close()

	if err != nil {
		os.Remove(filepath)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	err = s.SendMsg(sender, recipient, header.Filename, filepath, msgID)
	if err != nil {
		http.Error(w, "couldnt send msg"+err.Error(), http.StatusInternalServerError)
		return
	}

	// del local file after  S3 upload
	err = os.Remove(filepath)
	if err != nil {
		fmt.Println("Warning: Could not delete local file:", err)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Sent to receiver successfully",
	})
}

func (s *Server) VerificationHandler(w http.ResponseWriter, r *http.Request) {

	origin := r.Header.Get("Origin")
	EnableCors(w, r, origin)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get token
	// verify with cognito
	var req struct {
		Username         string `json:"username"`
		VerificationCode string `json:"verificationCode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	v, err := s.VerifyCognito(req.Username, req.VerificationCode)
	if err != nil && !v {
		http.Error(w, "user not verified"+err.Error(), http.StatusForbidden)
		return

	}

	if !v {
		http.Error(w, "user not verified", http.StatusForbidden)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Verification successful",
	})
}
