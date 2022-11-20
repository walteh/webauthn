package protocol

import (
	"encoding/base64"
	"fmt"
	"nugg-auth/core/pkg/hex"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SavedCredential struct {
	// A probabilistically-unique byte sequence identifying a public key credential source and its authentication assertions.
	RawID hex.Hash `dynamodbav:"-" json:"raw_credential_id"`

	Type CredentialType `dynamodbav:"credential_type" json:"credential_type"`
	// The public key portion of a Relying Party-specific credential key pair, generated by an authenticator and returned to
	// a Relying Party at registration time (see also public key credential). The private key portion of the credential key
	// pair is known as the credential private key. Note that in the case of self attestation, the credential key pair is also
	// used as the attestation key pair, see self attestation for details.
	PublicKey hex.Hash `dynamodbav:"public_key" json:"public_key"`
	// The attestation format used (if any) by the authenticator when creating the credential.
	AttestationType string `dynamodbav:"attestation_type" json:"attestation_type"`
	// The Authenticator information for a given certificate

	Receipt hex.Hash `dynamodbav:"receipt" json:"receipt"`

	// The AAGUID of the authenticator. An AAGUID is defined as an array containing the globally unique
	// identifier of the authenticator model being sought.
	AAGUID hex.Hash `dynamodbav:"aaguid" json:"aaguid"`
	// SignCount -Upon a new login operation, the Relying Party compares the stored signature counter value
	// with the new signCount value returned in the assertion’s authenticator data. If this new
	// signCount value is less than or equal to the stored value, a cloned authenticator may
	// exist, or the authenticator may be malfunctioning.
	SignCount uint64 `dynamodbav:"sign_count" json:"sign_count"`
	// CloneWarning - This is a signal that the authenticator may be cloned, i.e. at least two copies of the
	// credential private key may exist and are being used in parallel. Relying Parties should incorporate
	// this information into their risk scoring. Whether the Relying Party updates the stored signature
	// counter value in this case, or not, or fails the authentication ceremony or not, is Relying Party-specific.
	CloneWarning bool `dynamodbav:"clone_warning" json:"clone_warning"`

	CreatedAt uint64 `dynamodbav:"created_at"          json:"created_at"`
	UpdatedAt uint64 `dynamodbav:"updated_at"          json:"updated_at"`

	SessionId hex.Hash `dynamodbav:"session_id" json:"session_id"`
}

// type Marshaler interface {
// 	MarshalDynamoDBAttributeValue() (types.AttributeValue, error)
// }

func (s SavedCredential) ID() string {
	return (s.RawID.Hex())
}

func (s SavedCredential) MarshalDynamoDBAttributeValue() (*types.AttributeValueMemberM, error) {
	av := types.AttributeValueMemberM{}
	av.Value = make(map[string]types.AttributeValue)
	av.Value["credential_id"] = &types.AttributeValueMemberS{Value: s.ID()}
	av.Value["credential_type"] = &types.AttributeValueMemberS{Value: string(s.Type)}
	av.Value["public_key"] = &types.AttributeValueMemberS{Value: s.PublicKey.Hex()}
	av.Value["attestation_type"] = &types.AttributeValueMemberS{Value: s.AttestationType}
	av.Value["receipt"] = &types.AttributeValueMemberS{Value: s.Receipt.Hex()}
	av.Value["aaguid"] = &types.AttributeValueMemberS{Value: s.AAGUID.Hex()}
	av.Value["sign_count"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", s.SignCount)}
	av.Value["clone_warning"] = &types.AttributeValueMemberBOOL{Value: s.CloneWarning}
	av.Value["created_at"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", s.CreatedAt)}
	av.Value["updated_at"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", s.UpdatedAt)}
	av.Value["session_id"] = &types.AttributeValueMemberS{Value: s.SessionId.Hex()}
	return &av, nil
}

//type Unmarshaler interface {
// 	UnmarshalDynamoDBAttributeValue(types.AttributeValue) error
// }

// func Parse[AV types.AttributeValue, T interface{}](m types.AttributeValueMemberM, field string) (T, bool) {
// 	var v types.AttributeValue
// 	if v, ok := m.Value[field].(AV); ok {
// 		s.ID = m.Value[field].(AV).isAttributeValue()
// 	}
// }

func (s SavedCredential) UnmarshalDynamoDBAttributeValue(m *types.AttributeValueMemberM) (err error) {

	if s.RawID, err = base64.RawURLEncoding.DecodeString(m.Value["credential_id"].(*types.AttributeValueMemberS).Value); err != nil {
		return err
	}

	s.RawID = hex.HexToHash(m.Value["credential_id"].(*types.AttributeValueMemberS).Value)

	s.Type = CredentialType(m.Value["credential_type"].(*types.AttributeValueMemberS).Value)
	s.PublicKey = hex.HexToHash(m.Value["public_key"].(*types.AttributeValueMemberS).Value)
	s.AttestationType = m.Value["attestation_type"].(*types.AttributeValueMemberS).Value
	s.Receipt = hex.HexToHash(m.Value["receipt"].(*types.AttributeValueMemberS).Value)
	s.AAGUID = hex.HexToHash(m.Value["aaguid"].(*types.AttributeValueMemberS).Value)
	s.SignCount, err = strconv.ParseUint(m.Value["sign_count"].(*types.AttributeValueMemberN).Value, 10, 64)
	s.CloneWarning = m.Value["clone_warning"].(*types.AttributeValueMemberBOOL).Value
	s.CreatedAt, err = strconv.ParseUint(m.Value["created_at"].(*types.AttributeValueMemberN).Value, 10, 64)
	s.UpdatedAt, err = strconv.ParseUint(m.Value["updated_at"].(*types.AttributeValueMemberN).Value, 10, 64)
	s.SessionId = hex.HexToHash(m.Value["session_id"].(*types.AttributeValueMemberS).Value)
	return
}

func (s SavedCredential) Update(table *string, counter uint64) (*types.TransactWriteItem, error) {
	s.updateCounter(counter)

	s.UpdatedAt = uint64(time.Now().Unix())
	av, err := s.MarshalDynamoDBAttributeValue()
	if err != nil {
		return nil, err
	}

	return &types.TransactWriteItem{
		Update: &types.Update{
			TableName: table,
			Key: map[string]types.AttributeValue{
				"credential_id": &types.AttributeValueMemberS{Value: s.ID()},
			},
			ExpressionAttributeNames: map[string]string{
				"#e": "sign_count",
				"#f": "clone_warning",
				"#g": "updated_at",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":e": av.Value["sign_count"],
				":f": av.Value["clone_warning"],
				":g": av.Value["updated_at"],
			},
			UpdateExpression: aws.String("SET #e = :e, #f = :f, #g = :g"),
		}}, nil
}

func SavedCredentialGet(table *string, id string) *types.Get {
	return &types.Get{
		TableName: table,
		Key: map[string]types.AttributeValue{
			"credential_id": &types.AttributeValueMemberS{Value: id},
		},
	}
}

func (a *SavedCredential) updateCounter(authDataCount uint64) error {
	if authDataCount <= a.SignCount && (authDataCount != 0 || a.SignCount != 0) {
		a.CloneWarning = true
		return fmt.Errorf("authDataCount %d <= a.SignCount %d", authDataCount, a.SignCount)
	}
	a.SignCount = authDataCount
	return nil
}

func (s SavedCredential) Get() *types.Get {
	return &types.Get{
		Key: map[string]types.AttributeValue{
			"credential_id": &types.AttributeValueMemberS{Value: s.ID()},
		},
	}
}

func (s SavedCredential) Put() (*types.Put, error) {

	av, err := s.MarshalDynamoDBAttributeValue()
	if err != nil {
		return nil, err
	}

	return &types.Put{
		Item: av.Value,
	}, nil
}

func TableType() string {
	return "credential"
}

func NewUnsafeGettableCredential(id hex.Hash) *SavedCredential {
	return &SavedCredential{
		RawID: id,
	}
}
