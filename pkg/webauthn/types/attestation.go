package types

import (
	"github.com/nuggxyz/webauthn/pkg/errors"
	"github.com/nuggxyz/webauthn/pkg/hex"
	"github.com/nuggxyz/webauthn/pkg/webauthn/extensions"
)

func NewDefaultCredentialIdentifier(CredentialID hex.Hash) *CredentialIdentifier {
	return &CredentialIdentifier{
		ID:               CredentialID,
		Type:             "public-key",
		ExtensionResults: nil,
	}
}

type CredentialIdentifier struct {
	// ID of the credential
	ID   hex.Hash
	Type CredentialType

	ExtensionResults extensions.ClientOutputs
}

type VerifyAttestationInputArgs struct {
	Provider           AttestationProvider
	Input              AttestationInput
	StoredChallenge    hex.Hash
	SessionId          hex.Hash
	VerifyUser         bool
	RelyingPartyID     string
	RelyingPartyOrigin string
}

// From §5.2.1 (https://www.w3.org/TR/webauthn/#authenticatorattestationresponse)
// "The authenticator's response to a client’s request for the creation
// of a new public key credential. It contains information about the new credential
// that can be used to identify it for later use, and metadata that can be used by
// the WebAuthn Relying Party to assess the characteristics of the credential
// during registration."

// The initial unpacked 'response' object received by the relying party. This
// contains the clientDataJSON object, which will be marshalled into
// CollectedClientData, and the 'attestationObject', which contains
// information about the authenticator, and the newly minted
// public key credential. The information in both objects are used
// to verify the authenticity of the ceremony and new credential
type AttestationInput struct {
	// The byte slice of clientDataJSON, which becomes CollectedClientData
	// AuthenticatorResponse
	UTF8ClientDataJSON string `json:"client_data_json"`
	// The byte slice version of AttestationObject
	// This attribute contains an attestation object, which is opaque to, and
	// cryptographically protected against tampering by, the client. The
	// attestation object contains both authenticator data and an attestation
	// statement. The former contains the AAGUID, a unique credential ID, and
	// the credential public key. The contents of the attestation statement are
	// determined by the attestation statement format used by the authenticator.
	// It also contains any additional information that the Relying Party's server
	// requires to validate the attestation statement, as well as to decode and
	// validate the authenticator data along with the JSON-serialized client data.
	AttestationObject hex.Hash `json:"attestation_object"`

	CredentialID     hex.Hash                 `json:"credential_id"`
	CredentialType   CredentialType           `json:"credential_type"`
	ClientExtensions extensions.ClientOutputs `json:"client_extensions"`
}

// From §6.4. Authenticators MUST also provide some form of attestation. The basic requirement is that the
// authenticator can produce, for each credential public key, an attestation statement verifiable by the
// WebAuthn Relying Party. Typically, this attestation statement contains a signature by an attestation
// private key over the attested credential public key and a challenge, as well as a certificate or similar
// data providing provenance information for the attestation public key, enabling the Relying Party to make
// a trust decision. However, if an attestation key pair is not available, then the authenticator MUST
// perform self attestation of the credential public key with the corresponding credential private key.
// All this information is returned by authenticators any time a new public key credential is generated, in
// the overall form of an attestation object. (https://www.w3.org/TR/webauthn/#attestation-object)
type AttestationObject struct {
	// The authenticator data, including the newly created public key. See AuthenticatorData for more info
	AuthData AuthenticatorData
	// The byteform version of the authenticator data, used in part for signature validation
	RawAuthData hex.Hash `json:"authData" cbor:"authData"`
	// The format of the Attestation data.
	Format string `json:"fmt" cbor:"fmt"`
	// The attestation statement data sent back if attestation is requested.
	AttStatement map[string]interface{} `json:"attStmt,omitempty"`

	ClientData CollectedClientData `json:"-"`

	Extensions extensions.ClientOutputs `json:"-"`
}

type AttestationProvider interface {
	Attest(AttestationObject, []byte) (hex.Hash, string, []interface{}, error)
	ID() string
}

func (me CredentialIdentifier) Verify() error {
	if me.ID.IsZero() {
		return errors.ErrBadRequest.WithMessage("Parse error for Registration").WithInfo("Missing ID")
	}

	// testB64, err := base64.RawURLEncoding.DecodeString(ccr.ID)
	// if err != nil || !(len(testB64) > 0) {
	// 	return errors.ErrBadRequest.WithDetails("Parse error for Registration").WithInfo("ID not base64.RawURLEncoded").WithRoot(err)
	// }

	if me.Type == "" {
		return errors.ErrBadRequest.WithMessage("Parse error for Registration").WithInfo("Missing type")
	}

	if me.Type != "public-key" {
		return errors.ErrBadRequest.WithMessage("Parse error for Registration").WithInfo("Type not public-key")
	}

	return nil
}

type CredentialAttestationType string

const (
	NotFidoAttestationType CredentialAttestationType = "none"
	FidoAttestationType    CredentialAttestationType = "fido-u2f"
)

// GetAppID takes a AuthenticationExtensions object or nil. It then performs the following checks in order:
//
// 1. Check that the Session Data's AuthenticationExtensions has been provided and return a blank appid if it hasn't been.
// 2. Check that the AuthenticationExtensionsClientOutputs contains the extensions output and return a blank appid if it doesn't.
// 3. Check that the Credential AttestationType is `fido-u2f` and return a blank appid if it isn't.
// 4. Check that the AuthenticationExtensionsClientOutputs contains the appid key and return a blank appid if it doesn't.
// 5. Check that the AuthenticationExtensionsClientOutputs appid is a bool and return an error if it isn't.
// 6. Check that the appid output is true and return a blank appid if it isn't.
// 7. Check that the Session Data has an appid extension defined and return an error if it doesn't.
// 8. Check that the appid extension in Session Data is a string and return an error if it isn't.
// 9. Return the appid extension value from the Session Data.
func (ppkc CredentialIdentifier) GetAppID(authExt extensions.ClientInputs, credentialAttestationType CredentialAttestationType) (appID string, err error) {
	var (
		value, clientValue interface{}
		enableAppID, ok    bool
	)

	if ppkc.ExtensionResults == nil {
		return "", nil
	}

	// If the credential does not have the correct attestation type it is assumed to NOT be a fido-u2f credential.
	// https://w3c.github.io/webauthn/#sctn-fido-u2f-attestation
	if credentialAttestationType != FidoAttestationType {
		return "", nil
	}

	if clientValue, ok = ppkc.ExtensionResults["appid"]; !ok {
		return "", nil
	}

	if enableAppID, ok = clientValue.(bool); !ok {
		return "", errors.ErrBadRequest.WithMessage("Client Output appid did not have the expected type")
	}

	if !enableAppID {
		return "", nil
	}

	if value, ok = authExt["appid"]; !ok {
		return "", errors.ErrBadRequest.WithMessage("Session Data does not have an appid but Client Output indicates it should be set")
	}

	if appID, ok = value.(string); !ok {
		return "", errors.ErrBadRequest.WithMessage("Session Data appid did not have the expected type")
	}

	return appID, nil
}
