package dynamo

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Client struct {
	*dynamodb.Client
	CeremonyTableName string
	UserTableName     string
	// contains filtered or unexported fields
}

func NewClient(config aws.Config, userTableName string, ceremonyTableName string) *Client {
	return &Client{
		Client:            dynamodb.NewFromConfig(config),
		UserTableName:     userTableName,
		CeremonyTableName: ceremonyTableName,
	}
}

func (client *Client) MustUserTableName() *string {
	if client.UserTableName == "" {
		panic("user table name is empty")
	}
	return aws.String(client.UserTableName)
}

func (client *Client) MustCeremonyTableName() *string {
	if client.CeremonyTableName == "" {
		panic("ceremony table name is empty")
	}
	return aws.String(client.CeremonyTableName)
}

func IsConditionalCheckFailed(err error) bool {
	var bne *types.ConditionalCheckFailedException
	return errors.As(err, &bne)
}
