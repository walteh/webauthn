package devicecheck

import (
	"context"

	"git.nugg.xyz/go-sdk/errors"

	"git.nugg.xyz/webauthn/pkg/dynamo"
	cerrors "git.nugg.xyz/webauthn/pkg/errors"
	"git.nugg.xyz/webauthn/pkg/hex"
	"git.nugg.xyz/webauthn/pkg/webauthn/clientdata"
	"git.nugg.xyz/webauthn/pkg/webauthn/credential"
	"git.nugg.xyz/webauthn/pkg/webauthn/providers"
	"git.nugg.xyz/webauthn/pkg/webauthn/types"
)

type DeviceCheckAttestationInput struct {
	RawAttestationObject hex.Hash
	UTF8ClientDataJSON   string
	RawCredentialID      hex.Hash
	RawSessionID         hex.Hash
}

type DeviceCheckAttestationOutput struct {
	SuggestedStatusCode int
	OK                  bool
}

func Attest(ctx context.Context, dynamoClient *dynamo.Client, input DeviceCheckAttestationInput) (DeviceCheckAttestationOutput, error) {
	var err error

	if input.RawAttestationObject.IsZero() || input.UTF8ClientDataJSON == "" || input.RawCredentialID.IsZero() {
		return DeviceCheckAttestationOutput{400, false}, cerrors.Err0x67InvalidInput.WithCaller()
	}
	parsedResponse := types.AttestationInput{
		AttestationObject:  input.RawAttestationObject,
		UTF8ClientDataJSON: input.UTF8ClientDataJSON,
		CredentialID:       input.RawCredentialID,
		CredentialType:     types.PublicKeyCredentialType,
		ClientExtensions:   nil,
	}

	cd, err := clientdata.ParseClientData(parsedResponse.UTF8ClientDataJSON)
	if err != nil {
		return DeviceCheckAttestationOutput{400, false}, cerrors.Err0x67InvalidInput.WithCaller()
	}

	cer := types.NewUnsafeGettableCeremony(cd.Challenge)

	err = dynamoClient.TransactGet(ctx, cer)
	if err != nil {
		return DeviceCheckAttestationOutput{502, false}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	if !cer.SessionID.Equals(input.RawSessionID) {
		return DeviceCheckAttestationOutput{401, false}, errors.NewError(0x99).WithMessage("session id mismatch").WithCaller()
	}

	pk, err := credential.VerifyAttestationInput(ctx, types.VerifyAttestationInputArgs{
		Provider:           providers.NewAppAttestSandbox(),
		Input:              parsedResponse,
		SessionId:          cer.SessionID,
		StoredChallenge:    cer.ChallengeID,
		VerifyUser:         false,
		RelyingPartyID:     "4497QJSAD3.xyz.nugg.app",
		RelyingPartyOrigin: "https://nugg.xyz",
	})

	if err != nil {
		return DeviceCheckAttestationOutput{401, false}, errors.NewError(0x99).WithMessage("problem verifying attestation").WithRoot(err).WithCaller()
	}

	if !input.RawCredentialID.Equals(pk.RawID) {
		err := errors.NewError(0x93).WithCaller().WithKV("attestationKey", input.RawCredentialID.Hex()).WithKV("pk.RawID", pk.RawID.Hex())
		return DeviceCheckAttestationOutput{401, false}, err
	}

	del, err := dynamo.MakeDelete(dynamoClient.MustCeremonyTableName(), cer)
	if err != nil {
		return DeviceCheckAttestationOutput{502, false}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	putter, err := dynamo.MakePut(dynamoClient.MustCredentialTableName(), pk)
	if err != nil {
		return DeviceCheckAttestationOutput{500, false}, errors.NewError(0x99).WithMessage("problem making put").WithRoot(err).WithCaller()
	}

	err = dynamoClient.TransactWrite(ctx, *putter, *del)
	if err != nil {
		return DeviceCheckAttestationOutput{502, false}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	return DeviceCheckAttestationOutput{204, true}, nil
}
