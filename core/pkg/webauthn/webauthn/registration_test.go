package webauthn

import (
	"testing"

	"nugg-auth/core/pkg/webauthn/protocol"
)

func TestRegistration_FinishRegistrationFailure(t *testing.T) {

	session := SessionData{
		UserID: []byte("ABC"),
	}

	webauthn := &WebAuthn{}
	credential, err := webauthn.FinishRegistration("123", session, nil)
	if err == nil {
		t.Errorf("FinishRegistration() error = nil, want %v", protocol.ErrBadRequest.Type)
	}
	if credential != nil {
		t.Errorf("FinishRegistration() credential = %v, want nil", credential)
	}
}