package dynamo

import (
	"context"
	"errors"
	"time"

	"github.com/nuggxyz/webauthn/pkg/user"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoUser struct {
	UserId    string `dynamodbav:"user_id"`
	CreatedAt int64  `dynamodbav:"created_at"`
	UpdatedAt int64  `dynamodbav:"updated_at"`
}

func (client *Client) NewUserPut(id string) (*types.Put, error) {
	usr := &DynamoUser{
		UserId:    id,
		CreatedAt: (time.Now().Unix()),
		UpdatedAt: time.Now().Unix(),
	}
	av, err := attributevalue.MarshalMap(usr)
	if err != nil {
		return nil, err
	}

	return &types.Put{
		Item:      av,
		TableName: client.MustUserTableName(),
	}, nil
}

func (client *Client) LoadUser(ctx context.Context, userId string) (*user.User, error) {

	if client.UserTableName == "" {
		return nil, errors.New("user table name is not set")
	}

	var usr user.User

	input := &dynamodb.GetItemInput{
		TableName: client.MustUserTableName(),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userId},
		},
	}

	result, err := client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(result.Item, &usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (client *Client) FindUserInGetResult(result []*GetOutput) (cer *DynamoUser, err error) {

	usr, err := FindInOnePerTableGetResult[DynamoUser](result, client.MustUserTableName())
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// func (client *Client) SaveNewUser(ctx context.Context, usr *user.User, ceremony *Ceremony) error {

// 	item, err := attributevalue.MarshalMap(usr)
// 	if err != nil {
// 		log.Printf("failed to marshal challenge, %v", err)
// 		return err
// 	}

// 	cer, err := attributevalue.MarshalMap(ceremony)
// 	if err != nil {
// 		log.Printf("failed to marshal challenge, %v", err)
// 		return err
// 	}

// 	_, err = client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
// 		TransactItems: []types.TransactWriteItem{
// 			{
// 				Put: &types.Put{
// 					TableName:           client.MustUserTableName(),
// 					Item:                item,
// 					ConditionExpression: aws.String("attribute_not_exists(user_id) AND attribute_not_exists(apple_id) AND attribute_not_exists(username)"),
// 				},
// 			},
// 			{
// 				Put: &types.Put{
// 					TableName: client.MustCeremonyTableName(),
// 					Item:      cer,
// 				},
// 			},
// 		},
// 	})

// 	if err != nil {
// 		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
// 		return err
// 	}
// 	return nil
// }

// func (client *Client) SaveFirstUserLogin(ctx context.Context, usr *user.User, ceremony *Ceremony) error {

// 	item, err := attributevalue.MarshalMap(usr)
// 	if err != nil {
// 		log.Printf("failed to marshal challenge, %v", err)
// 		return err
// 	}

// 	cer, err := attributevalue.MarshalMap(ceremony)
// 	if err != nil {
// 		log.Printf("failed to marshal challenge, %v", err)
// 		return err
// 	}

// 	_, err = client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
// 		TransactItems: []types.TransactWriteItem{
// 			{
// 				Update: &types.Update{
// 					TableName: client.MustUserTableName(),
// 					Key:       item,
// 				},
// 			},
// 			{
// 				Put: &types.Put{
// 					TableName: client.MustCeremonyTableName(),
// 					Item:      cer,
// 				},
// 			},
// 		},
// 	})

// 	if err != nil {
// 		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
// 		return err
// 	}
// 	return nil
// }
