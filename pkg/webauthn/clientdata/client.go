package clientdata

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/walteh/terrors"
	"github.com/walteh/webauthn/pkg/errd"
	"github.com/walteh/webauthn/pkg/webauthn/types"
)

func ParseClientData(clientData string) (types.CollectedClientData, error) {
	var cd types.CollectedClientData
	err := json.Unmarshal([]byte(clientData), &cd)
	if err != nil {
		log.Printf("failed to unmarshal client data, %v", err)
		return types.CollectedClientData{}, err
	}
	return cd, nil
}

// Handles steps 3 through 6 of verfying the registering client data of a
// new credential and steps 7 through 10 of verifying an authentication assertion
// See https://www.w3.org/TR/webauthn/#registering-a-new-credential
// and https://www.w3.org/TR/webauthn/#verifying-assertion
func Verify(ctx context.Context, expected types.VerifyClientDataArgs) error {

	r := expected.ClientData
	// Registration Step 3. Verify that the value of C.type is webauthn.create.

	// Assertion Step 7. Verify that the value of C.type is the string webauthn.get.
	if r.Type != expected.CeremonyType {
		return errd.Wrap(ctx, ErrInvalidCeremonyType)
	}

	// Registration Step 4. Verify that the value of C.challenge matches the challenge
	// that was sent to the authenticator in the create() call.

	// Assertion Step 8. Verify that the value of C.challenge matches the challenge
	// that was sent to the authenticator in the PublicKeyCredentialRequestOptions
	// passed to the get() call.

	// challenge := r.Challenge

	// rdata, err := base64.RawURLEncoding.DecodeString(challenge)
	// if err != nil {
	// 	return err
	// }

	// abc := base64.RawURLEncoding.EncodeToString(storedChallenge)

	// log.Println(abc)

	if subtle.ConstantTimeCompare(expected.StoredChallenge, r.Challenge) != 1 {
		return terrors.Mismatch(expected.StoredChallenge.Ref().Hex(), r.Challenge.Ref().Hex())
		// return errd.Mismatch(ctx, ErrChallengeMismatch, expected.StoredChallenge.Ref().Hex(), string(r.Challenge.Ref().Hex()))
	}

	// Registration Step 5 & Assertion Step 9. Verify that the value of C.origin matches
	// the Relying Party's origin.
	clientDataOrigin, err := url.Parse(r.Origin)
	if err != nil {
		return errd.Wrap(ctx, ErrOriginNotParsableAsURL)
	}

	if !strings.EqualFold(types.FullyQualifiedOrigin(clientDataOrigin), expected.RelyingPartyOrigin) {
		return terrors.Mismatch(expected.RelyingPartyOrigin, types.FullyQualifiedOrigin(clientDataOrigin))
		// return terrors.Mismatch(ctx, ErrOriginMismatch, expected.RelyingPartyOrigin, types.FullyQualifiedOrigin(clientDataOrigin))
	}

	// Registration Step 6 and Assertion Step 10. Verify that the value of C.tokenBinding.status
	// matches the state of Token Binding for the TLS connection over which the assertion was
	// obtained. If Token Binding was used on that TLS connection, also verify that C.tokenBinding.id
	// matches the base64url encoding of the Token Binding ID for the connection.
	if r.TokenBinding != nil {
		if r.TokenBinding.Status == "" {
			return errd.Wrap(ctx, ErrTokenMissingStatus)
		}
		if r.TokenBinding.Status != types.Present && r.TokenBinding.Status != types.Supported && r.TokenBinding.Status != types.NotSupported {
			return errd.Wrap(ctx, ErrTokenInvalidStatus)
		}
	}
	// Not yet fully implemented by the spec, browsers, and me.

	return nil
}
