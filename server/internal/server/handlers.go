package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// Your signup endpoint
func (s *Server) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&req)
	res, err := s.CallCognito(req.Username, req.Password, req.Email)
	if err != nil {
		// Handle specific Cognito errors
		if strings.Contains(err.Error(), "UsernameExistsException") {
			http.Error(w, "Username or email already exists", http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "InvalidPasswordException") {
			http.Error(w, "Password does not meet requirements", http.StatusBadRequest)
			return
		}
		http.Error(w, "Signup failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Signup successful! Please check your email to verify your account.",
		"userSub": res.UserSub, // cognito user id
	})
}

func (s *Server) CallCognito(username string, password string, email string) (*cognitoidentityprovider.SignUpOutput, error) {

	resp, err := s.cognitoClient.SignUp(context.TODO(), &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(s.userPoolClientId),
		Username: aws.String(username), // Can be username or email
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(email)},
			{Name: aws.String("name"), Value: aws.String(username)},
		},
	})

	return resp, err

}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Username string `json:"username"` // can be email or username
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// authenticate with cognito
	authResult, err := s.cognitoClient.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(s.userPoolClientId),
		AuthParameters: map[string]string{
			"USERNAME": req.Username,
			"PASSWORD": req.Password,
		},
	})

	if err != nil {
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		if strings.Contains(err.Error(), "UserNotConfirmedException") {
			http.Error(w, "Please verify your email before logging in", http.StatusForbidden)
			return
		}
		http.Error(w, "Login failed: "+err.Error(), http.StatusInternalServerError)
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

	json.NewEncoder(w).Encode(map[string]interface{}{
		"msgs":  msgs,
		"count": count,
	})

}
func (s *Server) UserHandler(w http.ResponseWriter, r *http.Request) {
	// this is a get request

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

func (s *Server) VerificationHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err := s.cognitoClient.ConfirmSignUp(context.TODO(), &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(s.userPoolClientId),
		Username:         aws.String(req.Username),
		ConfirmationCode: aws.String(req.VerificationCode),
	})

	if err != nil {
		http.Error(w, "Verification failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Verification successful",
	})
}
