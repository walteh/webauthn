package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/challenge"
)

type CeremonyID hex.Hash

func (s CeremonyID) Ref() hex.Hash {
	return hex.Hash(s)
}

func (s CeremonyID) MarshalText() ([]byte, error) {
	return s.Ref().MarshalText()
}

type Ceremony struct {
	ChallengeID  CeremonyID   `dynamodbav:"challenge_id" json:"challenge_id"`
	SessionID    hex.Hash     `dynamodbav:"session_id" json:"session_id"`
	CredentialID CredentialID `dynamodbav:"credential_id,omitempty" json:"credential_id,omitempty"`
	CeremonyType CeremonyType `dynamodbav:"ceremony_type" json:"ceremony_type"`
	CreatedAt    uint64       `dynamodbav:"created_at" json:"created_at"`
	Ttl          uint64       `dynamodbav:"ttl" json:"ttl"`
}

// type Marshaler interface {
// 	MarshalDynamoDBAttributeValue() (types.AttributeValue, error)
// }

func (s Ceremony) MarshalDynamoDBAttributeValue() (*types.AttributeValueMemberM, error) {
	av := types.AttributeValueMemberM{}
	av.Value = make(map[string]types.AttributeValue)
	av.Value["challenge_id"] = &types.AttributeValueMemberS{Value: s.ChallengeID.Ref().Hex()}
	av.Value["session_id"] = &types.AttributeValueMemberS{Value: s.SessionID.Hex()}
	av.Value["credential_id"] = &types.AttributeValueMemberS{Value: s.CredentialID.Ref().Hex()}
	av.Value["ceremony_type"] = &types.AttributeValueMemberS{Value: string(s.CeremonyType)}
	av.Value["created_at"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", s.CreatedAt)}
	av.Value["ttl"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", s.Ttl)}
	return &av, nil
}

//type Unmarshaler interface {
// 	UnmarshalDynamoDBAttributeValue(types.AttributeValue) error
// }

func (s *Ceremony) UnmarshalDynamoDBAttributeValue(av *types.AttributeValueMemberM) (err error) {

	if av.Value == nil {
		return errors.New("empty attribute value")
	}

	if chal, err := GetSHashNotZero(av, "challenge_id"); err != nil {
		return err
	} else {
		s.ChallengeID = CeremonyID(chal)
	}

	if s.SessionID, err = GetSHashNotZero(av, "session_id"); err != nil {
		return err
	}

	if cred, err := GetSHash(av, "credential_id"); err != nil {
		return err
	} else {
		s.CredentialID = CredentialID(cred)
	}

	if r, err := GetS(av, "ceremony_type"); err != nil {
		return err
	} else {
		s.CeremonyType = CeremonyType(r)
	}

	if s.CreatedAt, err = GetNUint64(av, "created_at"); err != nil {
		return err
	}

	if s.Ttl, err = GetNUint64(av, "ttl"); err != nil {
		return err
	}

	return nil
}

func NewCeremony(credentialID CredentialID, sessionId hex.Hash, ceremonyType CeremonyType) *Ceremony {

	chal, err := challenge.CreateChallenge()
	if err != nil {
		panic(err)
	}

	cer := &Ceremony{
		CredentialID: credentialID,
		SessionID:    sessionId,
		ChallengeID:  CeremonyID(chal),
		CeremonyType: ceremonyType,
		CreatedAt:    Now(),
		Ttl:          Now() + 300,
	}

	return cer
}

func Now() uint64 {
	return uint64(time.Now().Unix())
}

// func (s Ceremony) Get() *types.Get {
// 	return &types.Get{
// 		Key: map[string]types.AttributeValue{
// 			"challenge_id": &types.AttributeValueMemberS{Value: s.ChallengeID.Hex()},
// 		},
// 	}
// }

func (s Ceremony) WasGot() bool {
	return s.CreatedAt != 0
}

// func NewUnsafeGettableCeremony(id hex.Hash) *Ceremony {
// 	return &Ceremony{
// 		ChallengeID: id,
// 	}
// }

func (s Ceremony) Put() (*types.Put, error) {

	av, err := s.MarshalDynamoDBAttributeValue()
	if err != nil {
		return nil, err
	}

	return &types.Put{
		Item: av.Value,
	}, nil
}
