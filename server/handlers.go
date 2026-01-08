package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (s *Server) NotificationHandler(w http.ResponseWriter, r *http.Request) {

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

	fmt.Println("Receiver:", receiver.Username)

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

	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string][]server.QueryResult{"msgs": msgs})
	json.NewEncoder(w).Encode(map[string]interface{}{
		"msgs":  msgs,
		"count": count,
	})

}

func (s *Server) FileHandler(w http.ResponseWriter, r *http.Request) {
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

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	// Parse multipart form (10 MB max)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get file from request
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get recipient and sender
	recipient := r.FormValue("recipient")
	sender := r.FormValue("sender")
	msgID := r.FormValue("msgID")

	fmt.Printf("Received file: %s for recipient: %s by sender: %s\n", header.Filename, recipient, sender)

	// Save file temporarily
	filepath := "./uploads/" + header.Filename
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(dst, file)
	dst.Close() // Close immediately after writing

	if err != nil {
		os.Remove(filepath) // Clean up on error
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	err = s.SendMsg(sender, recipient, header.Filename, filepath, msgID)
	if err != nil {
		fmt.Println(err)
	}

	// Delete local file after successful S3 upload
	err = os.Remove(filepath)
	if err != nil {
		fmt.Println("Warning: Could not delete local file:", err)
		// Don't fail the request if cleanup fails
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Sent to receiver successfully",
	})
}
