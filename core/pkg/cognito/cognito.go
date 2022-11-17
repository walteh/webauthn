package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

type Client struct {
	*cognitoidentity.Client

	PoolName string
}

func NewClient(config aws.Config, poolName string) *Client {
	return &Client{
		Client:   cognitoidentity.NewFromConfig(config),
		PoolName: poolName,
	}
}

func (c *Client) GetIdentityId(ctx context.Context, token string) (string, error) {

	resp, err := c.GetId(ctx, &cognitoidentity.GetIdInput{
		IdentityPoolId: aws.String(c.PoolName),
		Logins: map[string]string{
			"nuggid.nugg.xyz": token,
		},
	})

	if err != nil {
		return "", err
	}

	return *resp.IdentityId, nil
}

func (c *Client) GetCredentials(ctx context.Context, identityId string, token string) (aws.Credentials, error) {

	resp, err := c.GetCredentialsForIdentity(ctx, &cognitoidentity.GetCredentialsForIdentityInput{
		IdentityId: aws.String(identityId),
		Logins: map[string]string{
			"nuggid.nugg.xyz": token,
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

func (c *Client) GetDevCreds(ctx context.Context, nuggId string) (*cognitoidentity.GetOpenIdTokenForDeveloperIdentityOutput, error) {

	resp, err := c.GetOpenIdTokenForDeveloperIdentity(ctx, &cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
		Logins: map[string]string{
			"nuggid.nugg.xyz": nuggId,
		},
		IdentityPoolId: aws.String(c.PoolName),
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// validate a users open id token

// func (c *Client) ValidateToken(ctx context.Context, token string) (*cognitoidentity.ValidateIdentityInput, error) {

// 	resp, err := c.GetCredentialsForIdentity()(ctx, &cognitoidentity.ValidateIdentityInput{
// 		Token: aws.String(token),
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }
