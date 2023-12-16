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
	SuggestedStatusCode int
	OK                  bool
}

func Assert(ctx context.Context, dynamoClient storage.Provider, rp relyingparty.Provider, input DeviceCheckAssertionInput) (DeviceCheckAssertionOutput, error) {
	var err error

	if input.RawAssertionObject.IsZero() || input.ClientDataToValidate.IsZero() {
		return DeviceCheckAssertionOutput{400, false}, err
	}

	parsed, err := assertion.ParseFidoAssertionInput(ctx, input.RawAssertionObject)
	if err != nil {
		return DeviceCheckAssertionOutput{400, false}, err
	}

	cd, err := clientdata.ParseClientData(parsed.RawClientDataJSON)
	if err != nil {
		return DeviceCheckAssertionOutput{400, false}, err
	}

	cerem, cred, err := dynamoClient.GetExisting(ctx, cd.Challenge, parsed.CredentialID)
	if err != nil {
		return DeviceCheckAssertionOutput{502, false}, err
	}

	if cred.RawID.Ref().Hex() != cerem.CredentialID.Ref().Hex() {
		return DeviceCheckAssertionOutput{401, false}, terrors.Errorf("credential id does not match ceremony id")
	}

	if !cerem.ChallengeID.Ref().Equals(cd.Challenge.Ref()) {
		// err :=
		// zerolog.Ctx(ctx).Error().Err(err).Msg("assertion failed")
		return DeviceCheckAssertionOutput{401, false}, terrors.Errorf("challenge ids do not match")
	}

	attestationProvider := providers.NewAppAttestSandbox()

	if attestationProvider.ID() != cred.AttestationType {
		return DeviceCheckAssertionOutput{401, false}, terrors.Errorf("attestation type does not match")
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
		return DeviceCheckAssertionOutput{401, false}, validError
	}

	err = dynamoClient.IncrementExistingCredential(ctx, cerem.ChallengeID, parsed.CredentialID)
	if err != nil {
		return DeviceCheckAssertionOutput{502, false}, err
	}

	return DeviceCheckAssertionOutput{204, true}, nil
}
