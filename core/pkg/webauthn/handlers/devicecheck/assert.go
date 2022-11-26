package devicecheck

import (
	"context"
	"log"
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/env"
	"nugg-webauthn/core/pkg/errors"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/webauthn/assertion"
	"nugg-webauthn/core/pkg/webauthn/clientdata"
	"nugg-webauthn/core/pkg/webauthn/extensions"
	"nugg-webauthn/core/pkg/webauthn/providers"
	"nugg-webauthn/core/pkg/webauthn/types"
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
		return DeviceCheckAssertionOutput{400, false}, errors.Err0x67InvalidInput.WithCaller()
	}

	parsed, err := assertion.ParseFidoAssertionInput(input.RawAssertionObject)
	if err != nil {
		return DeviceCheckAssertionOutput{400, false}, errors.Err0x67InvalidInput.WithCaller()
	}

	cd, err := clientdata.ParseClientData(parsed.RawClientDataJSON)
	if err != nil {
		return DeviceCheckAssertionOutput{400, false}, errors.Err0x67InvalidInput.WithCaller()
	}

	cred := types.NewUnsafeGettableCredential(parsed.CredentialID)
	cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	err = dynamoClient.TransactGet(ctx, cred, cerem)
	if err != nil {
		return DeviceCheckAssertionOutput{502, false}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	if cred.RawID.Hex() != cerem.CredentialID.Hex() {
		log.Println(cred.RawID.Hex(), cerem.CredentialID.Hex())
		return DeviceCheckAssertionOutput{401, false}, errors.Err0x67InvalidInput.WithMessage("credential ids do not match").WithCaller()
	}

	if !cerem.ChallengeID.Equals(cd.Challenge) {
		return DeviceCheckAssertionOutput{401, false}, errors.Err0x67InvalidInput.WithMessage("challenge ids do not match").WithCaller()
	}

	attestationProvider := providers.NewAppAttestSandbox()

	if attestationProvider.ID() != cred.AttestationType {
		return DeviceCheckAssertionOutput{401, false}, errors.Err0x67InvalidInput.WithMessage("invalid attestation provider").WithCaller()
	}

	// Handle steps 4 through 16
	if validError := assertion.VerifyAssertionInput(types.VerifyAssertionInputArgs{
		Input:                          parsed,
		StoredChallenge:                cerem.ChallengeID,
		RelyingPartyID:                 "4497QJSAD3.xyz.nugg.app",
		RelyingPartyOrigin:             env.RPOrigin(),
		AAGUID:                         cred.AAGUID,
		CredentialAttestationType:      types.FidoAttestationType,
		AttestationProvider:            attestationProvider,
		VerifyUser:                     false,
		CredentialPublicKey:            cred.PublicKey,
		Extensions:                     extensions.ClientInputs{},
		DataSignedByClient:             append(input.ClientDataToValidate, cerem.ChallengeID...),
		UseSavedAttestedCredentialData: true,
	}); validError != nil {
		return DeviceCheckAssertionOutput{401, false}, errors.Err0x67InvalidInput.WithMessage("invalid assertion").WithCaller()
	}

	ceremonyDelete, err := dynamoClient.BuildDelete(cerem)
	if err != nil {
		return DeviceCheckAssertionOutput{500, false}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	credentialUpdate, err := cred.UpdateIncreasingCounter(dynamoClient.MustCredentialTableName())
	if err != nil {
		return DeviceCheckAssertionOutput{500, false}, errors.NewError(0x99).WithMessage("problem creating dynamo put").WithRoot(err).WithCaller()
	}

	err = dynamoClient.TransactWrite(ctx, *credentialUpdate, *ceremonyDelete)
	if err != nil {
		return DeviceCheckAssertionOutput{502, false}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	return DeviceCheckAssertionOutput{204, true}, nil
}