package dynamo

import (
	"nugg-webauthn/core/pkg/webauthn/protocol"

	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Ceremony struct {
	Id        string `dynamodbav:"ceremony_id"`
	SessionId string `dynamodbav:"session_id"`
	Ttl       int64  `dynamodbav:"ttl"`
}

func newCeremony(challenge, sessionId string) *Ceremony {
	return &Ceremony{
		Id:        challenge,
		SessionId: sessionId,
		Ttl:       (time.Now().Unix()) + 300,
	}
}

func (client *Client) makeCeremonyPut(c interface{}) (*types.Put, error) {
	av, err := attributevalue.MarshalMap(c)
	if err != nil {
		return nil, err
	}

	return &types.Put{Item: av, TableName: client.MustCeremonyTableName()}, nil
}

func (client *Client) NewCeremonyPut(challenge, sessionId string) (*types.Put, error) {
	return client.makeCeremonyPut(newCeremony(challenge, sessionId))
}

func (client *Client) FindCeremonyInGetResult(result []*GetOutput) (cer *Ceremony, err error) {
	cred, err := FindInOnePerTableGetResult[Ceremony](result, client.MustCeremonyTableName())
	if err != nil {
		return nil, err
	}

	return cred, nil
}

func (client *Client) NewCeremonyGet(challenge string) *types.Get {
	return &types.Get{
		Key: map[string]types.AttributeValue{
			"ceremony_id": &types.AttributeValueMemberS{Value: protocol.ResolveToRawURLEncoding(challenge)},
		},
		TableName: client.MustCeremonyTableName(),
	}
}
