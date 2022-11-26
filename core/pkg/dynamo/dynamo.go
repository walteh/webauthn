package dynamo

import (
	"context"
	"errors"
	"fmt"
	"nugg-webauthn/core/pkg/webauthn/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

type Puttable interface {
	MarshalDynamoDBAttributeValue() (*dtypes.AttributeValueMemberM, error)
	Put() (*dtypes.Put, error)
}

func MakePut(table *string, d Puttable) (*dtypes.Put, error) {
	put, err := d.Put()
	if err != nil {
		return nil, err
	}
	put.TableName = table
	return put, nil
}

func MakeDelete(table *string, d Gettable) (*dtypes.TransactWriteItem, error) {
	return &dtypes.TransactWriteItem{Delete: &dtypes.Delete{
		Key:       d.Get().Key,
		TableName: table,
	}}, nil
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

func (c *Client) TransactWrite(ctx context.Context, items ...dtypes.TransactWriteItem) error {
	_, err := c.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: items,
	})
	return err
}

type GetOutput struct {
	Item    map[string]dtypes.AttributeValue
	Request dtypes.TransactGetItem
}

type Gettable interface {
	UnmarshalDynamoDBAttributeValue(*dtypes.AttributeValueMemberM) error
	Get() *dtypes.Get
}

func (c *Client) BuildPut(d Puttable) (*dtypes.Put, error) {
	switch d.(type) {
	case *types.Ceremony:
		return MakePut(c.MustCeremonyTableName(), d)
	// case types.Cser:
	// 	return MakePut(c.MustUserTableName(), d)
	case *types.Credential:
		return MakePut(c.MustCredentialTableName(), d)
	default:
		return nil, fmt.Errorf("unknown type %T", d)
	}
}

func (c *Client) BuildDelete(d Gettable) (*dtypes.TransactWriteItem, error) {
	switch d.(type) {
	case *types.Ceremony:
		return MakeDelete(c.MustCeremonyTableName(), d)
	// case types.Cser:
	// 	return MakeDelete(c.MustUserTableName(), d)
	case *types.Credential:
		return MakeDelete(c.MustCredentialTableName(), d)
	default:
		return nil, fmt.Errorf("unknown type %T", d)
	}
}

func (c *Client) TransactGet(ctx context.Context, items ...Gettable) error {

	var itms []dtypes.TransactGetItem
	for _, item := range items {
		var tbl *string
		switch item.(type) {
		case *types.Ceremony:
			tbl = c.MustCeremonyTableName()
		// case types.Ctity:
		// 	tbl = c.MustUserTableName()
		case *types.Credential:
			tbl = c.MustCredentialTableName()
		default:
			return fmt.Errorf("unknown type %T", item)
		}

		this := item.Get()
		this.TableName = tbl

		itms = append(itms, dtypes.TransactGetItem{
			Get: this,
		})
	}

	res, err := c.TransactGetItems(ctx, &dynamodb.TransactGetItemsInput{
		TransactItems: itms,
	})

	if err != nil {
		return err
	}

	if len(res.Responses) != len(items) {
		return errors.New("unexpected number of responses")
	}

	for i := range items {

		err = items[i].UnmarshalDynamoDBAttributeValue(&dtypes.AttributeValueMemberM{Value: res.Responses[i].Item})
		if err != nil {
			return err
		}
	}

	return nil
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
	var bne *dtypes.ConditionalCheckFailedException
	return errors.As(err, &bne)
}
