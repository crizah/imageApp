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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SendMessageRequest struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	FileName  string `json:"fileName"`
	ImageData string `json:"imageData"` // base64 encoded
}

func (s *Server) SendMsg(sender string, receiver string, fileName string, filePath string, msgID string) error {
	BUCKET := os.Getenv("BUCKET_NAME")

	// get the recipients kms key from dynamo Table

	snsTopicARN, err := s.getRecipientARN(receiver)
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
	encryptedImage, err := encryptAES(imageBytes, dataKey)
	if err != nil {
		return err
	}

	// 4. Encrypt the data key with recipient's KMS key

	encryptedDK, err := s.encryptDataKey(dataKey)
	if err != nil {
		return err
	}

	// upload encrypted image to s3
	s3Key := fmt.Sprintf("%s/%s/%s", sender, receiver, fileName)
	err = s.UploadToS3(encryptedImage, s3Key, BUCKET)
	if err != nil {

		return err
	}

	// store message metadata in DynamoDB
	err = s.putIntoMessagesTable(sender, receiver, s3Key, fileName, encryptedDK, msgID)
	if err != nil {
		// if error, delete from S3
		err = s.DeletFromS3(s3Key, BUCKET)
		if err != nil {
			return err
		}

		return err
	}

	// 7. Send SNS notification to recipient
	err = s.sendSNS(sender, snsTopicARN)

	// send notification to frontend of the receiver

	return err

}

func (s *Server) sendSNS(sender string, snsTopicARN string) error {

	// recipientTopic := fmt.Sprintf("arn:aws:sns:region:account:user-%s-notifications", receiver)
	// "arn:aws:sns:eu-north-1:YOUR_ACCOUNT_ID:user-%s-notifications"
	_, err := s.snsClient.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(snsTopicARN),
		Message:  aws.String(fmt.Sprintf("New message from %s", sender)),
		Subject:  aws.String("New Message"),
	})

	return err

}

func (s *Server) encryptDataKey(dataKey []byte) (string, error) {

	// encryptedDataKeyResult, err := s.kmsClient.Encrypt(context.TODO(), &kms.EncryptInput{
	// 	KeyId:     aws.String(recKey),
	// 	Plaintext: dataKey,
	// })

	// if err != nil {
	// 	return "", err
	// }

	encryptedDataKey := base64.StdEncoding.EncodeToString(dataKey)

	return encryptedDataKey, nil

}

func (s *Server) getRecipientARN(username string) (string, error) {

	//       TableName: "Users",
	//   Item: {
	//     username: { S: username },
	//     email: { S: email },
	//     kmsKeyId: { S: kmsKeyId },  // Fix: proper DynamoDB format
	//     snsTopicArn: { S: topicArn },
	//     createdAt: { S: new Date().toISOString() }
	//   }

	result, err := s.dynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: username},
		},
	})
	if err != nil {
		return "", err
	}

	// kmsKeyId := result.Item["kmsKeyId"].(*types.AttributeValueMemberS).Value
	snsTopicArn := result.Item["snsTopicArn"].(*types.AttributeValueMemberS).Value
	return snsTopicArn, nil
}

func encryptAES(data []byte, key []byte) ([]byte, error) {
	// Implement AES-256-GCM encryption
	// ... encryption logic ...

	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.

	var result []byte
	block, err := aes.NewCipher(key)
	if err != nil {
		return result, err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return result, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return result, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	res := append(nonce, ciphertext...)

	return res, nil

}
