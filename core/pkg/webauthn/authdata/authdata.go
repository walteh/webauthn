package authdata

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/webauthn/errors"
	"nugg-webauthn/core/pkg/webauthn/types"
	"nugg-webauthn/core/pkg/webauthn/webauthncbor"

	"github.com/k0kubun/pp"
)

const (
	minAuthDataLength     = 37
	minAttestedAuthLength = 55

	// https://w3c.github.io/webauthn/#attested-credential-data
	maxCredentialIDLength = 1023
)

// Unmarshal will take the raw Authenticator Data and marshalls it into AuthenticatorData for further validation.
// The authenticator data has a compact but extensible encoding. This is desired since authenticators can be
// devices with limited capabilities and low power requirements, with much simpler software stacks than the client platform.
// The authenticator data structure is a byte array of 37 bytes or more, and is laid out in this table:
// https://www.w3.org/TR/webauthn/#table-authData

func ParseAuthenticatorData(rawAuthData hex.Hash) (a types.AuthenticatorData, err error) {
	return ParseAuthenticatorDataSavedAttestedCredential(rawAuthData, false)
}
func ParseAuthenticatorDataSavedAttestedCredential(rawAuthData hex.Hash, savedAttestedCredentials bool) (a types.AuthenticatorData, err error) {
	a = types.AuthenticatorData{}
	byt := rawAuthData.Bytes()
	if minAuthDataLength > len(rawAuthData) {
		err := errors.ErrBadRequest.WithMessage("Authenticator data length too short")
		info := fmt.Sprintf("Expected data greater than %d bytes. Got %d bytes\n", minAuthDataLength, len(rawAuthData))
		return a, err.WithInfo(info)
	}

	a.RPIDHash = byt[:32]
	a.Flags = types.AuthenticatorFlags(byt[32])
	a.Counter = uint64(binary.BigEndian.Uint32(byt[33:37]))

	remaining := len(byt) - minAuthDataLength

	if a.Flags.HasAttestedCredentialData() {
		if len(byt) > minAttestedAuthLength {
			att, validError := ParseAttestedAuthData(byt)
			if validError != nil {
				return a, validError
			}
			a.AttData = att
			attDataLen := len(a.AttData.AAGUID) + 2 + len(a.AttData.CredentialID) + len(a.AttData.CredentialPublicKey)
			remaining = remaining - attDataLen
		} else if !savedAttestedCredentials {
			return a, errors.ErrBadRequest.WithMessage("Attested credential flag set but data is missing")
		}
	} else {
		if !a.Flags.HasExtensions() && len(byt) != 37 {
			return a, errors.ErrBadRequest.WithMessage("Attested credential flag not set")
		}
	}

	if a.Flags.HasExtensions() {
		if remaining != 0 {
			a.ExtData = byt[len(byt)-remaining:]
			remaining -= len(a.ExtData)
		} else {
			return a, errors.ErrBadRequest.WithMessage("Extensions flag set but extensions data is missing")
		}
	}

	if remaining != 0 {
		return a, errors.ErrBadRequest.WithMessage("Leftover bytes decoding AuthenticatorData")
	}

	return a, nil
}

// If Attestation Data is present, unmarshall that into the appropriate public key structure
func ParseAttestedAuthData(rawAuthData []byte) (types.AttestedCredentialData, error) {
	a := types.AttestedCredentialData{}
	a.AAGUID = rawAuthData[37:53]
	idLength := binary.BigEndian.Uint16(rawAuthData[53:55])
	if len(rawAuthData) < int(55+idLength) {
		return a, errors.ErrBadRequest.WithMessage("Authenticator attestation data length too short")
	}
	if idLength > maxCredentialIDLength {
		return a, errors.ErrBadRequest.WithMessage("Authenticator attestation data credential id length too long")
	}
	a.CredentialID = rawAuthData[55 : 55+idLength]
	a.CredentialPublicKey = unmarshalCredentialPublicKey(rawAuthData[55+idLength:])
	return a, nil
}

// Unmarshall the credential's Public Key into CBOR encoding
func unmarshalCredentialPublicKey(keyBytes []byte) []byte {
	var m interface{}
	webauthncbor.Unmarshal(keyBytes, &m)
	rawBytes, _ := webauthncbor.Marshal(m)
	if bytes.Equal(rawBytes, keyBytes) {
		log.Println("Credential Public Key is not CBOR encoded")
	}
	return rawBytes
}

// ResidentKeyRequired - Require that the key be private key resident to the client device
func ResidentKeyRequired() *bool {
	required := true
	return &required
}

// ResidentKeyUnrequired - Do not require that the private key be resident to the client device.
func ResidentKeyUnrequired() *bool {
	required := false
	return &required
}

// Verify on AuthenticatorData handles Steps 9 through 12 for Registration
// and Steps 11 through 14 for Assertion.
func VerifyAuenticatorData(args types.VerifyAuenticatorDataArgs) error {

	// Begin Step 11. Verify that the rpIdHash in authData is the SHA-256 hash of the RP ID expected by the RP.
	rpIDHash := sha256.Sum256([]byte(args.RelyingPartyID))

	var appIDHash [32]byte
	if args.AppId != "" {
		appIDHash = sha256.Sum256([]byte(args.AppId))
	}

	hasAttestationCredsSaved := !args.OptionalAttestedCredentialData.CredentialPublicKey.IsZero()

	data, err := ParseAuthenticatorDataSavedAttestedCredential(args.Data, hasAttestationCredsSaved)
	if err != nil {
		return err
	}

	if hasAttestationCredsSaved && data.AttData.CredentialPublicKey.IsZero() {
		data.AttData = args.OptionalAttestedCredentialData
	}

	pp.Println(data)

	// Registration Step 9 & Assertion Step 11
	// Verify that the RP ID hash in authData is indeed the SHA-256
	// hash of the RP ID expected by the RP.
	if !bytes.Equal(data.RPIDHash, rpIDHash[:]) && !bytes.Equal(data.RPIDHash[:], appIDHash[:]) {
		return errors.ErrVerification.
			WithInfo(fmt.Sprintf("RP Hash mismatch. Expected %x and Received %x\n", data.RPIDHash, rpIDHash)).
			WithKV("AuthenticatorData", data)
	}

	// Registration Step 10 & Assertion Step 12
	// Verify that the User Present bit of the flags in authData is set.
	if args.RequireUserPresence && !data.Flags.UserPresent() {
		return errors.ErrVerification.WithInfo(fmt.Sprintln("User presence flag not set by authenticator"))
	}

	// Registration Step 11 & Assertion Step 13
	// If user verification is required for this assertion, verify that
	// the User Verified bit of the flags in authData is set.
	if args.RequireUserVerification && !data.Flags.UserVerified() {
		return errors.ErrVerification.WithInfo(fmt.Sprintln("User verification required but flag not set by authenticator"))
	}

	if data.Counter < args.LastSignCount {
		return errors.ErrVerification.WithInfo(fmt.Sprintf("Counter was not not greater than previous  %d\n", data.Counter))
	}

	// Registration Step 12 & Assertion Step 14
	// Verify that the values of the client extension outputs in clientExtensionResults
	// and the authenticator extension outputs in the extensions in authData are as
	// expected, considering the client extension input values that were given as the
	// extensions option in the create() call. In particular, any extension identifier
	// values in the clientExtensionResults and the extensions in authData MUST be also be
	// present as extension identifier values in the extensions member of options, i.e., no
	// extensions are present that were not requested. In the general case, the meaning
	// of "are as expected" is specific to the Relying Party and which extensions are in use.

	// This is not yet fully implemented by the spec or by browsers

	return nil
}
