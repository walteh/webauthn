package dynamo

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"log"
	"nugg-auth/core/pkg/cwebauthn"
	"nugg-auth/core/pkg/webauthn/webauthn"

	"strings"

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

func Base64ToDynamoPrimaryKey(s string) string {
	return base32.StdEncoding.EncodeToString([]byte(s))
}

func DynamoPrimaryKeyToBase64(s string) string {
	b, err := base32.StdEncoding.DecodeString(s)
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

func clientDataChallengeToDynamoPrimaryKey(clientData string, expectedChallengeType string) (string, error) {

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

	// this prob will not do anything, it is just to be safe because we are raw encoding below
	base64Encoded := strings.TrimRight(cd.Challenge, "=")

	decoded, err := base64.RawStdEncoding.DecodeString(base64Encoded)
	if err != nil {
		log.Printf("failed to decode base64, %v", err)
		return "", err
	}

	return string(decoded), nil
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
		Id:          Base64ToDynamoPrimaryKey(session.Challenge),
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
func (client *Client) LoadChallenge(ctx context.Context, clientData string, expectedChallengeType string) (*webauthn.SessionData, *User, *cwebauthn.User, error) {

	key, err := clientDataChallengeToDynamoPrimaryKey(clientData, expectedChallengeType)
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
