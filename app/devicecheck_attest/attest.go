package devicecheck

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/walteh/webauthn/pkg/errd"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/relyingparty"
	"github.com/walteh/webauthn/pkg/storage"
	"github.com/walteh/webauthn/pkg/webauthn/clientdata"
	"github.com/walteh/webauthn/pkg/webauthn/credential"
	"github.com/walteh/webauthn/pkg/webauthn/providers"
	"github.com/walteh/webauthn/pkg/webauthn/types"
)

type DeviceCheckAttestationInput struct {
	RawAttestationObject hex.Hash
	UTF8ClientDataJSON   string
	RawCredentialID      hex.Hash
	RawSessionID         hex.Hash
	Production           bool
	Time                 *time.Time
	RootCert             string
}

type DeviceCheckAttestationOutput struct {
	SuggestedStatusCode int
	OK                  bool
}

var (
	ErrDeviceCheckAttestInvalidInput = errors.New("ErrDeviceCheckAttestInvalidInput")

	ErrDeviceCheckAttestInvalidSessionID = errors.New("ErrDeviceCheckAttestInvalidSessionID")

	ErrDeviceCheckAttestInvalidCredentialID = errors.New("ErrDeviceCheckAttestInvalidCredentialID")

	ErrDeviceCheckAttestInvalidChallenge = errors.New("ErrDeviceCheckAttestInvalidChallenge")

	ErrDeviceCheckAttestInvalidCounter = errors.New("ErrDeviceCheckAttestInvalidCounter")

	ErrDeviceCheckAttestDataRead = errors.New("ErrDeviceCheckAttestDataRead")

	ErrDeviceCheckAttestDataWrite = errors.New("ErrDeviceCheckAttestDataWrite")
)

func Attest(ctx context.Context, dynamoClient storage.Provider, rp relyingparty.Provider, input DeviceCheckAttestationInput) (DeviceCheckAttestationOutput, error) {
	var err error

	if input.RawAttestationObject.IsZero() || input.UTF8ClientDataJSON == "" || input.RawCredentialID.IsZero() {
		return DeviceCheckAttestationOutput{400, false}, errd.Wrap(ctx, ErrDeviceCheckAttestInvalidInput)
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
		return DeviceCheckAttestationOutput{400, false}, errd.Wrap(ctx, ErrDeviceCheckAttestInvalidInput)
	}

	cer, _, err := dynamoClient.GetExisting(ctx, cd.Challenge.String(), "")
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to transact get")
		return DeviceCheckAttestationOutput{502, false}, errd.Wrap(ctx, ErrDeviceCheckAttestDataRead)
	}

	if !cer.SessionID.Equals(input.RawSessionID) {
		return DeviceCheckAttestationOutput{401, false}, errd.Mismatch(ctx, ErrDeviceCheckAttestInvalidSessionID, cer.SessionID.Hex(), input.RawSessionID.Hex())
	}

	prov := providers.NewAppAttestSandbox()
	if input.Production {
		prov = providers.NewAppAttestProduction()
	}

	if input.Time != nil {
		prov = prov.WithTime(*input.Time)
	}

	if input.RootCert != "" {
		prov = prov.WithRootCert(input.RootCert)
	}

	pk, err := credential.VerifyAttestationInput(ctx, types.VerifyAttestationInputArgs{
		Provider:           prov,
		Input:              parsedResponse,
		SessionId:          cer.SessionID,
		StoredChallenge:    cer.ChallengeID,
		VerifyUser:         false,
		RelyingPartyID:     rp.RPID(),
		RelyingPartyOrigin: rp.RPOrigin(),
	})

	if err != nil {
		return DeviceCheckAttestationOutput{401, false}, errd.Wrap(ctx, err)
	}

	if !input.RawCredentialID.Equals(pk.RawID) {
		return DeviceCheckAttestationOutput{401, false}, errd.Mismatch(ctx, ErrDeviceCheckAttestInvalidCredentialID, input.RawCredentialID.Hex(), pk.RawID.Hex())
	}

	err = dynamoClient.WriteNewCredential(ctx, cer, pk)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to write new credential")
		return DeviceCheckAttestationOutput{502, false}, errd.Wrap(ctx, ErrDeviceCheckAttestDataWrite)
	}

	// del, err := dynamo.MakeDelete(dynamoClient.MustCeremonyTableName(), cer)
	// if err != nil {
	// 	zerolog.Ctx(ctx).Error().Err(err).Msg("failed to make delete")
	// 	return DeviceCheckAttestationOutput{502, false}, errd.Wrap(ctx, ErrDeviceCheckAttestDataWrite)
	// }

	// putter, err := dynamo.MakePut(dynamoClient.MustCredentialTableName(), pk)
	// if err != nil {
	// 	zerolog.Ctx(ctx).Error().Err(err).Msg("failed to make put")
	// 	return DeviceCheckAttestationOutput{500, false}, errd.Wrap(ctx, ErrDeviceCheckAttestDataWrite)
	// }

	// err = dynamoClient.TransactWrite(ctx, *putter, *del)
	// if err != nil {
	// 	zerolog.Ctx(ctx).Error().Err(err).Msg("failed to transact write")
	// 	return DeviceCheckAttestationOutput{502, false}, errd.Wrap(ctx, ErrDeviceCheckAttestDataWrite)
	// }

	return DeviceCheckAttestationOutput{204, true}, nil
}
