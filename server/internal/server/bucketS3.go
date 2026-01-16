package server

import (
	"bytes"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// every user generates one key during sign up

const (
	BUCKET = "encrypted-files"
)

func (s *Server) UploadToS3(body []byte, key string) error {

	_, err := s.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})

	return err

}

func (s *Server) DeletFromS3(key string) error {

	_, err := s.s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(key),
	})

	return err

}

func (s *Server) GetfromS3(key string) (*s3.GetObjectOutput, error) {

	result, err := s.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(key),
	})

	return result, err

}
