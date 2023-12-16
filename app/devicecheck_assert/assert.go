package devicecheck_assert

import (
	"context"

	"github.com/walteh/webauthn/pkg/relyingparty"
	"github.com/walteh/webauthn/pkg/storage"

	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/assertion"
	"github.com/walteh/webauthn/pkg/webauthn/clientdata"
	"github.com/walteh/webauthn/pkg/webauthn/extensions"
	"github.com/walteh/webauthn/pkg/webauthn/providers"
	"github.com/walteh/webauthn/pkg/webauthn/types"

	"github.com/walteh/terrors"
)

type DeviceCheckAssertionInput struct {
	RawAssertionObject   hex.Hash
	ClientDataToValidate hex.Hash
}

type DeviceCheckAssertionOutput struct {
	OK bool
}

func Assert(ctx context.Context, store storage.Provider, rp relyingparty.Provider, input *DeviceCheckAssertionInput) (*DeviceCheckAssertionOutput, error) {
	var err error

	if input.RawAssertionObject.IsZero() || input.ClientDataToValidate.IsZero() {
		return nil, terrors.New("invalid input").WithCode(400)
	}

	parsed, err := assertion.ParseFidoAssertionInput(ctx, input.RawAssertionObject)
	if err != nil {
		return nil, terrors.Wrap(err, "failed to parse assertion input").WithCode(400)
	}

	cd, err := clientdata.ParseClientData(parsed.RawClientDataJSON)
	if err != nil {
		return nil, terrors.Wrap(err, "failed to parse client data").WithCode(400)
	}

	cerem, cred, err := store.GetExisting(ctx, cd.Challenge, parsed.CredentialID)
	if err != nil {
		return nil, terrors.Wrap(err, "failed to get existing ceremony").WithCode(502)
	}

	if cred.RawID.Ref().Hex() != cerem.CredentialID.Ref().Hex() {
		return nil, terrors.Errorf("credential id does not match ceremony id").WithCode(401)
	}

	if !cerem.ChallengeID.Ref().Equals(cd.Challenge.Ref()) {
		// err :=
		// zerolog.Ctx(ctx).Error().Err(err).Msg("assertion failed")
		return nil, terrors.Errorf("challenge ids do not match").WithCode(401)
	}

	attestationProvider := providers.NewAppAttestSandbox()

	if attestationProvider.ID() != cred.AttestationType {
		return nil, terrors.Errorf("attestation type does not match").WithCode(401)
	}

	// Handle steps 4 through 16
	if validError := assertion.VerifyAssertionInput(ctx, &types.VerifyAssertionInputArgs{
		Input:                          parsed,
		StoredChallenge:                cerem.ChallengeID,
		RelyingPartyID:                 rp.RPID(),
		RelyingPartyOrigin:             rp.RPOrigin(),
		AAGUID:                         cred.AAGUID,
		CredentialAttestationType:      types.FidoAttestationType,
		AttestationProvider:            attestationProvider,
		VerifyUser:                     false,
		CredentialPublicKey:            cred.PublicKey,
		Extensions:                     extensions.ClientInputs{},
		DataSignedByClient:             append(input.ClientDataToValidate, cerem.ChallengeID...),
		UseSavedAttestedCredentialData: true,
	}); validError != nil {
		return nil, terrors.Wrap(validError, "failed to verify assertion input").WithCode(401)
	}

	err = store.IncrementExistingCredential(ctx, cerem.ChallengeID, parsed.CredentialID)
	if err != nil {
		return nil, terrors.Wrap(err, "failed to increment existing credential").WithCode(502)
	}

	return &DeviceCheckAssertionOutput{true}, nil
}
