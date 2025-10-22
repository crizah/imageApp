package main

import (
	"encoding/json"
	"fmt"
	"imageApp_2/server"
	"io"
	"net/http"
	"os"
)

var SENDER string
var RECEIVER string
var MSGID string

func getHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET , OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msgID := r.FormValue("messageID")
	username := r.FormValue("username")
	if username == "" || msgID == "" {
		http.Error(w, "got no username or msgID", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Notification sent for %s (%s)", username, msgID)

}

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
func notificationHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var receiver struct {
		Username string `json:"username"` // or whatever key you're sending
	}

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Receiver:", receiver.Username)

	// query messages db to get the unread messages and return count to send back to frontend

	count, err := server.QueryForCount(receiver.Username)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": count})

}

// func notificationHandler(w http.ResponseWriter, r *http.Request) {

// 	var req MsgRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	msgID := req.MessageID
// 	if msgID == "" {
// 		http.Error(w, "missing messageID", http.StatusBadRequest)
// 		return
// 	}
// 	// update this to also check for status

// 	// when user clicks on msgs
// 	// query all of receiver + unread to send notifications DYNAMICALLY

// 	result, err := server.GetFromDynamo(RECEIVER, msgID)

// 	if err != nil {
// 		http.Error(w, "couldnt get message"+err.Error(), http.StatusInternalServerError)
// 	}

// 	s3Key := result.Item["s3Key"].(*types.AttributeValueMemberS).Value
// 	encryptedDataKey := result.Item["encryptedDataKey"].(*types.AttributeValueMemberS).Value
// 	fileName := result.Item["fileName"].(*types.AttributeValueMemberS).Value
// 	sender := result.Item["sender"].(*types.AttributeValueMemberS).Value

// 	// get decrypted image
// 	decImage_bytes, err := server.Decryption(s3Key, RECEIVER, encryptedDataKey)
// 	if err != nil {
// 		http.Error(w, "decryption error"+err.Error(), http.StatusInternalServerError)
// 	}

// 	// change status in dynamo table to read

// 	// return image_data
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"sender":    sender,
// 		"fileName":  fileName,
// 		"imageData": base64.StdEncoding.EncodeToString(decImage_bytes),
// 	})

// }

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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
	SENDER = sender
	RECEIVER = recipient
	MSGID = msgID
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

	// works till here

	err = server.EncryptMsg(sender, recipient, header.Filename, filepath, msgID)
	if err != nil {
		fmt.Println(err)
	}

	// Delete local file after successful S3 upload
	err = os.Remove(filepath)
	if err != nil {
		fmt.Println("Warning: Could not delete local file:", err)
		// Don't fail the request if cleanup fails
	}

	// send notification to the receiver on the app itself
	// such that when they log into their account, and clicke messages, they can see who al sent them msgs

	// Return success response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Sent to receiver successfully",
	})
}

func main() {

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/messages", withCORS(notificationHandler))
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
