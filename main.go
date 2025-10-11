package main

import (
	"encoding/json"
	"fmt"
	"imageApp_2/server"
	"io"
	"net/http"
	"os"
)

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

	err = server.EncryptMsg(sender, recipient, header.Filename, filepath)
	if err != nil {
		fmt.Println(err)
	}

	// Delete local file after successful S3 upload
	err = os.Remove(filepath)
	if err != nil {
		fmt.Println("Warning: Could not delete local file:", err)
		// Don't fail the request if cleanup fails
	}

	// Return success response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Sent to receiver successfully",
	})
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
