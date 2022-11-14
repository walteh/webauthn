package webauthn

import (
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

// Path: core/pkg/webauthn/config.go

// Config is the configuration for the webauthn package.
type WebAuthn struct {
	*webauthn.WebAuthn
}

func NewConfig() (*WebAuthn, error) {
	web, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "nugg.xyz",
		RPID:          "nugg.xyz",
		RPOrigin:      "https://auth.nugg.xyz",
		// AuthenticatorSelection: protocol.AuthenticatorSelection{
		// 	AuthenticatorAttachment: protocol.AuthenticatorAttachment("apple"),
		// 	UserVerification:        protocol.VerificationRequired,
		// 	ResidentKey:             protocol.ResidentKeyRequirementRequired,
		// 	RequireResidentKey:      protocol.ResidentKeyRequired(),
		// },
		AttestationPreference: protocol.PreferDirectAttestation,
	})
	return &WebAuthn{web}, err
}
