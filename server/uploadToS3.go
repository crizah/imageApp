package server

import (
	"bytes"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// every user generates one key during sign up

func UploadToS3(body []byte, key string) error {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-north-1"))
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)
	bucket := "non-encrypted-files"

	// file, err := os.Open(filePath)
	// if err != nil {
	// 	return err
	// }

	// defer file.Close()

	// var buf bytes.Buffer
	// if _, err := io.Copy(&buf, file); err != nil {

	// 	return err
	// }

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})

	return err

}

func DeletFromS3(key string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-north-1"))
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)
	bucket := "non-encrypted-files"

	_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return err

}
