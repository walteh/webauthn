package authdata

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"log"

	"github.com/rs/zerolog"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/types"
	"github.com/walteh/webauthn/pkg/webauthn/webauthncbor"
)

const (
	MinAuthDataLength     = 37
	MinAttestedAuthLength = 55

	// https://w3c.github.io/webauthn/#attested-credential-data
	MaxCredentialIDLength = 1023
)

// Unmarshal will take the raw Authenticator Data and marshalls it into AuthenticatorData for further validation.
// The authenticator data has a compact but extensible encoding. This is desired since authenticators can be
// devices with limited capabilities and low power requirements, with much simpler software stacks than the client platform.
// The authenticator data structure is a byte array of 37 bytes or more, and is laid out in this table:
// https://www.w3.org/TR/webauthn/#table-authData

func ParseAuthenticatorData(ctx context.Context, rawAuthData hex.Hash) (a types.AuthenticatorData, err error) {
	return ParseAuthenticatorDataSavedAttestedCredential(ctx, rawAuthData, false)
}
func ParseAuthenticatorDataSavedAttestedCredential(ctx context.Context, rawAuthData hex.Hash, savedAttestedCredentials bool) (a types.AuthenticatorData, err error) {
	a = types.AuthenticatorData{}
	byt := rawAuthData.Bytes()
	if MinAuthDataLength > len(rawAuthData) {
		err := errors.New("authenticator data length too short")
		zerolog.Ctx(ctx).Error().Err(err).Msg("authenticator data length too short")
		return a, err
	}

	a.RPIDHash = byt[:32]
	a.Flags = types.AuthenticatorFlags(byt[32])
	a.Counter = uint64(binary.BigEndian.Uint32(byt[33:37]))

	remaining := len(byt) - MinAuthDataLength

	if a.Flags.HasAttestedCredentialData() {
		if len(byt) > MinAttestedAuthLength {
			att, validError := ParseAttestedAuthData(ctx, byt)
			if validError != nil {
				return a, validError
			}
			a.AttData = att
			attDataLen := len(a.AttData.AAGUID) + 2 + len(a.AttData.CredentialID) + len(a.AttData.CredentialPublicKey)
			remaining = remaining - attDataLen
		} else if !savedAttestedCredentials {
			err := errors.New("attested credential flag set but data is missing")
			zerolog.Ctx(ctx).Error().Err(err).Send()
			return a, err
		}
	} else {
		if !a.Flags.HasExtensions() && len(byt) != 37 {
			err := errors.New("attested credential flag not set")
			zerolog.Ctx(ctx).Error().Err(err).Send()
			return a, err
		}
	}

	if a.Flags.HasExtensions() {
		if remaining != 0 {
			a.ExtData = byt[len(byt)-remaining:]
			remaining -= len(a.ExtData)
		} else {
			err := errors.New("extensions flag set but extensions data is missing")
			zerolog.Ctx(ctx).Error().Err(err).Send()
			return a, err
		}
	}

	if remaining != 0 {
		err := errors.New("leftover bytes decoding authenticatordata")
		zerolog.Ctx(ctx).Error().Err(err).Send()
		return a, err
	}

	return a, nil
}

// If Attestation Data is present, unmarshall that into the appropriate public key structure
func ParseAttestedAuthData(ctx context.Context, rawAuthData []byte) (types.AttestedCredentialData, error) {
	a := types.AttestedCredentialData{}
	a.AAGUID = rawAuthData[37:53]
	idLength := binary.BigEndian.Uint16(rawAuthData[53:55])
	if len(rawAuthData) < int(55+idLength) {
		err := errors.New("authenticator attestation data length too short")
		zerolog.Ctx(ctx).Error().Err(err).Send()
		return a, err
	}
	if idLength > MaxCredentialIDLength {
		err := errors.New("authenticator attestation data credential id length too long")
		zerolog.Ctx(ctx).Error().Err(err).Send()
		return a, err
	}
	a.CredentialID = rawAuthData[55 : 55+idLength]
	a.CredentialPublicKey = UnmarshalCredentialPublicKey(rawAuthData[55+idLength:])
	return a, nil
}

// Unmarshall the credential's Public Key into CBOR encoding
func UnmarshalCredentialPublicKey(keyBytes []byte) []byte {
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
func VerifyAuenticatorData(ctx context.Context, args types.VerifyAuenticatorDataArgs) error {

	// Begin Step 11. Verify that the rpIdHash in authData is the SHA-256 hash of the RP ID expected by the RP.
	rpIDHash := sha256.Sum256([]byte(args.RelyingPartyID))

	var appIDHash [32]byte
	if args.AppId != "" {
		appIDHash = sha256.Sum256([]byte(args.AppId))
	}

	hasAttestationCredsSaved := !args.OptionalAttestedCredentialData.CredentialPublicKey.IsZero()

	data, err := ParseAuthenticatorDataSavedAttestedCredential(ctx, args.Data, hasAttestationCredsSaved)
	if err != nil {
		return err
	}

	if hasAttestationCredsSaved && data.AttData.CredentialPublicKey.IsZero() {
		data.AttData = args.OptionalAttestedCredentialData
	}

	// Registration Step 9 & Assertion Step 11
	// Verify that the RP ID hash in authData is indeed the SHA-256
	// hash of the RP ID expected by the RP.
	if !bytes.Equal(data.RPIDHash, rpIDHash[:]) && !bytes.Equal(data.RPIDHash[:], appIDHash[:]) {
		err := errors.New("rp hash mismatch")

		zerolog.Ctx(ctx).Error().Err(err).
			Str("data.RPIDHash", data.RPIDHash.Hex()).
			Str("rpIDHash[:]", hex.Bytes2Hex(rpIDHash[:])).
			Str("appIDHash[:]", hex.Bytes2Hex(appIDHash[:])).
			Msg("RP Hash mismatch")

		return err
	}

	// Registration Step 10 & Assertion Step 12
	// Verify that the User Present bit of the flags in authData is set.
	if args.RequireUserPresence && !data.Flags.UserPresent() {
		err := errors.New("user presence flag not set by authenticator")
		zerolog.Ctx(ctx).Error().Err(err).Send()
		return err
	}

	// Registration Step 11 & Assertion Step 13
	// If user verification is required for this assertion, verify that
	// the User Verified bit of the flags in authData is set.
	if args.RequireUserVerification && !data.Flags.UserVerified() {
		err := errors.New("user verification required but flag not set by authenticator")
		zerolog.Ctx(ctx).Error().Err(err).Send()
		return err

	}

	if data.Counter < args.LastSignCount {
		err := errors.New("counter value too low")
		zerolog.Ctx(ctx).Error().Err(err).Uint64("data.Counter", data.Counter).Uint64("args.LastSignCount", args.LastSignCount).Msg("Counter value too low")
		return err
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
