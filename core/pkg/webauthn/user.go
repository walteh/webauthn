package webauthn

import "github.com/duo-labs/webauthn/webauthn"

// WebAuthn is the main struct for the WebAuthn library. It contains the
// configuration and state for the library.
type User struct {
	// Config contains the configuration for the library.
	id       []byte
	username string
	creds    []webauthn.Credential
}

func NewUser(id []byte, username string) *User {
	return &User{
		id:       id,
		username: username,
		creds:    []webauthn.Credential{},
	}
}

func NewCredentials(pub string) *webauthn.Credential {
	return &webauthn.Credential{
		PublicKey: []byte{},

		// AttestationType: protocol.AttestationTypeNone,
		// AttestationTrustPath: []byte{},
		// UserHandle:          []byte{},
	}
}

func (user *User) WebAuthnID() []byte {
	return user.id
}

func (user *User) WebAuthnName() string {
	return user.username
}

func (user *User) WebAuthnDisplayName() string {
	return user.username
}

func (user *User) WebAuthnIcon() string {
	return "https://pics.com/avatar.png"
}

func (user *User) WebAuthnCredentials() []webauthn.Credential {
	return user.creds
	// return []webauthn.Credential{
	// 	{
	// 		ID:        []byte("credentialID"),
	// 		PublicKey: []byte("credentialPublicKey"),
	// 		Authenticator: webauthn.Authenticator{
	// 			AAGUID:       []byte("authenticatorAAGUID"),
	// 			SignCount:    0,
	// 			CloneWarning: true,
	// 		},
	// 	},
	// }
}
