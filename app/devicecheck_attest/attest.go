package devicecheck

import (
	"context"
	"time"

	"github.com/walteh/terrors"
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
	RawCredentialID      types.CredentialID
	RawSessionID         hex.Hash
	Production           bool
	Time                 *time.Time
	RootCert             string
}

type DeviceCheckAttestationOutput struct {
	SuggestedStatusCode int
	OK                  bool
}

func Attest(ctx context.Context, store storage.Provider, rp relyingparty.Provider, input DeviceCheckAttestationInput) (DeviceCheckAttestationOutput, error) {
	var err error

	if input.RawAttestationObject.IsZero() || input.UTF8ClientDataJSON == "" || input.RawCredentialID.Ref().IsZero() {
		return DeviceCheckAttestationOutput{400, false}, terrors.New("invalid input")
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
		return DeviceCheckAttestationOutput{400, false}, terrors.Wrap(err, "failed to parse client data")
	}

	// cer, _, err := store.GetExisting(ctx, cd.Challenge, nil)
	// if err != nil {
	// 	zerolog.Ctx(ctx).Error().Err(err).Msg("failed to transact get")
	// 	return DeviceCheckAttestationOutput{502, false}, terrors.Wrap(err, "failed to transact get")
	// }

	// if !cer.SessionID.Equals(input.RawSessionID) {
	// 	return DeviceCheckAttestationOutput{401, false}, terrors.Mismatch(cer.SessionID.Hex(), input.RawSessionID.Hex())
	// }

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
		SessionId:          input.RawSessionID,
		StoredChallenge:    cd.Challenge,
		VerifyUser:         false,
		RelyingPartyID:     rp.RPID(),
		RelyingPartyOrigin: rp.RPOrigin(),
	})

	if err != nil {
		return DeviceCheckAttestationOutput{401, false}, terrors.Wrap(err, "failed to verify attestation input")
	}

	if !input.RawCredentialID.Ref().Equals(pk.RawID.Ref()) {
		return DeviceCheckAttestationOutput{401, false}, terrors.Mismatch(input.RawCredentialID.Ref().Hex(), pk.RawID.Ref().Hex())
	}

	err = store.WriteNewCredential(ctx, cd.Challenge, pk)
	if err != nil {
		return DeviceCheckAttestationOutput{502, false}, terrors.Wrap(err, "failed to write new credential")
	}

	return DeviceCheckAttestationOutput{204, true}, nil
}
