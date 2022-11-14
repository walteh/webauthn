package dynamo

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/duo-labs/webauthn/webauthn"
)

var ErrCeremonyExpired = errors.New("ceremony_expired")
var ErrDynamoConnectionError = errors.New("ceremony_expired")

type Ceremony struct {
	UserId string `dynamodbav:"ceremony_user_id"`
	Data   []byte `dynamodbav:"data"`
	Ttl    int64  `dynamodbav:"ttl"`
}

func (client *Client) StartWebAuthnCeremony(ctx context.Context, user_id string, session *webauthn.SessionData) error {

	j, err := json.Marshal(session)
	if err != nil {
		log.Printf("failed to marshal session, %v", err)
		return err
	}

	challenge := Ceremony{
		UserId: user_id,
		Data:   j,
		Ttl:    time.Now().Add(5 * time.Minute).Unix(),
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

func (client *Client) FindWebAuthnCeremony(ctx context.Context, user_id string) (*webauthn.SessionData, error) {

	res, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(client.TableName),
		Key: map[string]types.AttributeValue{
			"ceremony_user_id": &types.AttributeValueMemberS{Value: user_id},
		},
	})
	if err != nil {
		log.Printf("Couldn't get item to table. Here's why: %v\n", err)
		return nil, ErrDynamoConnectionError
	}

	if len(res.Item) == 0 {
		return nil, ErrCeremonyExpired
	}

	var ceremony Ceremony
	err = attributevalue.UnmarshalMap(res.Item, &ceremony)
	if err != nil {
		log.Printf("failed to unmarshal item, %v", err)
		return nil, ErrCeremonyExpired
	}

	var session webauthn.SessionData
	err = json.Unmarshal(ceremony.Data, &session)
	if err != nil {
		log.Printf("failed to unmarshal session, %v", err)
		return nil, ErrCeremonyExpired
	}

	return &session, nil

}
