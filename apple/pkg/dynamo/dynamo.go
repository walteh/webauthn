package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Client struct {
	*dynamodb.Client
	TableName string
	// contains filtered or unexported fields
}

func NewClient(config aws.Config, tableName string) *Client {
	return &Client{
		Client: dynamodb.NewFromConfig(config),
	}
}

type Challenge struct {
	Challenge string `json:"challenge"`
	UserId    string `json:"userId"`
	Ttl       int64  `json:"ttl"`
}
