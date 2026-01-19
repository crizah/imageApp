package server

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (s *Server) CallCognito(username string, password string, email string) (*cognitoidentityprovider.SignUpOutput, error) {

	resp, err := s.cognitoClient.SignUp(context.TODO(), &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(s.userPoolClientId),
		Username: aws.String(username), // Can be username or email
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(email)},
			{Name: aws.String("name"), Value: aws.String(username)},
		},
	})

	return resp, err

}

func (s *Server) AuthenticateCognito(username string, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	authResult, err := s.cognitoClient.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(s.userPoolClientId),
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	})

	if err != nil {
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			return nil, err
		}
		if strings.Contains(err.Error(), "UserNotConfirmedException") {
			return nil, err
		}
		return nil, err
	}

	return authResult, nil

}

func (s *Server) VerifyCognito(username string, password string) (bool, error) {

	_, err := s.cognitoClient.ConfirmSignUp(context.TODO(), &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(s.userPoolClientId),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(password),
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
