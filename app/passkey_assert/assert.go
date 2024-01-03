package passkey_assert

import (
	"context"

	"github.com/walteh/terrors"
	"github.com/walteh/webauthn/pkg/accesstoken"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/relyingparty"
	"github.com/walteh/webauthn/pkg/storage"
	"github.com/walteh/webauthn/pkg/webauthn/assertion"
	"github.com/walteh/webauthn/pkg/webauthn/clientdata"
	"github.com/walteh/webauthn/pkg/webauthn/extensions"
	"github.com/walteh/webauthn/pkg/webauthn/providers"
	"github.com/walteh/webauthn/pkg/webauthn/types"
	"golang.org/x/sync/errgroup"
)

type PasskeyAssertionInput struct {
	SessionID            hex.Hash `json:"userID"`
	CredentialID         hex.Hash `json:"credentialID"`
	UTF8ClientDataJSON   string   `json:"rawClientDataJSON"`
	RawAuthenticatorData hex.Hash `json:"rawAuthenticatorData"`
	RawSignature         hex.Hash `json:"signature"`
	// PublicKey            hex.Hash `json:"publicKey"`
	// AAGUID               hex.Hash `json:"aaguid"`
}

type PasskeyAssertionOutput struct {
	AccessToken string
}

func Assert(ctx context.Context, store storage.Provider, rp relyingparty.Provider, tknp accesstoken.Provider, assert *PasskeyAssertionInput) (*PasskeyAssertionOutput, error) {
	var err error

	input := types.AssertionInput{
		CredentialID:       types.CredentialID(assert.CredentialID),
		RawAssertionObject: hex.Hash{},
		AssertionObject: &types.AssertionObject{
			RawAuthenticatorData: assert.RawAuthenticatorData,
			Signature:            assert.RawSignature,
		},
		UserID:            assert.SessionID,
		RawClientDataJSON: assert.UTF8ClientDataJSON,
	}

	cd, err := clientdata.ParseClientData(input.RawClientDataJSON)
	if err != nil {
		return nil, err
	}

	cer, cred, err := store.GetExisting(ctx, cd.Challenge, input.CredentialID)
	if err != nil {
		return nil, terrors.Wrap(err, "failed to get existing ceremony").WithCode(502)
	}

	// Handle steps 4 through 16
	if validError := assertion.VerifyAssertionInput(ctx, &types.VerifyAssertionInputArgs{
		Input:                          input,
		StoredChallenge:                cer.ChallengeID,
		RelyingPartyID:                 rp.RPID(),
		RelyingPartyOrigin:             rp.RPOrigin(),
		CredentialAttestationType:      types.NotFidoAttestationType,
		AttestationProvider:            providers.NewNoneAttestationProvider(),
		AAGUID:                         cred.AAGUID,
		VerifyUser:                     false,
		CredentialPublicKey:            cred.PublicKey,
		Extensions:                     extensions.ClientInputs{},
		DataSignedByClient:             hex.Hash([]byte(input.RawClientDataJSON)),
		UseSavedAttestedCredentialData: false,
	}); validError != nil {
		return nil, validError
	}

	grp, ctx := errgroup.WithContext(ctx)

	grp.SetLimit(2)

	var tkn string

	grp.Go(func() error {
		tknd, err := tknp.AccessTokenForUserID(ctx, assert.SessionID.Hex())
		if err != nil {
			return terrors.Wrap(err, "failed to generate access token")
		}

		tkn = tknd

		return nil
	})

	grp.Go(func() error {
		// verify the aaguid matches
		// verify the public key matches
		err := store.IncrementExistingCredential(ctx, cd.Challenge, cred)
		if err != nil {
			return terrors.Wrap(err, "failed to increment credential")
		}
		return nil
	})

	err = grp.Wait()
	if err != nil {
		return nil, terrors.Wrap(err, "failed to write data")
	}

	return &PasskeyAssertionOutput{AccessToken: tkn}, nil
}
