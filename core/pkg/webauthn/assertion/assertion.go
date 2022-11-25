package assertion

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/webauthn/authdata"
	"nugg-webauthn/core/pkg/webauthn/clientdata"
	"nugg-webauthn/core/pkg/webauthn/errors"
	"nugg-webauthn/core/pkg/webauthn/types"
	"nugg-webauthn/core/pkg/webauthn/webauthncose"
)

// Follow the remaining steps outlined in §7.2 Verifying an authentication assertion
// (https://www.w3.org/TR/webauthn/#verifying-assertion) and return an error if there
// is a failure during each step.
func VerifyAssertionInput(args types.VerifyAssertionInputArgs) error {
	// Steps 4 through 6 in verifying the assertion data (https://www.w3.org/TR/webauthn/#verifying-assertion) are
	// "assertive" steps, i.e "Let JSONtext be the result of running UTF-8 decode on the value of cData."
	// We handle these steps in part as we verify but also beforehand
	var (
		key interface{}
		err error
	)

	credId := types.NewDefaultCredentialIdentifier(args.Input.CredentialID)

	appID, err := credId.GetAppID(args.Extensions, args.AttestationType)
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
		Data:                    args.Input.RawAuthenticatorData,
		AppId:                   appID,
		RelyingPartyID:          args.RelyingPartyID,
		RequireUserVerification: args.VerifyUser,
		RequireUserPresence:     true,
		LastSignCount:           args.LastSignCount,
	})
	if validError != nil {
		return validError
	}

	// Step 15. Let hash be the result of computing a hash over the cData using SHA-256.
	clientDataHash := sha256.Sum256([]byte(args.Input.RawClientDataJSON))

	// Step 16. Using the credential public key looked up in step 3, verify that sig is
	// a valid signature over the binary concatenation of authData and hash.

	sigData := append(args.Input.RawAuthenticatorData, clientDataHash[:]...)

	if appID == "" {
		key, err = webauthncose.ParsePublicKey(args.CredentialPublicKey)
	} else {
		key, err = webauthncose.ParseFIDOPublicKey(args.CredentialPublicKey)
	}

	if err != nil {
		return errors.ErrAssertionSignature.WithMessage(fmt.Sprintf("Error parsing the assertion public key: %+v", err)).WithCaller()
	}

	valid, err := webauthncose.VerifySignature(key, sigData.Bytes(), args.Input.Signature.Bytes())
	if !valid || err != nil {
		log.Println("valid", valid, "err", err)
		return errors.ErrAssertionSignature.WithMessage("Error validating the assertion signature").WithCaller().
			WithRoot(err).
			WithKV("valid", valid).
			WithKV("key", key).
			WithKV("signature", args.Input.Signature.Hex()).
			WithKV("sigData", sigData).
			WithKV("appID", appID)
	}

	return nil
}

// {
// 	"challenge_id":"\(challenge.hexEncodedString())",
// 	"assertion_object":"\(assertion.hexEncodedString())",
// 	"session_id":"\(sessionID.hexEncodedString())",
// 	"provider":"apple"
//  "credential_id":"\(credentialID.hexEncodedString())"
// }

type AssertionResponse struct {
	UTF8ClientDataJSON string   `json:"client_data_json"`
	AssertionObject    hex.Hash `json:"assertion_object"`
	SessionID          hex.Hash `json:"session_id"`
	Provider           string   `json:"provider"`
	CredentialID       hex.Hash `json:"credential_id"`
}

func ParseAssertionResponse(response []byte) (*AssertionResponse, error) {
	var ar AssertionResponse
	err := json.Unmarshal(response, &ar)
	if err != nil {
		return nil, errors.ErrParsingData.WithMessage("Error parsing assertion response").WithRoot(err).WithCaller()
	}
	return &ar, nil
}
