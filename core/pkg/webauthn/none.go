package protocol

import (
	"nugg-webauthn/core/pkg/hex"
)

type NoneAttestationProvider struct{}

func NewNoneAttestationProvider() *NoneAttestationProvider {
	return &NoneAttestationProvider{}
}

func (me *NoneAttestationProvider) ID() string {
	return "none"
}

func (me *NoneAttestationProvider) Attest(att AttestationObject, clientDataHash []byte) (hex.Hash, string, []interface{}, error) {
	return hex.Hash{}, "", []interface{}{}, nil
}
