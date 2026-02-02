package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type Server struct {
	dynamoClient     *dynamodb.Client
	s3Client         *s3.Client
	snsClient        *sns.Client
	kmsClient        *kms.Client
	cognitoClient    *cognitoidentityprovider.Client
	userPoolId       string
	userPoolClientId string
}

func InitialiseServer() (*Server, error) {
	awsRegion := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion))

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)
	snsclient := sns.NewFromConfig(cfg)
	s3client := s3.NewFromConfig(cfg)
	kmsclient := kms.NewFromConfig(cfg)
	cognitoClient := cognitoidentityprovider.NewFromConfig(cfg)
	userPoolId := os.Getenv("USER_POOL_ID")
	userPoolClientId := os.Getenv("USER_POOL_CLIENT_ID")

	server := &Server{
		dynamoClient:     client,
		snsClient:        snsclient,
		s3Client:         s3client,
		kmsClient:        kmsclient,
		cognitoClient:    cognitoClient,
		userPoolId:       userPoolId,
		userPoolClientId: userPoolClientId,
	}

	return server, nil

}

func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		EnableCors(w, r, origin)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// get token from cookie
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized - no token", http.StatusUnauthorized)
			return
		}

		// verify with cognito
		userInfo, err := s.cognitoClient.GetUser(context.TODO(), &cognitoidentityprovider.GetUserInput{
			AccessToken: aws.String(cookie.Value),
		})

		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// add username to context
		ctx := context.WithValue(r.Context(), "username", *userInfo.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (s *Server) CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	// get

	origin := r.Header.Get("Origin")
	EnableCors(w, r, origin)

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNoContent)
		return

	}

	cookie, err := r.Cookie("access_token")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	//verify token
	userInfo, err := s.cognitoClient.GetUser(context.TODO(), &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(cookie.Value),
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"username":      *userInfo.Username,
	})
}
