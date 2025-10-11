package server

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/google/uuid"

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

func EncryptMsg(sender string, receiver string, fileName string, filePath string) error {

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
	err = putIntoMessagesTable(sender, receiver, s3Key, fileName, encryptedDK, client)
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
	return err

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

func putIntoMessagesTable(s string, r string, key string, f string, dk string, client *dynamodb.Client) error {

	messageId := uuid.New()
	timestamp := time.Now().Format(time.RFC3339)
	_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Messages"),
		Item: map[string]types.AttributeValue{
			"messageID":        &types.AttributeValueMemberS{Value: messageId.String()},
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
