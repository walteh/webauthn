package assertion

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/json"
	"errors"

	"github.com/go-webauthn/webauthn/protocol/webauthncose"
	"github.com/walteh/terrors"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/authdata"
	"github.com/walteh/webauthn/pkg/webauthn/clientdata"
	"github.com/walteh/webauthn/pkg/webauthn/types"

	"github.com/rs/zerolog"
	"github.com/ugorji/go/codec"
)

// Follow the remaining steps outlined in §7.2 Verifying an authentication assertion
// (https://www.w3.org/TR/webauthn/#verifying-assertion) and return an error if there
// is a failure during each step.
func VerifyAssertionInput(ctx context.Context, args types.VerifyAssertionInputArgs) error {
	// Steps 4 through 6 in verifying the assertion data (https://www.w3.org/TR/webauthn/#verifying-assertion) are
	// "assertive" steps, i.e "Let JSONtext be the result of running UTF-8 decode on the value of cData."
	// We handle these steps in part as we verify but also beforehand
	var (
		key      interface{}
		err      error
		asserter types.AssertionObject
	)

	if args.Input.AssertionObject == nil && !args.Input.RawAssertionObject.IsZero() {
		asserter, err = ParseAssertionObject(ctx, args.Input.RawAssertionObject)
		if err != nil {
			return err
		}
	} else {
		if args.Input.AssertionObject == nil {
			err := errors.New("Assertion object is missing")
			zerolog.Ctx(ctx).Error().Err(err).Send()
			return err
		}
		asserter = *args.Input.AssertionObject
	}

	credId := types.NewDefaultCredentialIdentifier(args.Input.CredentialID)

	appID, err := credId.GetAppID(args.Extensions, args.CredentialAttestationType)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Error getting appID")
		return err
	}

	clientd, err := clientdata.ParseClientData(args.Input.RawClientDataJSON)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Error parsing client data")
		return err
	}

	// Handle steps 7 through 10 of assertion by verifying stored data against the Collected Client Data
	// returned by the authenticator
	validError := clientdata.Verify(ctx, types.VerifyClientDataArgs{
		ClientData:         clientd,
		StoredChallenge:    args.StoredChallenge,
		CeremonyType:       types.AssertCeremony,
		RelyingPartyOrigin: args.RelyingPartyOrigin,
	})
	if validError != nil {
		return validError
	}

	// // Begin Step 11. Verify that the rpIdHash in authData is the SHA-256 hash of the RP ID expected by the RP.
	// rpIDHash := sha256.Sum256([]byte(args.RelyingPartyID))

	// var appIDHash [32]byte
	// if appID != "" {
	// 	appIDHash = sha256.Sum256([]byte(appID))
	// }

	// authdata, err := types.ParseAuthenticatorData(args.Input.RawAuthenticatorData)
	// if err != nil {
	// 	return errors.ErrParsingData.WithMessage("Error parsing authenticator data").WithRoot(err).WithCaller()
	// }

	// Handle steps 11 through 14, verifying the authenticator data.
	// validError = authdata.Verify(rpIDHash[:], appIDHash[:], args.VerifyUser, true, args.LastSignCount)
	validError = authdata.VerifyAuenticatorData(ctx, types.VerifyAuenticatorDataArgs{
		Data:                    asserter.RawAuthenticatorData,
		AppId:                   appID,
		RelyingPartyID:          args.RelyingPartyID,
		RequireUserVerification: args.VerifyUser,
		RequireUserPresence:     false,
		LastSignCount:           args.LastSignCount,
		OptionalAttestedCredentialData: types.AttestedCredentialData{
			CredentialID:        args.Input.CredentialID,
			AAGUID:              args.AAGUID,
			CredentialPublicKey: args.CredentialPublicKey,
		},
		UseSavedAttestedCredentialData: args.UseSavedAttestedCredentialData,
	})
	if validError != nil {
		return validError
	}

	// Step 15. Let hash be the result of computing a hash over the cData using SHA-256.
	// clientDataHash := sha256.Sum256(args.DataSignedByClient.Bytes())

	// combo := append(args.DataSignedByClient, args.StoredChallenge...)

	onemore := sha256.Sum256(args.DataSignedByClient)

	// Step 16. Using the credential public key looked up in step 3, verify that sig is
	// a valid signature over the binary concatenation of authData and hash.

	sigData := append(asserter.RawAuthenticatorData, onemore[:]...)

	if args.AttestationProvider.ID() == "apple-appattest" {
		// Apple's App Attest uses a different signature format
		// https://developer.apple.com/documentation/devicecheck/accessing_and_modifying_per-device_keys

		// The signature is a DER-encoded ECDSA signature with the following structure:
		// SEQUENCE {
		//   r INTEGER,
		//   s INTEGER
		// }

		// key, err = webauthncose.ParseFIDOPublicKey(args.CredentialPublicKey)
		nonce := sha256.Sum256(sigData)

		// 3. Use the public key that you stored from the attestation object to verify that the assertion’s signature is valid for nonce.
		x, y := elliptic.Unmarshal(elliptic.P256(), args.CredentialPublicKey)
		if x == nil {
			return terrors.New("error parsing the public key")
		}

		pubkey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
		nonceHash := sha256.Sum256(nonce[:])

		valid := ecdsa.VerifyASN1(pubkey, nonceHash[:], asserter.Signature)
		if !valid {
			return terrors.New("terror validating the assertion signature")
		}
	} else {
		if appID == "" {
			key, err = webauthncose.ParsePublicKey(args.CredentialPublicKey)
			if err != nil {
				return terrors.Wrap(err, "error parsing the public key")
			}
		} else {
			key, err = webauthncose.ParseFIDOPublicKey(args.CredentialPublicKey)
			if err != nil {
				return terrors.Wrap(err, "error parsing the fido public key")
			}
		}

		valid, err := webauthncose.VerifySignature(key, sigData[:], asserter.Signature)
		if !valid || err != nil {
			return terrors.Wrap(err, "error validating the assertion signature").Event(func(e *zerolog.Event) *zerolog.Event {
				return e.Bool("valid", valid).
					Any("key", key).
					Str("signature", asserter.Signature.Hex()).
					Str("sigData", sigData.Hex()).
					Str("appID", appID)
			})

		}
	}

	return nil
}

type AssertionResponse struct {
	UTF8ClientDataJSON string   `json:"client_data_json"`
	AssertionObject    hex.Hash `json:"assertion_object"`
	SessionID          hex.Hash `json:"session_id"`
	Provider           string   `json:"provider"`
	CredentialID       hex.Hash `json:"credential_id"` // not set to types.CredentialID so it can be decoded from JSON (otherwise we get - illegal base64 data at input byte 64)
}

func ParseNotFidoAssertionInput(ctx context.Context, response []byte) (types.AssertionInput, error) {
	return ParseAssertionInput(ctx, response, types.NotFidoAttestationType)
}

func ParseFidoAssertionInput(ctx context.Context, response []byte) (types.AssertionInput, error) {
	return ParseAssertionInput(ctx, response, types.FidoAttestationType)
}

func ParseAssertionInput(ctx context.Context, response []byte, attestationType types.CredentialAttestationType) (types.AssertionInput, error) {
	var parsed AssertionResponse
	err := json.Unmarshal(response, &parsed)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Error parsing assertion response")
		return types.AssertionInput{}, err
	}

	abc := types.AssertionInput{
		UserID:            parsed.SessionID,
		CredentialID:      types.CredentialID(parsed.CredentialID),
		RawClientDataJSON: parsed.UTF8ClientDataJSON,
		// RawAuthenticatorData: nil,
		RawAssertionObject: parsed.AssertionObject,
		AssertionObject:    nil,
		// Type:            attestationType,
	}

	return abc, nil
}

func ParseAssertionObject(ctx context.Context, input hex.Hash) (types.AssertionObject, error) {
	b := types.AssertionObject{}

	var a struct {
		// AuthenticatorData    AuthenticatorData
		RawAuthenticatorData []byte `json:"authenticatorData"`
		Signature            []byte `json:"signature"`
	}

	cborHandler := codec.CborHandle{}

	// Decode the attestation data with unmarshalled auth data
	err := codec.NewDecoderBytes(input, &cborHandler).Decode(&a)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("input", input.Hex())
		return b, err
	}

	b.RawAuthenticatorData = a.RawAuthenticatorData
	b.Signature = a.Signature

	// ad, err := authdata.ParseAuthenticatorData(a.RawAuthenticatorData)
	// if err != nil {
	// 	return a, terrors.Errorf("error decoding auth data: %v", err)
	// }

	// a.AuthenticatorData = ad

	// if err := json.Unmarshal(aar.RawClientData, &aar.ClientDataJSON); err != nil {
	// 	return nil, terrors.Errorf("error decoding client data: %v", err)
	// }

	return b, nil
}
