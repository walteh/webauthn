package providers

import (
	"time"

	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/types"
)

type NoneAttestationProvider struct{}

func NewNoneAttestationProvider() *NoneAttestationProvider {
	return &NoneAttestationProvider{}
}

func (me *NoneAttestationProvider) Time() time.Time {
	return time.Unix(0, 0)
}

func (me *NoneAttestationProvider) ID() string {
	return "none"
}

func (me *NoneAttestationProvider) Attest(att types.AttestationObject, clientDataHash []byte) (hex.Hash, string, []interface{}, error) {
	return hex.Hash{}, "", []interface{}{}, nil
}
