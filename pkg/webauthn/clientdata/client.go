package clientdata

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"git.nugg.xyz/webauthn/pkg/errors"
	"git.nugg.xyz/webauthn/pkg/webauthn/types"
	"github.com/rs/zerolog"
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

	real := expected.ClientData
	// Registration Step 3. Verify that the value of C.type is webauthn.create.

	// Assertion Step 7. Verify that the value of C.type is the string webauthn.get.
	if real.Type != expected.CeremonyType {
		err := errors.ErrVerification.WithMessage("Error validating ceremony type")
		zerolog.Ctx(ctx).Error().Err(err).Msg("ceremony type mismatch")
		return err
	}

	// Registration Step 4. Verify that the value of C.challenge matches the challenge
	// that was sent to the authenticator in the create() call.

	// Assertion Step 8. Verify that the value of C.challenge matches the challenge
	// that was sent to the authenticator in the PublicKeyCredentialRequestOptions
	// passed to the get() call.

	// challenge := real.Challenge

	// rdata, err := base64.RawURLEncoding.DecodeString(challenge)
	// if err != nil {
	// 	return err
	// }

	// abc := base64.RawURLEncoding.EncodeToString(storedChallenge)

	// log.Println(abc)

	if subtle.ConstantTimeCompare(expected.StoredChallenge, real.Challenge) != 1 {
		err := errors.ErrVerification.WithMessage("Error validating challenge")
		zerolog.Ctx(ctx).Error().Err(err).Str("expected", string(expected.StoredChallenge)).Str("received", string(real.Challenge)).Msg("challenge mismatch")
		return err
	}

	// Registration Step 5 & Assertion Step 9. Verify that the value of C.origin matches
	// the Relying Party's origin.
	clientDataOrigin, err := url.Parse(real.Origin)
	if err != nil {
		return errors.ErrParsingData.WithMessage("Error decoding clientData origin as URL")
	}

	if !strings.EqualFold(types.FullyQualifiedOrigin(clientDataOrigin), expected.RelyingPartyOrigin) {
		err := errors.ErrVerification.WithMessage("Error validating origin")
		zerolog.Ctx(ctx).Error().Err(err).Str("expected", expected.RelyingPartyOrigin).Str("received", types.FullyQualifiedOrigin(clientDataOrigin)).Msg("origin mismatch")
		return err
	}

	// Registration Step 6 and Assertion Step 10. Verify that the value of C.tokenBinding.status
	// matches the state of Token Binding for the TLS connection over which the assertion was
	// obtained. If Token Binding was used on that TLS connection, also verify that C.tokenBinding.id
	// matches the base64url encoding of the Token Binding ID for the connection.
	if real.TokenBinding != nil {
		if real.TokenBinding.Status == "" {
			return errors.ErrParsingData.WithMessage("Error decoding clientData, token binding present without status")
		}
		if real.TokenBinding.Status != types.Present && real.TokenBinding.Status != types.Supported && real.TokenBinding.Status != types.NotSupported {
			err := errors.ErrParsingData.WithMessage("Error decoding clientData, token binding present with invalid status")
			zerolog.Ctx(ctx).Error().Err(err).Msg("invalid token binding status")
			return err
		}
	}
	// Not yet fully implemented by the spec, browsers, and me.

	return nil
}