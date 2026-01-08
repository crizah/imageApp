package server

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type Server struct {
	dynamoClient *dynamodb.Client
	s3Client     *s3.Client
	snsClient    *sns.Client
	kmsClient    *kms.Client
}

func InitialiseServer() (*Server, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-east-1"))

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)
	snsclient := sns.NewFromConfig(cfg)
	s3client := s3.NewFromConfig(cfg)
	kmsclient := kms.NewFromConfig(cfg)

	server := &Server{
		dynamoClient: client,
		snsClient:    snsclient,
		s3Client:     s3client,
		kmsClient:    kmsclient,
	}

	return server, nil

}
