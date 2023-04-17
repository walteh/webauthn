package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	"git.nugg.xyz/webauthn/pkg/hex"
)

type VerifyClientDataArgs struct {
	ClientData         CollectedClientData
	StoredChallenge    hex.Hash
	CeremonyType       CeremonyType
	RelyingPartyOrigin string
}

// CollectedClientData represents the contextual bindings of both the WebAuthn Relying Party
// and the client. It is a key-value mapping whose keys are strings. Values can be any type
// that has a valid encoding in JSON. Its structure is defined by the following Web IDL.
// https://www.w3.org/TR/webauthn/#sec-client-data
type CollectedClientData struct {
	// Type the string "webauthn.create" when creating new credentials,
	// and "webauthn.get" when getting an assertion from an existing credential. The
	// purpose of this member is to prevent certain types of signature confusion attacks
	//(where an attacker substitutes one legitimate signature for another).
	Type         CeremonyType  `json:"type"`
	Challenge    hex.Hash      `json:"challenge"`
	Origin       string        `json:"origin"`
	TokenBinding *TokenBinding `json:"tokenBinding,omitempty"`
	// Chromium (Chrome) returns a hint sometimes about how to handle clientDataJSON in a safe manner
	Hint string `json:"new_keys_may_be_added_here,omitempty"`
}

// custom unmarshal to handle converting challenge from base64 to hex
func (c *CollectedClientData) UnmarshalJSON(data []byte) error {
	type Alias CollectedClientData
	aux := &struct {
		Challenge string `json:"challenge"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	challenge, err := base64.RawURLEncoding.DecodeString(aux.Challenge)
	if err != nil {
		return err
	}
	c.Challenge = hex.Hash(challenge)
	return nil
}

type CeremonyType string

const (
	CreateCeremony CeremonyType = "webauthn.create"
	AssertCeremony CeremonyType = "webauthn.get"
)

type TokenBinding struct {
	Status TokenBindingStatus `json:"status"`
	ID     string             `json:"id,omitempty"`
}

type TokenBindingStatus string

const (
	// Indicates token binding was used when communicating with the
	// Relying Party. In this case, the id member MUST be present.
	Present TokenBindingStatus = "present"
	// Indicates token binding was used when communicating with the
	// negotiated when communicating with the Relying Party.
	Supported TokenBindingStatus = "supported"
	// Indicates token binding not supported
	// when communicating with the Relying Party.
	NotSupported TokenBindingStatus = "not-supported"
)

// Returns the origin per the HTML spec: (scheme)://(host)[:(port)]
func FullyQualifiedOrigin(u *url.URL) string {
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}
