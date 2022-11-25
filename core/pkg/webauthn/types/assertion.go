package types

import (
	"encoding/json"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/webauthn/extensions"
)

type VerifyAssertionInputArgs struct {
	Input               AssertionInput
	StoredChallenge     hex.Hash
	AttestationType     CredentialAttestationType
	VerifyUser          bool
	CredentialPublicKey hex.Hash
	Extensions          extensions.ClientInputs
	LastSignCount       uint64
	RelyingPartyID      string
	RelyingPartyOrigin  string
	DataSignedByClient  hex.Hash
}

type AssertionInput struct {
	UserID               hex.Hash `json:"userID"`
	CredentialID         hex.Hash `json:"credentialID"`
	RawClientDataJSON    string   `json:"rawClientDataJSON"`
	RawAuthenticatorData hex.Hash `json:"rawAuthenticatorData"`
	Signature            hex.Hash `json:"signature"`
	Type                 string   `json:"credentialType"`
}

type DataAssertion struct {
	SessionID         hex.Hash `json:"session_id"`
	RawClientDataJSON string   `json:"client_data_json"`
	AssertionObject   hex.Hash `json:"assersion_object"`
	ChallengeID       hex.Hash `json:"challenge_id"`
	CredentialID      hex.Hash `json:"credential_id"`
	Body              hex.Hash `json:"-"`
}

// convert stirng into DataAssertion

func NewDataAssertion(header hex.Hash, body hex.Hash) (*DataAssertion, error) {
	var car DataAssertion
	err := json.Unmarshal(body.Bytes(), &car)
	if err != nil {
		return nil, err
	}
	car.Body = body
	return &car, nil
}
