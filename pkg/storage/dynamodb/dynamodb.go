package dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/walteh/webauthn/pkg/storage"
	"github.com/walteh/webauthn/pkg/webauthn/types"
)

var (
	_ storage.Provider = (*DynamoDBStorageClient)(nil)
)

type DynamoDBStorageClient struct {
	ceremonyTable   string
	credentialTable string
	client          *dynamodb.Client
}

func NewDynamoDBStorageClient(client aws.Config, ceremonyTable, credentialTable string) *DynamoDBStorageClient {
	return &DynamoDBStorageClient{
		ceremonyTable:   ceremonyTable,
		credentialTable: credentialTable,
		client:          dynamodb.NewFromConfig(client),
	}
}

func (me *DynamoDBStorageClient) GetExisting(ctx context.Context, crm types.CeremonyID, credid types.CredentialID) (*types.Ceremony, *types.Credential, error) {

	tgi := []dtypes.TransactGetItem{}

	tgi = append(tgi, dtypes.TransactGetItem{
		Get: &dtypes.Get{
			TableName: &me.ceremonyTable,
			Key: map[string]dtypes.AttributeValue{
				"id": &dtypes.AttributeValueMemberS{Value: string(crm.Ref().Hex())},
			},
		},
	})

	tgi = append(tgi, dtypes.TransactGetItem{
		Get: &dtypes.Get{
			TableName: &me.credentialTable,
			Key: map[string]dtypes.AttributeValue{
				"id": &dtypes.AttributeValueMemberS{Value: string(credid.Ref().Hex())},
			},
		},
	})

	resp, err := me.client.TransactGetItems(ctx, &dynamodb.TransactGetItemsInput{
		TransactItems:          tgi,
		ReturnConsumedCapacity: dtypes.ReturnConsumedCapacityNone,
	})
	if err != nil {
		return nil, nil, err
	}

	var credd *types.Credential
	var crmd *types.Ceremony

	for _, item := range resp.Responses {
		if str, ok := item.Item["id"].(*dtypes.AttributeValueMemberS); ok {
			if str.Value == string(crm) {
				err := crmd.UnmarshalDynamoDBAttributeValue(&dtypes.AttributeValueMemberM{Value: item.Item})
				if err != nil {
					return nil, nil, err
				}
			} else if str.Value == string(credid) {
				err := credd.UnmarshalDynamoDBAttributeValue(&dtypes.AttributeValueMemberM{Value: item.Item})
				if err != nil {
					return nil, nil, err
				}
			}
		}
	}

	return crmd, credd, nil
}

func (me *DynamoDBStorageClient) IncrementExistingCredential(ctx context.Context, crm types.CeremonyID, credid *types.Credential) error {

	twi := []dtypes.TransactWriteItem{}

	twi = append(twi, dtypes.TransactWriteItem{
		Update: &dtypes.Update{
			TableName: &me.credentialTable,
			Key: map[string]dtypes.AttributeValue{
				"id": &dtypes.AttributeValueMemberS{Value: string(credid.RawID.Ref().Hex())},
			},
			UpdateExpression: aws.String("ADD #count :inc"),
			ExpressionAttributeNames: map[string]string{
				"#count": "count",
			},
			ExpressionAttributeValues: map[string]dtypes.AttributeValue{
				":inc": &dtypes.AttributeValueMemberN{Value: "1"},
				"#sid": &dtypes.AttributeValueMemberS{Value: string(credid.SessionId.Hex())},
			},
			ConditionExpression: aws.String("attribute_exists(.id) AND .sessionid = #sid"),
		},
	})

	twi = append(twi, dtypes.TransactWriteItem{
		Delete: &dtypes.Delete{
			TableName: &me.ceremonyTable,
			Key: map[string]dtypes.AttributeValue{
				"id": &dtypes.AttributeValueMemberS{Value: string(crm.Ref().Hex())},
			},
			ConditionExpression: aws.String("attribute_exists(.id) AND .sessionid = #sid"),
			ExpressionAttributeValues: map[string]dtypes.AttributeValue{
				"#sid": &dtypes.AttributeValueMemberS{Value: credid.SessionId.Hex()},
			},
		},
	})

	_, err := me.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems:          twi,
		ReturnConsumedCapacity: dtypes.ReturnConsumedCapacityNone,
	})

	return err

}

func (me *DynamoDBStorageClient) WriteNewCeremony(ctx context.Context, crm *types.Ceremony) error {

	marshcer, err := crm.MarshalDynamoDBAttributeValue()
	if err != nil {
		return err
	}

	_, err = me.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &me.ceremonyTable,
		Item:                marshcer.Value,
		ConditionExpression: aws.String("attribute_not_exists(.id)"),
	})

	return err
}

// REQUIREMENT: the ceremony must exist in the db
// REQUIREMENT: the credential must not exist in the db
// REQUIREMENT: the sessionid of the credential must match the sessionid of the ceremony in the db
func (me *DynamoDBStorageClient) WriteNewCredential(ctx context.Context, crm types.CeremonyID, cred *types.Credential) error {

	twi := []dtypes.TransactWriteItem{}

	marshcred, err := cred.MarshalDynamoDBAttributeValue()
	if err != nil {
		return err
	}

	if cred.SessionId.IsZero() {
		return errors.New("sessionid is zero")
	}

	twi = append(twi, dtypes.TransactWriteItem{
		Delete: &dtypes.Delete{
			TableName: &me.ceremonyTable,
			Key: map[string]dtypes.AttributeValue{
				"id": &dtypes.AttributeValueMemberS{Value: string(crm)},
			},
			ConditionExpression: aws.String("attribute_exists(.id) AND .sessionid = #sid"),
			ExpressionAttributeValues: map[string]dtypes.AttributeValue{
				"#sid": &dtypes.AttributeValueMemberS{Value: cred.SessionId.Hex()},
			},
		},
	})

	twi = append(twi, dtypes.TransactWriteItem{
		Put: &dtypes.Put{
			TableName:           &me.credentialTable,
			Item:                marshcred.Value,
			ConditionExpression: aws.String("attribute_not_exists(.id)"),
		},
	})

	_, err = me.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems:          twi,
		ReturnConsumedCapacity: dtypes.ReturnConsumedCapacityNone,
	})

	return err
}
