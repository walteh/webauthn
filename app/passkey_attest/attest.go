package passkey_attest

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/walteh/terrors"
	"github.com/walteh/webauthn/pkg/accesstoken"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/relyingparty"
	"github.com/walteh/webauthn/pkg/storage"
	"github.com/walteh/webauthn/pkg/webauthn/clientdata"
	"github.com/walteh/webauthn/pkg/webauthn/credential"
	"github.com/walteh/webauthn/pkg/webauthn/providers"
	"github.com/walteh/webauthn/pkg/webauthn/types"
)

type PasskeyAttestationInput struct {
	RawAttestationObject hex.Hash
	UTF8ClientDataJSON   string
	RawCredentialID      hex.Hash
	RawSessionID         hex.Hash
}

type PasskeyAttestationOutput struct {
	AccessToken string
}

func Attest(ctx context.Context, store storage.Provider, rp relyingparty.Provider, tknp accesstoken.Provider, assert *PasskeyAttestationInput) (*PasskeyAttestationOutput, error) {
	var err error

	parsedResponse := types.AttestationInput{
		AttestationObject:  assert.RawAttestationObject,
		UTF8ClientDataJSON: assert.UTF8ClientDataJSON,
		CredentialID:       types.CredentialID(assert.RawCredentialID),
		ClientExtensions:   nil,
	}

	cd, err := clientdata.ParseClientData(parsedResponse.UTF8ClientDataJSON)
	if err != nil {
		return nil, terrors.Wrap(err, "failed to parse client data").WithCode(400)
	}

	// cerem, _, err := store.GetExisting(ctx, cd.Challenge, nil)
	// if err != nil {
	// 	return nil, terrors.Wrap(err, "failed to get existing ceremony").WithCode(502)
	// }

	cred, invalidErr := credential.VerifyAttestationInput(ctx, types.VerifyAttestationInputArgs{
		Provider:           providers.NewNoneAttestationProvider(),
		Input:              parsedResponse,
		StoredChallenge:    cd.Challenge,
		SessionId:          assert.RawSessionID,
		VerifyUser:         false,
		RelyingPartyID:     rp.RPID(),
		RelyingPartyOrigin: rp.RPOrigin(),
	})

	if invalidErr != nil {
		return nil, terrors.Wrap(invalidErr, "failed to verify attestation input").WithCode(401)
	}

	// make error group

	grp, ctx := errgroup.WithContext(ctx)

	grp.SetLimit(2)

	var tkn string

	grp.Go(func() error {
		tknd, err := tknp.AccessTokenForUserID(ctx, assert.RawSessionID.Hex())
		if err != nil {
			return terrors.Wrap(err, "failed to generate access token").WithCode(502)
		}

		tkn = tknd

		return nil
	})

	grp.Go(func() error {

		err = store.WriteNewCredential(ctx, cd.Challenge, cred)
		if err != nil {
			return terrors.Wrap(err, "failed to write new credential")
		}
		return nil
	})

	err = grp.Wait()
	if err != nil {
		return nil, err
	}

	return &PasskeyAttestationOutput{AccessToken: tkn}, nil
}
