package providers

import (
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/webauthn/types"
)

type NoneAttestationProvider struct{}

func NewNoneAttestationProvider() *NoneAttestationProvider {
	return &NoneAttestationProvider{}
}

func (me *NoneAttestationProvider) ID() string {
	return "none"
}

func (me *NoneAttestationProvider) Attest(att types.AttestationObject, clientDataHash []byte) (hex.Hash, string, []interface{}, error) {
	return hex.Hash{}, "", []interface{}{}, nil
}