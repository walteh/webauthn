package passkey

import (
	"context"
	"log"
	"nugg-webauthn/core/pkg/cognito"
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

type PasskeyAssertionInput struct {
	SessionID            hex.Hash `json:"userID"`
	CredentialID         hex.Hash `json:"credentialID"`
	UTF8ClientDataJSON   string   `json:"rawClientDataJSON"`
	RawAuthenticatorData hex.Hash `json:"rawAuthenticatorData"`
	RawSignature         hex.Hash `json:"signature"`
}

func Assert(ctx context.Context, dynamoClient *dynamo.Client, cognitoClient cognito.Client, assert PasskeyAssertionInput, body hex.Hash) (int, string, error) {
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
		return 400, "", errors.NewError(0x67).WithMessage("invalid client data").WithRoot(err).WithCaller()
	}

	cred := types.NewUnsafeGettableCredential(input.CredentialID)
	cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	if err = dynamoClient.TransactGet(ctx, cred, cerem); err != nil {
		return 502, "", errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	if cred.RawID.Hex() != cerem.CredentialID.Hex() {
		log.Println(cred.RawID.Hex(), cerem.CredentialID.Hex())
		return 401, "", errors.NewError(0x67).WithMessage("credential id does not match").WithCaller()
	}

	z, err := cognitoClient.GetDevCreds(ctx, cerem.CredentialID)
	if err != nil {
		return 502, "", errors.NewError(0x99).WithMessage("problem calling cognito").WithRoot(err).WithCaller()
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
		return 401, "", errors.NewError(0x67).WithMessage("invalid assertion").WithRoot(validError).WithCaller()
	}

	credentialUpdate, err := cred.UpdateIncreasingCounter(dynamoClient.MustCredentialTableName())
	if err != nil {
		return 500, "", errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	err = dynamoClient.TransactWrite(ctx, *credentialUpdate)
	if err != nil {
		return 502, "", errors.NewError(0x99).WithMessage("problem calling dynamo").WithRoot(err).WithCaller()
	}

	return 200, *z.Token, nil
}
