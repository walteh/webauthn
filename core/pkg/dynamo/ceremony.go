package dynamo

import (
	"context"
	"log"
	"nugg-auth/core/pkg/webauthn/webauthn"

	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Ceremony struct {
	Id          string                `dynamodbav:"id"`
	UserId      string                `dynamodbav:"userId"`
	SessionData *webauthn.SessionData `dynamodbav:"session_data"`
	Ttl         int64                 `dynamodbav:"ttl"`
}

func NewCeremony(userId string, session *webauthn.SessionData) *Ceremony {
	return &Ceremony{
		Id:          session.Challenge,
		UserId:      userId,
		SessionData: session,
		Ttl:         (time.Now().Unix()) + 300000,
	}
}

func (client *Client) SaveCeremony(ctx context.Context, cer *Ceremony) error {

	item, err := attributevalue.MarshalMap(cer)
	if err != nil {
		log.Printf("failed to marshal challenge, %v", err)
		return err
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: client.MustCeremonyTableName(), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
		return err
	}
	return nil
}

// / load challenge
func (client *Client) LoadCeremony(ctx context.Context, challenge string) (*Ceremony, error) {

	input := &dynamodb.GetItemInput{
		TableName: client.MustCeremonyTableName(),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: challenge},
		},
	}

	result, err := client.GetItem(ctx, input)
	if err != nil {
		log.Printf("Got error calling GetItem: %s	", err)
		return nil, err
	}

	if result.Item == nil {
		return nil, ErrNotFound
	}

	var cer Ceremony

	err = attributevalue.UnmarshalMap(result.Item, &cer)
	if err != nil {
		log.Printf("Got error unmarshalling: %s", err)
		return nil, err
	}

	return &cer, nil
}
