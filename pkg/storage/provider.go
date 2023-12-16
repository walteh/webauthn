package storage

import (
	"context"

	"github.com/walteh/webauthn/pkg/webauthn/types"
)

type Provider interface {
	WriteNewCeremony(ctx context.Context, crm *types.Ceremony) error
	// GetExistingCeremony(ctx context.Context, challenge string) (*types.Ceremony, error)
	GetExisting(ctx context.Context, crm types.CeremonyID, credid types.CredentialID) (*types.Ceremony, *types.Credential, error)
	// REQUIREMENT: the ceremony must exist in the db
	// REQUIREMENT: the credential must not exist in the db
	// REQUIREMENT: the sessionid of the credential must match the sessionid of the ceremony in the db
	WriteNewCredential(ctx context.Context, crm types.CeremonyID, cred *types.Credential) error
	IncrementExistingCredential(ctx context.Context, crm types.CeremonyID, credid types.CredentialID) error
}
