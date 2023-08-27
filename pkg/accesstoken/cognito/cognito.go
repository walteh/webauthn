package cognito

import (
	"context"

	"github.com/walteh/webauthn/pkg/hex"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

type AWSCognitoClient interface {
	GetId(ctx context.Context, params *cognitoidentity.GetIdInput, optFns ...func(*cognitoidentity.Options)) (*cognitoidentity.GetIdOutput, error)
	GetCredentialsForIdentity(ctx context.Context, params *cognitoidentity.GetCredentialsForIdentityInput, optFns ...func(*cognitoidentity.Options)) (*cognitoidentity.GetCredentialsForIdentityOutput, error)
	GetOpenIdTokenForDeveloperIdentity(ctx context.Context, params *cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput, optFns ...func(*cognitoidentity.Options)) (*cognitoidentity.GetOpenIdTokenForDeveloperIdentityOutput, error)
}

type Client interface {
	GetDevCreds(ctx context.Context, nuggId hex.Hash) (*cognitoidentity.GetOpenIdTokenForDeveloperIdentityOutput, error)
	GetIdentityId(ctx context.Context, token string) (string, error)
	GetCredentials(ctx context.Context, identityId string, token string) (aws.Credentials, error)
}

type DefaultClient struct {
	*cognitoidentity.Client

	PoolName     string
	ProviderName string
}

func NewClient(config aws.Config, poolName string, providerName string) Client {
	return &DefaultClient{
		Client:       cognitoidentity.NewFromConfig(config),
		PoolName:     poolName,
		ProviderName: providerName,
	}
}

func (c *DefaultClient) GetIdentityId(ctx context.Context, token string) (string, error) {

	resp, err := c.GetId(ctx, &cognitoidentity.GetIdInput{
		IdentityPoolId: aws.String(c.PoolName),
		Logins: map[string]string{
			c.ProviderName: token,
		},
	})

	if err != nil {
		return "", err
	}

	return *resp.IdentityId, nil
}

func (c *DefaultClient) GetCredentials(ctx context.Context, identityId string, token string) (aws.Credentials, error) {

	resp, err := c.GetCredentialsForIdentity(ctx, &cognitoidentity.GetCredentialsForIdentityInput{
		IdentityId: aws.String(identityId),
		Logins: map[string]string{
			c.ProviderName: token,
		},
	})

	if err != nil {
		return aws.Credentials{}, err
	}

	r := credentials.NewStaticCredentialsProvider(*resp.Credentials.AccessKeyId, *resp.Credentials.SecretKey, *resp.Credentials.SessionToken)
	res, err := r.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	return res, nil
}

func (c *DefaultClient) GetDevCreds(ctx context.Context, nuggId hex.Hash) (*cognitoidentity.GetOpenIdTokenForDeveloperIdentityOutput, error) {

	resp, err := c.GetOpenIdTokenForDeveloperIdentity(ctx, &cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
		Logins: map[string]string{
			c.ProviderName: nuggId.Hex(),
		},
		IdentityPoolId: aws.String(c.PoolName),
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
