package server

import (
	"bytes"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// every user generates one key during sign up

const (
	BUCKET = "encrypted-files"
)

func UploadToS3(body []byte, key string) error {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-east-1"))
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})

	return err

}

func DeletFromS3(key string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-east-1"))
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(key),
	})

	return err

}

func GetfromS3(key string) (*s3.GetObjectOutput, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-east-1"))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(key),
	})

	return result, err

}
