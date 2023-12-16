package storage

import (
	"context"

	"github.com/walteh/webauthn/pkg/webauthn/types"
)

type Provider interface {
	WriteNewCeremony(ctx context.Context, crm *types.Ceremony) error
	// GetExistingCeremony(ctx context.Context, challenge string) (*types.Ceremony, error)
	GetExisting(ctx context.Context, crm types.CeremonyID, credid types.CredentialID) (*types.Ceremony, *types.Credential, error)
	// GetExistingCredential(ctx context.Context, credid types.CredentialID) (*types.Credential, error)
	WriteNewCredential(ctx context.Context, crm types.CeremonyID, cred *types.Credential) error
	IncrementExistingCredential(ctx context.Context, crm types.CeremonyID, credid types.CredentialID) error
}
