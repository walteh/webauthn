package devicecheck

import (
	"context"

	"git.nugg.xyz/go-sdk/errors"
	"github.com/rs/zerolog"

	"git.nugg.xyz/webauthn/pkg/constants"
	"git.nugg.xyz/webauthn/pkg/dynamo"
	cerrors "git.nugg.xyz/webauthn/pkg/errors"

	"git.nugg.xyz/webauthn/pkg/hex"
	"git.nugg.xyz/webauthn/pkg/webauthn/assertion"
	"git.nugg.xyz/webauthn/pkg/webauthn/clientdata"
	"git.nugg.xyz/webauthn/pkg/webauthn/extensions"
	"git.nugg.xyz/webauthn/pkg/webauthn/providers"
	"git.nugg.xyz/webauthn/pkg/webauthn/types"
)

type DeviceCheckAssertionInput struct {
	RawAssertionObject   hex.Hash
	ClientDataToValidate hex.Hash
}

type DeviceCheckAssertionOutput struct {
	SuggestedStatusCode int
	OK                  bool
}

func Assert(ctx context.Context, dynamoClient *dynamo.Client, input DeviceCheckAssertionInput) (DeviceCheckAssertionOutput, error) {
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

	cred := types.NewUnsafeGettableCredential(parsed.CredentialID)
	cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	err = dynamoClient.TransactGet(ctx, cred, cerem)
	if err != nil {
		return DeviceCheckAssertionOutput{502, false}, err
	}

	if cred.RawID.Hex() != cerem.CredentialID.Hex() {
		err := errors.New("credential id does not match ceremony id")
		zerolog.Ctx(ctx).Error().Err(err).Msg("assertion failed")
		return DeviceCheckAssertionOutput{401, false}, err
	}

	if !cerem.ChallengeID.Equals(cd.Challenge) {
		err := errors.New("challenge ids do not match")
		zerolog.Ctx(ctx).Error().Err(err).Msg("assertion failed")
		return DeviceCheckAssertionOutput{401, false}, err
	}

	attestationProvider := providers.NewAppAttestSandbox()

	if attestationProvider.ID() != cred.AttestationType {
		return DeviceCheckAssertionOutput{401, false}, cerrors.Err0x67InvalidInput.WithMessage("invalid attestation provider").WithCaller()
	}

	// Handle steps 4 through 16
	if validError := assertion.VerifyAssertionInput(ctx, types.VerifyAssertionInputArgs{
		Input:                          parsed,
		StoredChallenge:                cerem.ChallengeID,
		RelyingPartyID:                 "4497QJSAD3.xyz.nugg.app",
		RelyingPartyOrigin:             constants.RPOrigin(),
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

	ceremonyDelete, err := dynamoClient.BuildDelete(cerem)
	if err != nil {
		return DeviceCheckAssertionOutput{500, false}, err
	}

	credentialUpdate, err := cred.UpdateIncreasingCounter(dynamoClient.MustCredentialTableName())
	if err != nil {
		return DeviceCheckAssertionOutput{500, false}, err
	}

	err = dynamoClient.TransactWrite(ctx, *credentialUpdate, *ceremonyDelete)
	if err != nil {
		return DeviceCheckAssertionOutput{502, false}, err
	}

	return DeviceCheckAssertionOutput{204, true}, nil
}
