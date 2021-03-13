package ecr

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecr"
)

var (
	ErrEmptyAuthorisationData = errors.New("empty authorization data")
	ErrCredentialFormat       = errors.New("invalid credential format")
	ErrRegistryProxy          = errors.New("invalid registry proxy")
)

type Credential struct {
	Username    string
	Password    string
	RegistryURL string
}

type credentialsProvider interface {
	GetAuthorizationTokenWithContext(ctx aws.Context, input *ecr.GetAuthorizationTokenInput, opts ...request.Option) (*ecr.GetAuthorizationTokenOutput, error)
}

func GetAuthorizationTokens(ctx context.Context, provider credentialsProvider) ([]*ecr.AuthorizationData, error) {
	resp, err := provider.GetAuthorizationTokenWithContext(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return nil, err
	}

	if len(resp.AuthorizationData) == 0 {
		return nil, ErrEmptyAuthorisationData
	}

	return resp.AuthorizationData, nil
}

func DecodeCredentials(data []*ecr.AuthorizationData) ([]Credential, error) {
	result := make([]Credential, 0, len(data))

	for _, cred := range data {
		token, err := base64.StdEncoding.DecodeString(*cred.AuthorizationToken)
		if err != nil {
			return nil, err
		}

		fields := strings.Split(string(token), ":")
		if len(fields) != 2 {
			return nil, ErrCredentialFormat
		}

		if cred.ProxyEndpoint == nil {
			return nil, ErrRegistryProxy
		}

		result = append(result, Credential{
			Username:    fields[0],
			Password:    fields[1],
			RegistryURL: *cred.ProxyEndpoint,
		})
	}

	return result, nil
}
