package server

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Item={ // users table
//                 "username": {"S": username},
//                 # "userId": {"S": user_sub},  # Use Cognito sub as primary ID
//                 "email": {"S": email},
//                 "snsTopicArn": {"S": topic_arn},
//                 "emailVerified": {"BOOL": True},
//                 "createdAt": {"S": datetime.now().isoformat()},
//               // "updatedAt": {"S": datetime.utcnow().isoformat()}
//             }

func (s *Server) putIntoMessagesTable(str string, r string, key string, f string, dk string, msgID string) error {

	// messageId := uuid.New()
	timestamp := time.Now().Format(time.RFC3339)
	_, err := s.dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Messages"),
		Item: map[string]types.AttributeValue{
			"messageID":        &types.AttributeValueMemberS{Value: msgID},
			"sender":           &types.AttributeValueMemberS{Value: str},
			"recipient":        &types.AttributeValueMemberS{Value: r},
			"s3Key":            &types.AttributeValueMemberS{Value: key},
			"timestamp":        &types.AttributeValueMemberS{Value: timestamp},
			"encryptedDataKey": &types.AttributeValueMemberS{Value: dk},
			"fileName":         &types.AttributeValueMemberS{Value: f},
			"status":           &types.AttributeValueMemberS{Value: "unread"},
		},
	})

	return err

}

func (s *Server) GetFromDynamo(receiver string, msgID string) (*dynamodb.GetItemOutput, error) {

	// query messages to get msgID

	result, err := s.dynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Messages"),
		Key: map[string]types.AttributeValue{
			"messageID": &types.AttributeValueMemberS{Value: msgID},
		},
	})

	if err != nil {
		return nil, err
	}

	// verify if the recipient of the msg is the username provided

	assumed := result.Item["recipient"]
	sa := assumed.(*types.AttributeValueMemberS).Value

	if sa != receiver {
		return nil, errors.New("receiver provided didnt match with the database. unauthorized")

	}

	return result, nil

}

type QueryResult struct {
	MessageID   string
	Receiver    string
	Sender      string
	EncryptedDK string
	S3Key       string
	FileName    string
}

func (s *Server) QueryForMsgs(receiver string) ([]QueryResult, error) {

	result, err := s.dynamoClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("Messages"),
		IndexName:              aws.String("recipientIndex"),
		KeyConditionExpression: aws.String("recipient = :r"), // Only partition key here
		FilterExpression:       aws.String("#s = :s"),        // Status goes in filter
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":r": &types.AttributeValueMemberS{Value: receiver},
			":s": &types.AttributeValueMemberS{Value: "unread"},
		},
		ExpressionAttributeNames: map[string]string{
			"#s": "status",
		},
	})

	if err != nil {
		return nil, err

	}

	if int(result.Count) == 0 {
		return nil, errors.New("no unread messages")
	}

	var answer []QueryResult

	for _, item := range result.Items {
		ans := QueryResult{
			MessageID:   getString(item["messageID"]),
			Receiver:    receiver,
			Sender:      getString(item["sender"]),
			EncryptedDK: getString(item["encryptedDataKey"]),
			S3Key:       getString(item["s3Key"]),
			FileName:    getString(item["fileName"]),
		}
		answer = append(answer, ans)

	}

	return answer, nil
}

func getString(av types.AttributeValue) string {
	if av == nil {
		return ""
	}
	if v, ok := av.(*types.AttributeValueMemberS); ok {
		return v.Value
	}
	return ""
}

func (s *Server) QueryForCount(receiver string) (int, error) {

	count := 0

	result, err := s.dynamoClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("Messages"),
		IndexName:              aws.String("recipientIndex"),
		KeyConditionExpression: aws.String("recipient = :r"), // Only partition key here
		FilterExpression:       aws.String("#s = :s"),        // Status goes in filter
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":r": &types.AttributeValueMemberS{Value: receiver},
			":s": &types.AttributeValueMemberS{Value: "unread"},
		},
		ExpressionAttributeNames: map[string]string{
			"#s": "status",
		},
	})

	if err != nil {
		return count, err
	}

	count = int(result.Count)
	return count, nil

}

func (s *Server) QueryForUsernames() ([]string, error) {
	var usernames []string

	result, err := s.dynamoClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:            aws.String("Users"),
		ProjectionExpression: aws.String("username"),
	})

	if err != nil {
		return nil, err
	}

	for _, item := range result.Items {
		username := getString(item["username"])
		usernames = append(usernames, username)
	}

	return usernames, nil
}
