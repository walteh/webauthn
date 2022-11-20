package cognito

import (
	"context"
	"nugg-auth/core/pkg/hex"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/smithy-go/middleware"
)

type MockClient struct {
}

func NewMockClient() Client {

	return &MockClient{}
}

func (c *MockClient) GetIdentityId(ctx context.Context, token string) (string, error) {
	return "", nil
}

func (c *MockClient) GetCredentials(ctx context.Context, identityId string, token string) (aws.Credentials, error) {
	return aws.Credentials{}, nil
}

func (c *MockClient) GetDevCreds(ctx context.Context, nuggId hex.Hash) (*cognitoidentity.GetOpenIdTokenForDeveloperIdentityOutput, error) {
	return &cognitoidentity.GetOpenIdTokenForDeveloperIdentityOutput{
		IdentityId:     aws.String("local:identity"),
		Token:          aws.String("OpenIdToken"),
		ResultMetadata: middleware.Metadata{},
	}, nil
}
