package server

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func (s *Server) Decryption(s3Key string, receiver string, encKey string) ([]byte, error) {

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

	// get receivers kms key from users table
	kmsKey, _, err := s.getRecipientKmsKey(receiver)

	if err != nil {
		return nil, err
	}

	// decrypt the dataKey with this kms key
	decKey, err := s.decryptKMS(kmsKey, encKey)
	if err != nil {
		return nil, err
	}
	dataKey := decKey.Plaintext

	// decrypt image from datakey
	decImage_bytes, err := decryptAES(enImage_bytes, dataKey)

	return decImage_bytes, err

}

func (s *Server) decryptKMS(kmsKey string, dataKey string) (*kms.DecryptOutput, error) {

	dataKey_bytes, err := base64.StdEncoding.DecodeString(dataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted data key: %w", err)
	}

	result, err := s.kmsClient.Decrypt(context.TODO(), &kms.DecryptInput{
		KeyId:          aws.String(kmsKey),
		CiphertextBlob: dataKey_bytes,
	})

	return result, err

}

func decryptAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract nonce from beginning
	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
