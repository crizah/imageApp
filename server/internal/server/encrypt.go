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

	snsTopicARN, err := s.getRecipientARN(receiver)
	if err != nil {
		return err
	}

	// generate ranbdim dek
	dataKey := make([]byte, 32)
	rand.Read(dataKey)

	// encrypt w dk

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

	encryptedImage, err := encryptAES(imageBytes, dataKey)
	if err != nil {
		return err
	}

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
	// SNS notification to recipient
	err = s.sendSNS(sender, snsTopicARN)

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

	var result []byte
	block, err := aes.NewCipher(key)
	if err != nil {
		return result, err
	}

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
