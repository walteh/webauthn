package storage

import (
	"context"

	"github.com/walteh/webauthn/pkg/webauthn/types"
)

type Provider interface {
	WriteNewCeremony(ctx context.Context, crm *types.Ceremony) error
	GetExistingCeremony(ctx context.Context, challenge string) (*types.Ceremony, error)
	WriteNewCredential(ctx context.Context, crm *types.Ceremony, cred *types.Credential) error
	IncrementExistingCredential(ctx context.Context, crm *types.Ceremony, cred *types.Credential) error
}
