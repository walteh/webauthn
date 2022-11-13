package dynamo

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Client struct {
	*dynamodb.Client
	TableName string
	// contains filtered or unexported fields
}

func NewClient(config aws.Config, tableName string) *Client {
	return &Client{
		Client:    dynamodb.NewFromConfig(config),
		TableName: tableName,
	}
}

func IsConditionalCheckFailed(err error) bool {
	var bne *types.ConditionalCheckFailedException
	return errors.As(err, &bne)
}
