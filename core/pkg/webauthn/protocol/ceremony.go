package protocol

import (
	"fmt"
	"strconv"
	"time"

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

func (s SavedCeremony) MarshalDynamoDBAttributeValue() (*types.AttributeValueMemberM, error) {
	av := types.AttributeValueMemberM{}
	av.Value = make(map[string]types.AttributeValue)
	av.Value["challenge_id"] = &types.AttributeValueMemberS{Value: s.ChallengeID.Hex()}
	av.Value["session_id"] = &types.AttributeValueMemberS{Value: s.ChallengeID.Hex()}
	av.Value["credential_id"] = &types.AttributeValueMemberS{Value: s.CredentialID.Hex()}
	av.Value["ceremony_type"] = &types.AttributeValueMemberS{Value: string(s.CeremonyType)}
	av.Value["created_at"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", s.CreatedAt)}
	av.Value["ttl"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", s.Ttl)}
	return &av, nil
}

//type Unmarshaler interface {
// 	UnmarshalDynamoDBAttributeValue(types.AttributeValue) error
// }

func (s *SavedCeremony) UnmarshalDynamoDBAttributeValue(av *types.AttributeValueMemberM) (err error) {
	s.ChallengeID = hex.HexToHash(av.Value["challenge_id"].(*types.AttributeValueMemberS).Value)
	s.SessionID = hex.HexToHash(av.Value["session_id"].(*types.AttributeValueMemberS).Value)
	s.CredentialID = hex.HexToHash(av.Value["credential_id"].(*types.AttributeValueMemberS).Value)
	s.CeremonyType = CeremonyType(av.Value["ceremony_type"].(*types.AttributeValueMemberS).Value)

	if s.CreatedAt, err = strconv.ParseUint(av.Value["created_at"].(*types.AttributeValueMemberN).Value, 10, 64); err != nil {
		return err
	}
	if s.Ttl, err = strconv.ParseUint(av.Value["ttl"].(*types.AttributeValueMemberN).Value, 10, 64); err != nil {
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
		CredentialID: credentialID,
		SessionID:    sessionId,
		ChallengeID:  chal,
		CeremonyType: ceremonyType,
		CreatedAt:    Now(),
		Ttl:          Now() + 300,
	}

	return cer
}

func Now() uint64 {
	return uint64(time.Now().Unix())
}

func (s SavedCeremony) Get() *types.Get {
	return &types.Get{
		Key: map[string]types.AttributeValue{
			"challenge_id": &types.AttributeValueMemberS{Value: s.ChallengeID.Hex()},
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

	return &types.Put{
		Item: av.Value,
	}, nil
}
