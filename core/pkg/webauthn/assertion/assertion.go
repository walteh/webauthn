package assertion

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nuggxyz/webauthn/pkg/errors"
	"github.com/nuggxyz/webauthn/pkg/hex"
	"github.com/nuggxyz/webauthn/pkg/webauthn/authdata"
	"github.com/nuggxyz/webauthn/pkg/webauthn/clientdata"
	"github.com/nuggxyz/webauthn/pkg/webauthn/types"
	"github.com/nuggxyz/webauthn/pkg/webauthn/webauthncose"

	"github.com/ugorji/go/codec"
)

// Follow the remaining steps outlined in §7.2 Verifying an authentication assertion
// (https://www.w3.org/TR/webauthn/#verifying-assertion) and return an error if there
// is a failure during each step.
func VerifyAssertionInput(args types.VerifyAssertionInputArgs) error {
	// Steps 4 through 6 in verifying the assertion data (https://www.w3.org/TR/webauthn/#verifying-assertion) are
	// "assertive" steps, i.e "Let JSONtext be the result of running UTF-8 decode on the value of cData."
	// We handle these steps in part as we verify but also beforehand
	var (
		key      interface{}
		err      error
		asserter types.AssertionObject
	)

	if args.Input.AssertionObject == nil && !args.Input.RawAssertionObject.IsZero() {
		asserter, err = ParseAssertionObject(args.Input.RawAssertionObject)
		if err != nil {
			return err
		}
	} else {
		if args.Input.AssertionObject == nil {
			return errors.ErrBadRequest.WithCaller().WithMessage("Assertion object is missing")
		}
		asserter = *args.Input.AssertionObject
	}

	credId := types.NewDefaultCredentialIdentifier(args.Input.CredentialID)

	appID, err := credId.GetAppID(args.Extensions, args.CredentialAttestationType)
	if err != nil {
		return errors.ErrParsingData.WithMessage("Error getting appID").WithRoot(err).WithCaller()
	}

	clientd, err := clientdata.ParseClientData(args.Input.RawClientDataJSON)
	if err != nil {
		return errors.ErrParsingData.WithMessage("Error parsing client data").WithRoot(err).WithCaller()
	}

	// Handle steps 7 through 10 of assertion by verifying stored data against the Collected Client Data
	// returned by the authenticator
	validError := clientdata.Verify(types.VerifyClientDataArgs{
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
	validError = authdata.VerifyAuenticatorData(types.VerifyAuenticatorDataArgs{
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
			return errors.ErrParsingData.WithCaller().WithMessage("Failed to parse the public key")
		}

		pubkey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
		nonceHash := sha256.Sum256(nonce[:])

		valid := ecdsa.VerifyASN1(pubkey, nonceHash[:], asserter.Signature)
		if !valid {
			return errors.ErrAssertionSignature.WithCaller().WithMessage("Error validating the assertion signature.\n")
		}
	} else {
		if appID == "" {
			key, err = webauthncose.ParsePublicKey(args.CredentialPublicKey)
		} else {
			key, err = webauthncose.ParseFIDOPublicKey(args.CredentialPublicKey)
		}

		if err != nil {
			return errors.ErrAssertionSignature.WithRoot(err).WithMessage(fmt.Sprintf("Error parsing the assertion public key: %+v", err)).WithCaller()
		}

		valid, err := webauthncose.VerifySignature(key, sigData[:], asserter.Signature)
		if !valid || err != nil {
			log.Println("valid", valid, "err", err)
			return errors.ErrAssertionSignature.WithMessage("Error validating the assertion signature").WithCaller().
				WithRoot(err).
				WithKV("valid", valid).
				WithKV("key", key).
				WithKV("signature", asserter.Signature.Hex()).
				WithKV("sigData", sigData.Hex()).
				WithKV("appID", appID)
		}
	}

	return nil
}

type AssertionResponse struct {
	UTF8ClientDataJSON string   `json:"client_data_json"`
	AssertionObject    hex.Hash `json:"assertion_object"`
	SessionID          hex.Hash `json:"session_id"`
	Provider           string   `json:"provider"`
	CredentialID       hex.Hash `json:"credential_id"`
}

func ParseNotFidoAssertionInput(response []byte) (types.AssertionInput, error) {
	return ParseAssertionInput(response, types.NotFidoAttestationType)
}

func ParseFidoAssertionInput(response []byte) (types.AssertionInput, error) {
	return ParseAssertionInput(response, types.FidoAttestationType)
}

func ParseAssertionInput(response []byte, attestationType types.CredentialAttestationType) (types.AssertionInput, error) {
	var parsed AssertionResponse
	err := json.Unmarshal(response, &parsed)
	if err != nil {
		return types.AssertionInput{}, errors.ErrParsingData.WithMessage("Error parsing assertion response").WithRoot(err).WithCaller()
	}

	abc := types.AssertionInput{
		UserID:            parsed.SessionID,
		CredentialID:      parsed.CredentialID,
		RawClientDataJSON: parsed.UTF8ClientDataJSON,
		// RawAuthenticatorData: nil,
		RawAssertionObject: parsed.AssertionObject,
		AssertionObject:    nil,
		// Type:            attestationType,
	}

	return abc, nil
}

func ParseAssertionObject(input hex.Hash) (types.AssertionObject, error) {
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
		return b, errors.Err0x66CborDecode.WithCaller().WithInfo(err.Error())
	}

	b.RawAuthenticatorData = a.RawAuthenticatorData
	b.Signature = a.Signature

	// ad, err := authdata.ParseAuthenticatorData(a.RawAuthenticatorData)
	// if err != nil {
	// 	return a, fmt.Errorf("error decoding auth data: %v", err)
	// }

	// a.AuthenticatorData = ad

	// if err := json.Unmarshal(aar.RawClientData, &aar.ClientDataJSON); err != nil {
	// 	return nil, fmt.Errorf("error decoding client data: %v", err)
	// }

	return b, nil
}
