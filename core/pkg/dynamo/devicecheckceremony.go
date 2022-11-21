package dynamo

import (
	"nugg-webauthn/core/pkg/webauthn/protocol"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DeviceCheckCeremony struct {
	Id        string `dynamodbav:"ceremony_id"`
	SessionID string `dynamodbav:"session_data"`
	Ttl       int64  `dynamodbav:"ttl"`
}

func newDeviceCheckCeremony(challenge, sessionId string) *DeviceCheckCeremony {
	return &DeviceCheckCeremony{
		Id:        challenge,
		SessionID: sessionId,
		Ttl:       (time.Now().Unix()) + 300,
	}
}

func (client *Client) makeDeviceCheckCeremonyPut(c interface{}) (*types.Put, error) {
	av, err := attributevalue.MarshalMap(c)
	if err != nil {
		return nil, err
	}

	return &types.Put{Item: av, TableName: client.MustCeremonyTableName()}, nil
}

func (client *Client) NewDeviceCheckCeremonyPut(challenge, sessionId string) (*types.Put, error) {
	return client.makeDeviceCheckCeremonyPut(newDeviceCheckCeremony(challenge, sessionId))
}

func (client *Client) FindDeviceCheckCeremonyInGetResult(result []*GetOutput) (cer *DeviceCheckCeremony, err error) {

	cred, err := FindInOnePerTableGetResult[DeviceCheckCeremony](result, client.MustCeremonyTableName())
	if err != nil {
		return nil, err
	}

	return cred, nil
}

func (client *Client) NewDeviceCheckCeremonyGet(challenge string) *types.Get {
	return &types.Get{
		Key: map[string]types.AttributeValue{
			"ceremony_id": &types.AttributeValueMemberS{Value: protocol.ResolveToRawURLEncoding(challenge)},
		},
		TableName: client.MustCeremonyTableName(),
	}
}
