package passkey

import (
	"context"
	"log"

	"github.com/nuggxyz/webauthn/pkg/cognito"
	"github.com/nuggxyz/webauthn/pkg/dynamo"
	"github.com/nuggxyz/webauthn/pkg/env"
	"github.com/nuggxyz/webauthn/pkg/errors"
	"github.com/nuggxyz/webauthn/pkg/hex"
	"github.com/nuggxyz/webauthn/pkg/webauthn/assertion"
	"github.com/nuggxyz/webauthn/pkg/webauthn/clientdata"
	"github.com/nuggxyz/webauthn/pkg/webauthn/extensions"
	"github.com/nuggxyz/webauthn/pkg/webauthn/providers"
	"github.com/nuggxyz/webauthn/pkg/webauthn/types"
)

type PasskeyAssertionInput struct {
	SessionID            hex.Hash `json:"userID"`
	CredentialID         hex.Hash `json:"credentialID"`
	UTF8ClientDataJSON   string   `json:"rawClientDataJSON"`
	RawAuthenticatorData hex.Hash `json:"rawAuthenticatorData"`
	RawSignature         hex.Hash `json:"signature"`
}

type PasskeyAssertionOutput struct {
	SuggestedStatusCode int
	AccessToken         string
}

func Assert(ctx context.Context, dynamoClient *dynamo.Client, cognitoClient cognito.Client, assert PasskeyAssertionInput) (PasskeyAssertionOutput, error) {
	var err error

	input := types.AssertionInput{
		CredentialID:       assert.CredentialID,
		RawAssertionObject: hex.Hash{},
		AssertionObject: &types.AssertionObject{
			RawAuthenticatorData: assert.RawAuthenticatorData,
			Signature:            assert.RawSignature,
		},
		UserID:            assert.SessionID,
		RawClientDataJSON: assert.UTF8ClientDataJSON,
	}

	cd, err := clientdata.ParseClientData(input.RawClientDataJSON)
	if err != nil {
		return PasskeyAssertionOutput{400, ""}, errors.NewError(0x67).WithMessage("invalid client data").WithRoot(err).WithCaller()
	}

	cred := types.NewUnsafeGettableCredential(input.CredentialID)
	cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	if err = dynamoClient.TransactGet(ctx, cred, cerem); err != nil {
		return PasskeyAssertionOutput{502, ""}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	if cred.RawID.Hex() != cerem.CredentialID.Hex() {
		log.Println(cred.RawID.Hex(), cerem.CredentialID.Hex())
		return PasskeyAssertionOutput{401, ""}, errors.NewError(0x67).WithMessage("credential id does not match").WithCaller()
	}

	z, err := cognitoClient.GetDevCreds(ctx, cerem.CredentialID)
	if err != nil {
		return PasskeyAssertionOutput{502, ""}, errors.NewError(0x99).WithMessage("problem calling cognito").WithRoot(err).WithCaller()
	}

	// Handle steps 4 through 16
	if validError := assertion.VerifyAssertionInput(types.VerifyAssertionInputArgs{
		Input:                          input,
		StoredChallenge:                cerem.ChallengeID,
		RelyingPartyID:                 env.RPID(),
		RelyingPartyOrigin:             env.RPOrigin(),
		CredentialAttestationType:      types.NotFidoAttestationType,
		AttestationProvider:            providers.NewNoneAttestationProvider(),
		AAGUID:                         cred.AAGUID,
		VerifyUser:                     false,
		CredentialPublicKey:            cred.PublicKey,
		Extensions:                     extensions.ClientInputs{},
		DataSignedByClient:             hex.Hash([]byte(input.RawClientDataJSON)),
		UseSavedAttestedCredentialData: false,
	}); validError != nil {
		return PasskeyAssertionOutput{401, ""}, errors.NewError(0x67).WithMessage("invalid assertion").WithRoot(validError).WithCaller()
	}

	credentialUpdate, err := cred.UpdateIncreasingCounter(dynamoClient.MustCredentialTableName())
	if err != nil {
		return PasskeyAssertionOutput{500, ""}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	ceremonyDelete, err := dynamoClient.BuildDelete(cerem)
	if err != nil {
		return PasskeyAssertionOutput{500, ""}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	err = dynamoClient.TransactWrite(ctx, *credentialUpdate, *ceremonyDelete)
	if err != nil {
		return PasskeyAssertionOutput{502, ""}, errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	return PasskeyAssertionOutput{204, *z.Token}, nil
}
