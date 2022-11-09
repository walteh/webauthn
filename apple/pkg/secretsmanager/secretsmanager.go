package secretsmanager

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://aws.github.io/aws-sdk-go-v2/docs/getting-started/

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Client struct {
	*secretsmanager.Client
	secretName string
	ttl        time.Time
	cache      *secretsmanager.GetSecretValueOutput
}

func NewClient(ctx context.Context, config aws.Config, secretName string) (client *Client) {
	return &Client{secretsmanager.NewFromConfig(config), secretName, time.Now(), nil}
}

func (c *Client) Refresh(ctx context.Context) (secretString string, err error) {
	if c.cache != nil && c.ttl.Before(time.Now()) {
		return *c.cache.SecretString, nil
	}

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(c.secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := c.GetSecretValue(ctx, input)
	if err != nil {
		return "", err
	}

	c.cache = result
	c.ttl = time.Now().Add(time.Minute * 5)

	return *result.SecretString, nil
}
