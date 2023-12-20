package types

import (
	"encoding/json"

	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/extensions"
)

type Asserter interface {
	VerifyAssertion(input VerifyAssertionInputArgs) (bool, error)
}

type VerifyAssertionInputArgs struct {
	Input                          AssertionInput
	StoredChallenge                CeremonyID
	CredentialAttestationType      CredentialAttestationType
	AttestationProvider            AttestationProvider
	VerifyUser                     bool
	AAGUID                         hex.Hash
	CredentialPublicKey            hex.Hash
	Extensions                     extensions.ClientInputs
	LastSignCount                  uint64
	RelyingPartyID                 string
	RelyingPartyOrigin             string
	DataSignedByClient             hex.Hash
	UseSavedAttestedCredentialData bool
}

type AssertionObject struct {
	// AuthenticatorData    AuthenticatorData
	RawAuthenticatorData hex.Hash
	Signature            hex.Hash
}

type AssertionInput struct {
	UserID            hex.Hash     `json:"userID"`
	CredentialID      CredentialID `json:"credentialID"`
	RawClientDataJSON string       `json:"rawClientDataJSON"`
	// RawAuthenticatorData hex.Hash `json:"rawAuthenticatorData"`
	RawAssertionObject hex.Hash `json:"signature"`
	// Type            string   `json:"credentialType"`
	AssertionObject *AssertionObject
}

func (a *AssertionInput) Validate() error {
	if len(a.UserID) == 0 || len(a.CredentialID) == 0 || len(a.RawClientDataJSON) == 0 || (len(a.RawAssertionObject) == 0 && a.AssertionObject == nil) {
		return ErrInvalidAssertionInput
	}
	return nil
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
