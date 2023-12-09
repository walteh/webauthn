package passkey_assert

import (
	"context"

	"github.com/walteh/webauthn/pkg/accesstoken/cognito"
	"github.com/walteh/webauthn/pkg/errd"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/relyingparty"
	"github.com/walteh/webauthn/pkg/storage"
	"github.com/walteh/webauthn/pkg/webauthn/assertion"
	"github.com/walteh/webauthn/pkg/webauthn/clientdata"
	"github.com/walteh/webauthn/pkg/webauthn/extensions"
	"github.com/walteh/webauthn/pkg/webauthn/providers"
	"github.com/walteh/webauthn/pkg/webauthn/types"
)

type PasskeyAssertionInput struct {
	SessionID            hex.Hash `json:"userID"`
	CredentialID         hex.Hash `json:"credentialID"`
	UTF8ClientDataJSON   string   `json:"rawClientDataJSON"`
	RawAuthenticatorData hex.Hash `json:"rawAuthenticatorData"`
	RawSignature         hex.Hash `json:"signature"`
	PublicKey            hex.Hash `json:"publicKey"`
	AAGUID               hex.Hash `json:"aaguid"`
}

type PasskeyAssertionOutput struct {
	SuggestedStatusCode int
	AccessToken         string
}

func Assert(ctx context.Context, dynamoClient storage.Provider, rp relyingparty.Provider, cognitoClient cognito.Client, assert PasskeyAssertionInput) (PasskeyAssertionOutput, error) {
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

	// cred := structure.NewCredentialQueryable(input.CredentialID.Hex())
	// cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	// if err = dynamoClient.TransactGet(ctx, cred, cerem); err != nil {
	// 	return PasskeyAssertionOutput{502, ""}, err
	// }

	// if cred.RawID.Hex() != cerem.CredentialID.Hex() {
	// 	log.Println(cred.RawID.Hex(), cerem.CredentialID.Hex())
	// 	return PasskeyAssertionOutput{401, ""}, errors.NewError(0x67).WithMessage("credential id does not match").WithCaller()
	// }

	z, err := cognitoClient.GetDevCreds(ctx, input.CredentialID)
	if err != nil {
		return PasskeyAssertionOutput{502, ""}, errd.Wrap(ctx, err)
	}

	// Handle steps 4 through 16
	if validError := assertion.VerifyAssertionInput(ctx, types.VerifyAssertionInputArgs{
		Input:                          input,
		StoredChallenge:                cd.Challenge,
		RelyingPartyID:                 rp.RPID(),
		RelyingPartyOrigin:             rp.RPOrigin(),
		CredentialAttestationType:      types.NotFidoAttestationType,
		AttestationProvider:            providers.NewNoneAttestationProvider(),
		AAGUID:                         assert.AAGUID,
		VerifyUser:                     false,
		CredentialPublicKey:            assert.PublicKey,
		Extensions:                     extensions.ClientInputs{},
		DataSignedByClient:             hex.Hash([]byte(input.RawClientDataJSON)),
		UseSavedAttestedCredentialData: false,
	}); validError != nil {
		return PasskeyAssertionOutput{401, ""}, validError
	}

	// verify the aaguid matches
	// verify the public key matches

	err = dynamoClient.IncrementExistingCredential(ctx, cd.Challenge.Hex(), input.CredentialID.Hex())
	if err != nil {
		return PasskeyAssertionOutput{502, ""}, err
	}

	// // add session to list of sessions for this credential
	// txs := indexable.IndexableIncrement(ctx, cred, indexable.NewCustomLastModifier(0, false), indexable.N(1))

	// // make sure the ceremony has not expired
	// // we should have some sort of ttl on it - need to make sure of that
	// // in our indexable increment or append of the session to a list, we will validate the ceremony has not been used
	// // tx2 := x.Ind

	// credentialUpdate, err := cred.UpdateIncreasingCounter(dynamoClient.MustCredentialTableName())
	// if err != nil {
	// 	return PasskeyAssertionOutput{500, ""}, err
	// }

	// ceremonyDelete, err := dynamoClient.BuildDelete(cerem)
	// if err != nil {
	// 	return PasskeyAssertionOutput{500, ""}, err
	// }

	// err = dynamoClient.TransactWrite(ctx, *credentialUpdate, *ceremonyDelete)
	// if err != nil {
	// 	return PasskeyAssertionOutput{502, ""}, err
	// }

	return PasskeyAssertionOutput{204, *z.Token}, nil
}
