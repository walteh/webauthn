package dynamo

import (
	"context"
	"errors"
	"fmt"
	"nugg-auth/core/pkg/webauthn/protocol"

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

type Puttable interface {
	attributevalue.Marshaler
	Put() (*types.Put, error)
}

func MakePut(table *string, d Puttable) (*types.Put, error) {
	put, err := d.Put()
	if err != nil {
		return nil, err
	}
	put.TableName = table
	return put, nil
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

type Gettable interface {
	attributevalue.Unmarshaler
	Get() *types.Get
}

func (c *Client) BuildPut(d Puttable) (*types.Put, error) {
	switch d.(type) {
	case protocol.SavedCeremony:
		return MakePut(c.MustCeremonyTableName(), d)
	// case protocol.SavedUser:
	// 	return MakePut(c.MustUserTableName(), d)
	case protocol.SavedCredential:
		return MakePut(c.MustCredentialTableName(), d)
	default:
		return nil, fmt.Errorf("unknown type %T", d)
	}
}

func (c *Client) TransactGet(ctx context.Context, items ...Gettable) error {

	var itms []types.TransactGetItem
	for _, item := range items {
		var tbl *string
		switch item.(type) {
		case protocol.SavedCeremony:
			tbl = c.MustCeremonyTableName()
		// case protocol.UserEntity:
		// 	tbl = c.MustUserTableName()
		case protocol.SavedCredential:
			tbl = c.MustCredentialTableName()
		default:
			return errors.New("unknown table type")
		}

		this := item.Get()
		this.TableName = tbl

		itms = append(itms, types.TransactGetItem{
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
		err = items[i].UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberM{Value: res.Responses[i].Item})
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
	var bne *types.ConditionalCheckFailedException
	return errors.As(err, &bne)
}
