package server

// import (
// 	"context"

// 	"encoding/base64"

// 	"github.com/aws/aws-sdk-go-v2/config"

// 	"encoding/json"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

// 	"net/http"
// )

// func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	username := r.URL.Query().Get("username")
// 	if username == "" {
// 		http.Error(w, "Missing username parameter", http.StatusBadRequest)
// 		return
// 	}

// 	cfg, err := config.LoadDefaultConfig(context.TODO(),
// 		config.WithRegion("eu-north-1"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	client := dynamodb.NewFromConfig(cfg)

// 	// Query messages where user is the recipient
// 	result, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
// 		TableName:        aws.String("Messages"),
// 		FilterExpression: aws.String("recipient = :username"),
// 		ExpressionAttributeValues: map[string]types.AttributeValue{
// 			":username": &types.AttributeValueMemberS{Value: username},
// 		},
// 	})

// 	if err != nil {
// 		http.Error(w, "Failed to get messages: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Convert to JSON-friendly format
// 	messages := []map[string]string{}
// 	for _, item := range result.Items {
// 		msg := map[string]string{
// 			"messageID": item["messageID"].(*types.AttributeValueMemberS).Value,
// 			"sender":    item["sender"].(*types.AttributeValueMemberS).Value,
// 			"fileName":  item["fileName"].(*types.AttributeValueMemberS).Value,
// 			"status":    item["status"].(*types.AttributeValueMemberS).Value,
// 		}
// 		if timestamp, ok := item["timestamp"]; ok && timestamp != nil {
// 			msg["timestamp"] = timestamp.(*types.AttributeValueMemberS).Value
// 		}
// 		messages = append(messages, msg)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"messages": messages,
// 	})
// }

// Get and decrypt a specific message
// func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	messageID := r.URL.Query().Get("messageID")
// 	username := r.URL.Query().Get("username")

// 	if messageID == "" || username == "" {
// 		http.Error(w, "Missing messageID or username", http.StatusBadRequest)
// 		return
// 	}

// 	cfg, err := config.LoadDefaultConfig(context.TODO(),
// 		config.WithRegion("eu-north-1"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Get message metadata from DynamoDB
// 	dynamoClient := dynamodb.NewFromConfig(cfg)
// 	result, err := dynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
// 		TableName: aws.String("Messages"),
// 		Key: map[string]types.AttributeValue{
// 			"messageID": &types.AttributeValueMemberS{Value: messageID},
// 		},
// 	})

// 	if err != nil {
// 		http.Error(w, "Failed to get message: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if result.Item == nil {
// 		http.Error(w, "Message not found", http.StatusNotFound)
// 		return
// 	}

// 	// Verify user is the recipient
// 	recipient := result.Item["recipient"].(*types.AttributeValueMemberS).Value
// 	if recipient != username {
// 		http.Error(w, "Unauthorized", http.StatusForbidden)
// 		return
// 	}

// 	sender := result.Item["sender"].(*types.AttributeValueMemberS).Value
// 	s3Key := result.Item["s3Key"].(*types.AttributeValueMemberS).Value
// 	encryptedDataKey := result.Item["encryptedDataKey"].(*types.AttributeValueMemberS).Value
// 	fileName := result.Item["fileName"].(*types.AttributeValueMemberS).Value

// 	// Decrypt and return image
// 	imageData, err := DecryptMsg(username, s3Key, encryptedDataKey, cfg)
// 	if err != nil {
// 		http.Error(w, "Failed to decrypt message: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Mark as read
// 	dynamoClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
// 		TableName: aws.String("Messages"),
// 		Key: map[string]types.AttributeValue{
// 			"messageID": &types.AttributeValueMemberS{Value: messageID},
// 		},
// 		UpdateExpression: aws.String("SET #status = :read"),
// 		ExpressionAttributeNames: map[string]string{
// 			"#status": "status",
// 		},
// 		ExpressionAttributeValues: map[string]types.AttributeValue{
// 			":read": &types.AttributeValueMemberS{Value: "read"},
// 		},
// 	})

// 	// Return image data
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"sender":    sender,
// 		"fileName":  fileName,
// 		"imageData": base64.StdEncoding.EncodeToString(imageData),
// 	})
// }
