package dynamo

import (
	"context"
	"encoding/json"
	"log"
	"nugg-auth/core/pkg/signinwithapple"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/duo-labs/webauthn/webauthn"
)

type SessionInfo struct {
	AuthProvider string `dynamodbav:"auth_provider" json:"auth_provider"`
	AccessToken  string `dynamodbav:"access_token"  json:"access_token`
	Ttl          int64  `dynamodbav:"ttl"           json:"ttl"`
	CreatedAt    int64  `dynamodbav:"created_at"    json:"created_at"`
	UpdatedAt    int64  `dynamodbav:"updated_at"    json:"updated_at"`
	CognitoId    string `dynamodbav:"cognito_id"    json:"cognito_id"`
}

type User struct {
	Id                string        `dynamodbav:"id" .                json:"id"`
	Username          string        `dynamodbav:"username"            json:"username"`
	AppleId           string        `dynamodbav:"apple_id"            json:"apple_id"`
	AppleRefreshToken string        `dynamodbav:"apple_refresh_token" json:"apple_refresh_token"`
	Sessions          []SessionInfo `dynamodbav:"sessions"            json:"sessions"`
	CreatedAt         int64         `dynamodbav:"created_at"          json:"created_at"`
	UpdatedAt         int64         `dynamodbav:"updated_at"          json:"updated_at"`
	SessionData       []byte        `dynamodbav:"session_data"        json:"session_data"`
}

func NewUser(newId string, username string, appleId string, cognitoId string, abc *signinwithapple.ValidationResponse) *User {
	now := time.Now().Unix()

	return &User{
		Id:                newId,
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
}

// AddMovie adds a movie the DynamoDB table.
func (client *Client) SaveNewUser(ctx context.Context, user *User, session *webauthn.SessionData) error {

	// convert session to bytes
	sessionBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	user.SessionData = sessionBytes

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		log.Printf("failed to marshal challenge, %v", err)
		return err
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(client.TableName), Item: item,
		ConditionExpression: aws.String("attribute_not_exists(user_id) AND attribute_not_exists(apple_id)"),
		// ConditionExpression: aws.String("attribute_not_exists(user_id) AND attribute_not_exists(apple_id) AND attribute_not_exists(username)"),
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
		return err
	}
	return nil
}
