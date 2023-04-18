package passkey

import (
	"context"

	"git.nugg.xyz/go-sdk/dynamo"
	"git.nugg.xyz/go-sdk/errors"
	"git.nugg.xyz/go-sdk/x"

	"git.nugg.xyz/webauthn/pkg/cognito"
	"git.nugg.xyz/webauthn/pkg/constants"
	"git.nugg.xyz/webauthn/pkg/hex"
	"git.nugg.xyz/webauthn/pkg/structure"
	"git.nugg.xyz/webauthn/pkg/webauthn/assertion"
	"git.nugg.xyz/webauthn/pkg/webauthn/clientdata"
	"git.nugg.xyz/webauthn/pkg/webauthn/extensions"
	"git.nugg.xyz/webauthn/pkg/webauthn/providers"
	"git.nugg.xyz/webauthn/pkg/webauthn/types"
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

func Assert(ctx context.Context, dynamoClient x.DynamoDBAPI, cognitoClient cognito.Client, assert PasskeyAssertionInput) (PasskeyAssertionOutput, error) {
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
		return PasskeyAssertionOutput{400, ""}, err
	}

	cred := structure.NewCredentialQueryable(input.CredentialID.Hex())
	cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	// if err = dynamoClient.TransactGet(ctx, cred, cerem); err != nil {
	// 	return PasskeyAssertionOutput{502, ""}, err
	// }

	// if cred.RawID.Hex() != cerem.CredentialID.Hex() {
	// 	log.Println(cred.RawID.Hex(), cerem.CredentialID.Hex())
	// 	return PasskeyAssertionOutput{401, ""}, errors.NewError(0x67).WithMessage("credential id does not match").WithCaller()
	// }

	z, err := cognitoClient.GetDevCreds(ctx, input.CredentialID)
	if err != nil {
		return PasskeyAssertionOutput{502, ""}, errors.NewError(0x99).WithMessage("problem calling cognito").WithRoot(err).WithCaller()
	}

	// Handle steps 4 through 16
	if validError := assertion.VerifyAssertionInput(ctx, types.VerifyAssertionInputArgs{
		Input:                          input,
		StoredChallenge:                cerem.ChallengeID,
		RelyingPartyID:                 constants.RPID(),
		RelyingPartyOrigin:             constants.RPOrigin(),
		CredentialAttestationType:      types.NotFidoAttestationType,
		AttestationProvider:            providers.NewNoneAttestationProvider(),
		AAGUID:                         cred.AAGUID,
		VerifyUser:                     false,
		CredentialPublicKey:            cred.PublicKey,
		Extensions:                     extensions.ClientInputs{},
		DataSignedByClient:             hex.Hash([]byte(input.RawClientDataJSON)),
		UseSavedAttestedCredentialData: false,
	}); validError != nil {
		return PasskeyAssertionOutput{401, ""}, validError
	}

	txs := x.IndexableIncrement(ctx, cred, x.NewCustomLastModifier(0, false), dynamo.N(1))

	credentialUpdate, err := cred.UpdateIncreasingCounter(dynamoClient.MustCredentialTableName())
	if err != nil {
		return PasskeyAssertionOutput{500, ""}, err
	}

	ceremonyDelete, err := dynamoClient.BuildDelete(cerem)
	if err != nil {
		return PasskeyAssertionOutput{500, ""}, err
	}

	err = dynamoClient.TransactWrite(ctx, *credentialUpdate, *ceremonyDelete)
	if err != nil {
		return PasskeyAssertionOutput{502, ""}, err
	}

	return PasskeyAssertionOutput{204, *z.Token}, nil
}
