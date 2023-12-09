package passkey_attest

import (
	"context"
	"errors"

	"github.com/walteh/webauthn/pkg/accesstoken"
	"github.com/walteh/webauthn/pkg/errd"
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
}

type PasskeyAttestationOutput struct {
	SuggestedStatusCode int
	AccessToken         string
}

var (
	ErrPasskeyAttestInvalidInput = errors.New("ErrPasskeyAttestInvalidInput")

	ErrPasskeyAttestInvalidSessionID = errors.New("ErrPasskeyAttestInvalidSessionID")

	ErrPasskeyAttestInvalidCredentialID = errors.New("ErrPasskeyAttestInvalidCredentialID")

	ErrPasskeyAttestInvalidChallenge = errors.New("ErrPasskeyAttestInvalidChallenge")

	ErrPasskeyAttestJWTGeneration = errors.New("ErrPasskeyAttestJWTGeneration")

	ErrPasskeyAttestDataRead = errors.New("ErrPasskeyAttestDataRead")

	ErrPasskeyAttestDataWrite = errors.New("ErrPasskeyAttestDataWrite")
)

func Attest(ctx context.Context, dynamoClient storage.Provider, rp relyingparty.Provider, tknp accesstoken.Provider, assert PasskeyAttestationInput) (PasskeyAttestationOutput, error) {
	var err error

	parsedResponse := types.AttestationInput{
		AttestationObject:  assert.RawAttestationObject,
		UTF8ClientDataJSON: assert.UTF8ClientDataJSON,
		CredentialID:       assert.RawCredentialID,
		ClientExtensions:   nil,
	}

	cd, err := clientdata.ParseClientData(parsedResponse.UTF8ClientDataJSON)
	if err != nil {
		return PasskeyAttestationOutput{400, ""}, err
	}

	// cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	cerem, err := dynamoClient.GetExistingCeremony(ctx, cd.Challenge.String())
	if err != nil {
		return PasskeyAttestationOutput{502, ""}, errd.Wrap(ctx, ErrPasskeyAttestDataRead)
	}

	cred, invalidErr := credential.VerifyAttestationInput(ctx, types.VerifyAttestationInputArgs{
		Provider:           providers.NewNoneAttestationProvider(),
		Input:              parsedResponse,
		StoredChallenge:    cerem.ChallengeID,
		SessionId:          cerem.SessionID,
		VerifyUser:         false,
		RelyingPartyID:     rp.RPID(),
		RelyingPartyOrigin: rp.RPOrigin(),
	})

	if invalidErr != nil {
		return PasskeyAttestationOutput{401, ""}, invalidErr
	}

	tkn, err := tknp.AccessTokenForUserID(ctx, cerem.CredentialID.String())
	if err != nil {
		return PasskeyAttestationOutput{502, ""}, errd.Wrap(ctx, ErrPasskeyAttestJWTGeneration)
	}

	// z, err := cognitoClient.GetDevCreds(ctx, cerem.CredentialID)
	// if err != nil {
	// 	return PasskeyAttestationOutput{502, ""}, errd.Wrap(ctx, ErrPasskeyAttestJWTGeneration)
	// }

	// // put
	// credput := indexable.IndexablePut(cred, false)

	// // should be a delete
	// ceremput := indexable.IndexablePut(cerem, true)

	err = dynamoClient.IncrementExistingCredential(ctx, string(cerem.ChallengeID), string(cred.RawID))

	if err != nil {
		return PasskeyAttestationOutput{502, ""}, errd.Wrap(ctx, ErrPasskeyAttestDataWrite)
	}

	return PasskeyAttestationOutput{204, tkn}, nil
}
