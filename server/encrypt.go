package server

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SendMessageRequest struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	FileName  string `json:"fileName"`
	ImageData string `json:"imageData"` // base64 encoded
}

func EncryptMsg(sender string, receiver string, fileName string, filePath string, msgID string) error {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-north-1"))
	if err != nil {
		return err
	}

	// get the recipients kms key from dynamo Table
	client := dynamodb.NewFromConfig(cfg)
	recKey, snsTopicARN, err := getRecipientKmsKey(client, receiver)
	if err != nil {
		return err
	}

	// generate ranbdim dek
	dataKey := make([]byte, 32) // 256-bit key
	rand.Read(dataKey)

	// 3. Encrypt the image with the data key (AES-256)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return err
	}

	imageBytes := buf.Bytes()

	// imageBytes, _ := base64.StdEncoding.DecodeString(buf)
	err, encryptedImage := encryptAES(imageBytes, dataKey)
	if err != nil {
		return err
	}

	// 4. Encrypt the data key with recipient's KMS key

	encryptedDK, err := encryptDataKey(cfg, recKey, dataKey)
	if err != nil {
		return err
	}

	// upload encrypted image to s3
	s3Key := fmt.Sprintf("%s/%s/%s", sender, receiver, fileName)
	err = UploadToS3(encryptedImage, s3Key)
	if err != nil {

		return err
	}

	// store message metadata in DynamoDB
	err = putIntoMessagesTable(sender, receiver, s3Key, fileName, encryptedDK, client, msgID)
	if err != nil {
		// if error, delete from S3
		err = DeletFromS3(s3Key)
		if err != nil {
			return err
		}

		return err
	}

	// 7. Send SNS notification to recipient
	err = sendSNS(cfg, sender, snsTopicARN)

	// send notification to frontend of the receiver

	return err

}

func Decryption(s3Key string, receiver string, encKey string) ([]byte, error) {

	s3Image, err := GetfromS3(s3Key)
	if err != nil {
		return nil, err
	}

	defer s3Image.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, s3Image.Body)
	if err != nil {
		return nil, err
	}
	enImage_bytes := buf.Bytes()

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-north-1"))
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)

	// get receivers kms key from users table
	kmsKey, _, err := getRecipientKmsKey(client, receiver)

	if err != nil {
		return nil, err
	}

	// decrypt the dataKey with this kms key
	decKey, err := decryptKMS(cfg, kmsKey, encKey)
	if err != nil {
		return nil, err
	}
	dataKey := decKey.Plaintext

	// decrypt image from datakey
	err, decImage_bytes := decryptAES(enImage_bytes, dataKey)

	return decImage_bytes, err

}

func sendSNS(cfg aws.Config, sender string, snsTopicARN string) error {
	snsClient := sns.NewFromConfig(cfg)
	// recipientTopic := fmt.Sprintf("arn:aws:sns:region:account:user-%s-notifications", receiver)
	// "arn:aws:sns:eu-north-1:YOUR_ACCOUNT_ID:user-%s-notifications"
	_, err := snsClient.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(snsTopicARN),
		Message:  aws.String(fmt.Sprintf("New message from %s", sender)),
		Subject:  aws.String("New Message"),
	})

	return err

}

func QueryForCount(receiver string) (int, error) {

	// error here
	// operation error DynamoDB: Query, https response error StatusCode: 400,
	// RequestID: 086KU10HEDTFK77H5L3JGT05T7VV4KQNSO5AEMVJF66Q9ASUAAJG,
	// api error ValidationException: Either the KeyConditions
	// or KeyConditionExpression parameter must be specified in the request.

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

func encryptDataKey(cfg aws.Config, recKey string, dataKey []byte) (string, error) {
	kmsClient := kms.NewFromConfig(cfg)
	encryptedDataKeyResult, err := kmsClient.Encrypt(context.TODO(), &kms.EncryptInput{
		KeyId:     aws.String(recKey),
		Plaintext: dataKey,
	})

	if err != nil {
		return "", err
	}

	encryptedDataKey := base64.StdEncoding.EncodeToString(encryptedDataKeyResult.CiphertextBlob)

	return encryptedDataKey, nil

}

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

func getRecipientKmsKey(client *dynamodb.Client, username string) (string, string, error) {

	//       TableName: "Users",
	//   Item: {
	//     username: { S: username },
	//     email: { S: email },
	//     kmsKeyId: { S: kmsKeyId },  // Fix: proper DynamoDB format
	//     snsTopicArn: { S: topicArn },
	//     createdAt: { S: new Date().toISOString() }
	//   }

	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: username},
		},
	})
	if err != nil {
		return "", "", err
	}

	kmsKeyId := result.Item["kmsKeyId"].(*types.AttributeValueMemberS).Value
	snsTopicArn := result.Item["snsTopicArn"].(*types.AttributeValueMemberS).Value
	return kmsKeyId, snsTopicArn, nil
}

func decryptKMS(cfg aws.Config, kmsKey string, dataKey string) (*kms.DecryptOutput, error) {
	kmsClient := kms.NewFromConfig(cfg)
	dataKey_bytes, err := base64.StdEncoding.DecodeString(dataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted data key: %w", err)
	}

	result, err := kmsClient.Decrypt(context.TODO(), &kms.DecryptInput{
		KeyId:          aws.String(kmsKey),
		CiphertextBlob: dataKey_bytes,
	})

	return result, err

}

func encryptAES(data []byte, key []byte) (error, []byte) {
	// Implement AES-256-GCM encryption
	// ... encryption logic ...

	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.

	var result []byte
	block, err := aes.NewCipher(key)
	if err != nil {
		return err, result
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err, result
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err, result
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	res := append(nonce, ciphertext...)

	return nil, res

}

// func DecryptMsg(username string, s3Key string, encryptedDataKeyB64 string, cfg aws.Config) ([]byte, error) {
// 	// Get user's KMS key
// 	dynamoClient := dynamodb.NewFromConfig(cfg)
// 	userKmsKey, _, err := getRecipientKmsKey(dynamoClient, username)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get user KMS key: %w", err)
// 	}

// 	// Decrypt the data key with user's KMS key
// 	kmsClient := kms.NewFromConfig(cfg)
// 	encryptedDataKeyBytes, err := base64.StdEncoding.DecodeString(encryptedDataKeyB64)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decode encrypted data key: %w", err)
// 	}

// 	decryptedKeyResult, err := kmsClient.Decrypt(context.TODO(), &kms.DecryptInput{
// 		KeyId:          aws.String(userKmsKey),
// 		CiphertextBlob: encryptedDataKeyBytes,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decrypt data key: %w", err)
// 	}
// 	dataKey := decryptedKeyResult.Plaintext

// 	// Download encrypted image from S3
// 	s3Client := s3.NewFromConfig(cfg)
// 	s3Result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: aws.String("non-encrypted-files"),
// 		Key:    aws.String(s3Key),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to download from S3: %w", err)
// 	}
// 	defer s3Result.Body.Close()

// 	var buf bytes.Buffer
// 	_, err = io.Copy(&buf, s3Result.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read S3 object: %w", err)
// 	}
// 	encryptedImageBytes := buf.Bytes()

// 	// Decrypt the image with the data key
// 	err, decryptedImage := decryptAES(encryptedImageBytes, dataKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decrypt image: %w", err)
// 	}

// 	return decryptedImage, nil
// }

func decryptAES(data []byte, key []byte) (error, []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err, nil
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err, nil
	}

	// Extract nonce from beginning
	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return fmt.Errorf("ciphertext too short"), nil
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err, nil
	}

	return nil, plaintext
}
