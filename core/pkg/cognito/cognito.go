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
			"appleid.apple.com": token,
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
			"appleid.apple.com": token,
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
