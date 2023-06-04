package passkey

import (
	"context"

	"git.nugg.xyz/go-sdk/errors"
	"git.nugg.xyz/go-sdk/x"

	"git.nugg.xyz/webauthn/pkg/cognito"
	"git.nugg.xyz/webauthn/pkg/hex"
	"git.nugg.xyz/webauthn/pkg/webauthn/clientdata"
	"git.nugg.xyz/webauthn/pkg/webauthn/credential"
	"git.nugg.xyz/webauthn/pkg/webauthn/providers"
	"git.nugg.xyz/webauthn/pkg/webauthn/types"
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

func Attest(ctx context.Context, dynamoClient x.DynamoDBAPI, cognitoClient cognito.Client, assert PasskeyAttestationInput) (PasskeyAttestationOutput, error) {
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

	cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	err = dynamoClient.TransactGet(ctx, cerem)
	if err != nil {
		return PasskeyAttestationOutput{502, ""}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	cred, invalidErr := credential.VerifyAttestationInput(ctx, types.VerifyAttestationInputArgs{
		Provider:           providers.NewNoneAttestationProvider(),
		Input:              parsedResponse,
		StoredChallenge:    cerem.ChallengeID,
		SessionId:          cerem.SessionID,
		VerifyUser:         false,
		RelyingPartyID:     "nugg.xyz",
		RelyingPartyOrigin: "https://nugg.xyz",
	})

	if invalidErr != nil {
		return PasskeyAttestationOutput{401, ""}, errors.NewError(0x99).WithMessage("invalid attestation").WithRoot(invalidErr).WithCaller()
	}

	z, err := cognitoClient.GetDevCreds(ctx, cerem.CredentialID)
	if err != nil {
		return PasskeyAttestationOutput{502, ""}, errors.NewError(0x99).WithMessage("problem calling cognito").WithRoot(err).WithCaller()
	}

	// put
	credput := x.IndexablePut(cred, false)

	// should be a delete
	ceremput := x.IndexablePut(cerem, true)

	x.N

	err = dynamoClient.TransactWrite(ctx, *credput, *ceremput)
	if err != nil {
		return PasskeyAttestationOutput{502, ""}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	return PasskeyAttestationOutput{204, *z.Token}, nil
}
