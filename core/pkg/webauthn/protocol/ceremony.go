package protocol

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"nugg-auth/core/pkg/hex"
)

type SavedCeremony struct {
	ChallengeID  hex.Hash     `dynamodbav:"challenge_id" json:"challenge_id"`
	SessionID    hex.Hash     `dynamodbav:"session_id" json:"session_id"`
	CredentialID hex.Hash     `dynamodbav:"credential_id,omitempty" json:"credential_id,omitempty"`
	CeremonyType CeremonyType `dynamodbav:"ceremony_type" json:"ceremony_type"`
	CreatedAt    uint64       `dynamodbav:"created_at" json:"created_at"`
	Ttl          uint64       `dynamodbav:"ttl" json:"ttl"`
}

// type Marshaler interface {
// 	MarshalDynamoDBAttributeValue() (types.AttributeValue, error)
// }

func (s SavedCeremony) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	av := types.AttributeValueMemberM{}
	av.Value = make(map[string]types.AttributeValue)
	av.Value["challenge_id"] = &types.AttributeValueMemberS{Value: s.ChallengeID.Hex()}
	av.Value["session_id"] = &types.AttributeValueMemberS{Value: s.ChallengeID.Hex()}
	av.Value["credential_id"] = &types.AttributeValueMemberS{Value: s.CredentialID.Hex()}
	av.Value["ceremony_type"] = &types.AttributeValueMemberS{Value: string(s.CeremonyType)}
	av.Value["created_at"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", s.CreatedAt)}
	av.Value["ttl"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", s.Ttl)}
	return &av, nil
}

//type Unmarshaler interface {
// 	UnmarshalDynamoDBAttributeValue(types.AttributeValue) error
// }

func (s SavedCeremony) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) (err error) {

	// plain unmarshal
	err = attributevalue.Unmarshal(av, &s)
	if err != nil {
		return err
	}

	return nil

}

func NewCeremony(credentialID hex.Hash, sessionId hex.Hash, ceremonyType CeremonyType) *SavedCeremony {

	chal, err := CreateChallenge()
	if err != nil {
		panic(err)
	}

	cer := &SavedCeremony{
		SessionID:    sessionId,
		ChallengeID:  chal,
		CeremonyType: ceremonyType,
		CreatedAt:    Now(),
		Ttl:          Now() + 300,
	}

	if credentialID.IsZero() {
		cer.CredentialID = chal
	}

	return cer
}

func Now() uint64 {
	return uint64(time.Now().Unix())
}

func (s SavedCeremony) Get() *types.Get {
	return &types.Get{
		Key: map[string]types.AttributeValue{
			"credential_id": &types.AttributeValueMemberS{Value: s.ChallengeID.Hex()},
		},
	}
}

func NewUnsafeGettableCeremony(id hex.Hash) *SavedCeremony {
	return &SavedCeremony{
		ChallengeID: id,
	}
}

func (s SavedCeremony) Put() (*types.Put, error) {

	av, err := s.MarshalDynamoDBAttributeValue()
	if err != nil {
		return nil, err
	}

	r, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return nil, fmt.Errorf("expected *types.AttributeValueMemberM, got %T", av)
	}

	return &types.Put{
		Item: r.Value,
	}, nil
}
