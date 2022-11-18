package dynamo

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Client struct {
	*dynamodb.Client
	CeremonyTableName   string
	UserTableName       string
	CredentialTableName string
}

func NewClient(config aws.Config, userTableName string, ceremonyTableName string, credentialTableName string) *Client {
	return &Client{
		Client:              dynamodb.NewFromConfig(config),
		UserTableName:       userTableName,
		CeremonyTableName:   ceremonyTableName,
		CredentialTableName: credentialTableName,
	}
}

func FindInOnePerTableGetResult[I interface{}](result []*GetOutput, tableName *string) (*I, error) {
	if len(result) == 0 {
		return nil, ErrCeremonyNotFound
	}

	var ceremony I
	for _, item := range result {
		if *item.Request.Get.TableName != *tableName {
			continue
		}

		err := attributevalue.UnmarshalMap(item.Item, &ceremony)
		if err != nil {
			return nil, err
		}

		return &ceremony, nil

	}
	return nil, ErrCeremonyNotFound
}

func (c *Client) TransactWrite(ctx context.Context, items ...types.TransactWriteItem) error {
	_, err := c.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: items,
	})
	return err
}

type GetOutput struct {
	Item    map[string]types.AttributeValue
	Request types.TransactGetItem
}

func (c *Client) TransactGet(ctx context.Context, items ...types.TransactGetItem) ([]*GetOutput, error) {
	res, err := c.TransactGetItems(ctx, &dynamodb.TransactGetItemsInput{
		TransactItems: items,
	})

	if err != nil {
		return nil, err
	}

	if len(res.Responses) != len(items) {
		return nil, errors.New("unexpected number of responses")
	}

	output := make([]*GetOutput, len(res.Responses))
	for i, item := range items {
		output[i] = &GetOutput{
			Item:    res.Responses[i].Item,
			Request: item,
		}
	}
	return output, nil
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

func (client *Client) MustCredentialTableName() *string {
	if client.CeremonyTableName == "" {
		panic("ceremony table name is empty")
	}
	return aws.String(client.CredentialTableName)
}

func IsConditionalCheckFailed(err error) bool {
	var bne *types.ConditionalCheckFailedException
	return errors.As(err, &bne)
}
