package server

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

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
