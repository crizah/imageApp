package server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
)

func (s *Server) Decryption(s3Key string, receiver string, encKey string) ([]byte, error) {

	s3Image, err := s.GetfromS3(s3Key)
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
	// kmsKey, _, err := s.getRecipientKmsKey(receiver)

	// decrypt the dataKey with this kms key
	dataKey, err := s.decryptKMS(encKey)
	if err != nil {
		return nil, err
	}

	// dataKey := decKey.Plaintext

	// decrypt image from datakey
	decImage_bytes, err := decryptAES(enImage_bytes, dataKey)

	return decImage_bytes, err

}

func (s *Server) decryptKMS(dataKey string) ([]byte, error) {

	dataKey_bytes, err := base64.StdEncoding.DecodeString(dataKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to decode encrypted data key: %w", err)
	// }

	// result, err := s.kmsClient.Decrypt(context.TODO(), &kms.DecryptInput{
	// 	KeyId:          aws.String(kmsKey),
	// 	CiphertextBlob: dataKey_bytes,
	// })

	return dataKey_bytes, err

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
