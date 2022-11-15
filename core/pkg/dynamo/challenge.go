package dynamo

import (
	"context"
	"encoding/json"
	"log"
	"nugg-auth/core/pkg/random"
	wan "nugg-auth/core/pkg/webauthn"

	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/duo-labs/webauthn/webauthn"
)

type Challenge struct {
	Id          string `dynamodbav:"id"`
	SessionData []byte `dynamodbav:"session_data"`
	Ttl         int64  `dynamodbav:"ttl"`
	UserData    []byte `dynamodbav:"user_data"`
}

// AddMovie adds a movie the DynamoDB table.
func (client *Client) SaveChallenge(ctx context.Context, session *webauthn.SessionData, userData *User) error {

	cer, err := json.Marshal(session)
	if err != nil {
		log.Printf("failed to marshal session, %v", err)
		return err
	}

	jsonUserData, err := json.Marshal(userData)
	if err != nil {
		log.Printf("failed to marshal user data, %v", err)
		return err
	}

	challenge := Challenge{
		Id:          random.KSUID().String(),
		SessionData: cer,
		Ttl:         (time.Now().Unix()) + 300,
		UserData:    jsonUserData,
	}

	log.Println("challenge", challenge)

	item, err := attributevalue.MarshalMap(challenge)
	if err != nil {
		log.Printf("failed to marshal challenge, %v", err)
		return err
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(client.TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
		return err
	}
	return nil
}

// / load challenge
func (client *Client) LoadChallenge(ctx context.Context, challengeId string) (*webauthn.SessionData, *User, *wan.User, error) {

	var challenge Challenge

	input := &dynamodb.GetItemInput{
		TableName: aws.String(client.TableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: challengeId},
		},
	}

	result, err := client.GetItem(ctx, input)
	if err != nil {
		log.Printf("Got error calling GetItem: %s	", err)
		return nil, nil, nil, err
	}

	if result.Item == nil {

		return nil, nil, nil, ErrNotFound
	}

	err = attributevalue.UnmarshalMap(result.Item, &challenge)
	if err != nil {
		log.Printf("Got error unmarshalling: %s", err)
		return nil, nil, nil, err
	}

	var session webauthn.SessionData
	err = json.Unmarshal(challenge.SessionData, &session)
	if err != nil {
		log.Printf("failed to unmarshal session, %v", err)
		return nil, nil, nil, err
	}

	var user User
	err = json.Unmarshal(challenge.UserData, &user)
	if err != nil {
		log.Printf("failed to unmarshal user, %v", err)
		return nil, nil, nil, err
	}

	webauthnUser := wan.NewUser([]byte(user.AppleId), user.Username)

	return &session, &user, webauthnUser, nil
}
