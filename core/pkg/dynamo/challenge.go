package dynamo

import (
	"context"
	"encoding/json"
	"log"
	"nugg-auth/core/pkg/cwebauthn"
	"nugg-auth/core/pkg/webauthn/webauthn"

	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Challenge struct {
	Id          string `dynamodbav:"id"`
	SessionData []byte `dynamodbav:"session_data"`
	Ttl         int64  `dynamodbav:"ttl"`
	UserData    []byte `dynamodbav:"user_data"`
}

func clientDataToSafeID(clientData string, expectedChallengeType string, expectedOrigin string) (string, error) {
	var cd ClientData
	err := json.Unmarshal([]byte(clientData), &cd)
	if err != nil {
		log.Printf("failed to unmarshal client data, %v", err)
		return "", err
	}

	if cd.Type != expectedChallengeType {
		log.Printf("unexpected challenge type, %v", cd.Type)
		return "", ErrUnexpectedChallengeType
	}

	if cd.Origin != expectedOrigin {
		log.Printf("unexpected origin, %v", cd.Origin)
		return "", ErrUnexpectedOrigin
	}

	// id, err := safeid.ParseFromChallengeString(cd.Challenge)
	// if err != nil {
	// 	log.Printf("failed to parse challenge [%s], %v", cd.Challenge, err)
	// 	return nil, err
	// }

	return cd.Challenge, nil
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
		Id:          session.Challenge,
		SessionData: cer,
		Ttl:         (time.Now().Unix()) + 300,
		UserData:    jsonUserData,
	}

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

type ClientData struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
	Origin    string `json:"origin"`
}

// / load challenge
func (client *Client) LoadChallenge(ctx context.Context, clientData string, expectedChallengeType string, expectedOrigin string) (*webauthn.SessionData, *User, *cwebauthn.User, error) {

	key, err := clientDataToSafeID(clientData, expectedChallengeType, expectedOrigin)
	if err != nil {
		log.Printf("failed to get key from client data, %v", err)
		return nil, nil, nil, err
	}

	var challenge Challenge
	input := &dynamodb.GetItemInput{
		TableName: aws.String(client.TableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: key},
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

	webauthnUser := cwebauthn.NewUser([]byte(user.AppleId), user.Username)

	return &session, &user, webauthnUser, nil
}
