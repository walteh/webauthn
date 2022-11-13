package dynamo

import (
	"context"
	"log"
	"nugg-auth/core/pkg/random"
	"nugg-auth/core/pkg/signinwithapple"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type SessionInfo struct {
	AuthProvider string `dynamodbav:"auth_provider"`
	AccessToken  string `dynamodbav:"access_token"`
	Ttl          int64  `dynamodbav:"ttl"`
	CreatedAt    int64  `dynamodbav:"created_at"`
	UpdatedAt    int64  `dynamodbav:"updated_at"`
	CognitoId    string `dynamodbav:"cognito_id"`
}

type User struct {
	Id                string        `dynamodbav:"id"`
	AppleId           string        `dynamodbav:"apple_id"`
	AppleRefreshToken string        `dynamodbav:"apple_refresh_token"`
	Sessions          []SessionInfo `dynamodbav:"sessions"`
	CreatedAt         int64         `dynamodbav:"created_at"`
	UpdatedAt         int64         `dynamodbav:"updated_at"`
}

// AddMovie adds a movie the DynamoDB table.
func (client *Client) GenerateUser(ctx context.Context, appleId string, cognitoId string, abc *signinwithapple.ValidationResponse) error {

	now := time.Now().Unix()

	challenge := User{
		Id:                random.KSUID(),
		AppleId:           appleId,
		AppleRefreshToken: abc.RefreshToken,
		CreatedAt:         now,
		UpdatedAt:         now,
		Sessions: []SessionInfo{
			{
				AuthProvider: "apple/signinwithapple",
				AccessToken:  abc.AccessToken,
				Ttl:          (now) + int64(abc.ExpiresIn),
				CreatedAt:    now,
				UpdatedAt:    now,
				CognitoId:    cognitoId,
			},
		},
	}

	item, err := attributevalue.MarshalMap(challenge)
	if err != nil {
		log.Printf("failed to marshal challenge, %v", err)
		return err
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(client.TableName), Item: item,
		ConditionExpression: aws.String("attribute_not_exists(user_id) AND attribute_not_exists(apple_id)"),
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
		return err
	}
	return nil
}
