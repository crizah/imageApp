package server

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func putIntoMessagesTable(s string, r string, key string, f string, dk string, client *dynamodb.Client, msgID string) error {

	// messageId := uuid.New()
	timestamp := time.Now().Format(time.RFC3339)
	_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Messages"),
		Item: map[string]types.AttributeValue{
			"messageID":        &types.AttributeValueMemberS{Value: msgID},
			"sender":           &types.AttributeValueMemberS{Value: s},
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

func GetFromDynamo(receiver string, msgID string) (*dynamodb.GetItemOutput, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-north-1"))
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)

	// query messages to get msgID

	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
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

func QueryForMsgs(receiver string) ([]QueryResult, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-north-1"))
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)
	result, err := client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("Messages"),
		IndexName:              aws.String("recipient-index"),
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

func QueryForCount(receiver string) (int, error) {

	count := 0
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-north-1"))
	if err != nil {
		return count, err
	}

	client := dynamodb.NewFromConfig(cfg)
	result, err := client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("Messages"),
		IndexName:              aws.String("recipient-index"),
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
