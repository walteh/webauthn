package credential

import (
	"context"
	"crypto/sha256"
	"errors"

	"github.com/rs/zerolog"

	"github.com/walteh/terrors"
	"github.com/walteh/webauthn/pkg/webauthn/authdata"
	"github.com/walteh/webauthn/pkg/webauthn/clientdata"
	"github.com/walteh/webauthn/pkg/webauthn/types"
	"github.com/walteh/webauthn/pkg/webauthn/webauthncbor"
)

// Parse the values returned in the authenticator response and perform attestation verification
// Step 8. This returns a fully decoded struct with the data put into a format that can be
// used to verify the user and credential that was created
func ParseAttestationInput(ctx context.Context, ccr types.AttestationInput) (*types.AttestationObject, error) {
	p := types.AttestationObject{}

	abc, err := clientdata.ParseClientData(ccr.UTF8ClientDataJSON)
	if err != nil {
		return nil, err
	}

	p.ClientData = abc

	err = webauthncbor.Unmarshal(ccr.AttestationObject, &p)
	if err != nil {
		return nil, terrors.Wrap(err, "error unmarshalling cbor attestation object")
	}

	// p.RawAuthData = hex.Hash(p.RawAuthData)

	// Step 8. Perform CBOR decoding on the attestationObject field of the AuthenticatorAttestationResponse
	// structure to obtain the attestation statement format fmt, the authenticator data authData, and
	// the attestation statement attStmt.
	dat, err := authdata.ParseAuthenticatorData(ctx, p.RawAuthData)
	if err != nil {
		return nil, terrors.Wrap(err, "error unmarshalling cbor auth data")
	}

	p.AuthData = dat

	p.Extensions = ccr.ClientExtensions

	if !p.AuthData.Flags.HasAttestedCredentialData() {
		return nil, terrors.New("Attestation missing attested credential data flag")
	}

	return &p, nil
}

// // Verifies the Client and Attestation data as laid out by §7.1. Registering a new credential
// // https://www.w3.org/TR/webauthn/#registering-a-new-credential
// Verify - Perform Steps 9 through 14 of registration verification, delegating Steps
func VerifyAttestationInput(ctx context.Context, args types.VerifyAttestationInputArgs) (*types.Credential, error) {

	attestationObject, err := ParseAttestationInput(ctx, args.Input)
	if err != nil {
		return nil, err
	}

	// Handles steps 3 through 6 - Verifying the Client Data against the Relying Party's stored data
	verifyError := clientdata.Verify(ctx, types.VerifyClientDataArgs{
		ClientData:         attestationObject.ClientData,
		StoredChallenge:    args.StoredChallenge,
		CeremonyType:       types.CreateCeremony,
		RelyingPartyOrigin: args.RelyingPartyOrigin,
	})

	if verifyError != nil {
		return nil, verifyError
	}

	// Step 7. Compute the hash of response.clientDataJSON using SHA-256.
	clientDataHash := sha256.Sum256([]byte(args.Input.UTF8ClientDataJSON))

	// Steps 9 through 12 are verified against the auth data.
	// These steps are identical to 11 through 14 for assertion
	// so we handle them with AuthData

	// Begin Step 9. Verify that the rpIdHash in authData is
	// the SHA-256 hash of the RP ID expected by the RP.
	// rpIDHash := sha256.Sum256([]byte(relyingPartyID))
	// Handle Steps 9 through 12

	authDataVerificationError := authdata.VerifyAuenticatorData(ctx, types.VerifyAuenticatorDataArgs{
		Data:                    attestationObject.RawAuthData,
		RelyingPartyID:          args.RelyingPartyID,
		AppId:                   "",
		RequireUserPresence:     false,
		RequireUserVerification: args.VerifyUser,
		LastSignCount:           0,
	})

	if authDataVerificationError != nil {
		return nil, authDataVerificationError
	}

	// Step 13. Determine the attestation statement format by performing a
	// USASCII case-sensitive match on fmt against the set of supported
	// WebAuthn Attestation Statement Format Identifier values. The up-to-date
	// list of registered WebAuthn Attestation Statement Format Identifier
	// values is maintained in the IANA registry of the same name
	// [WebAuthn-Registries] (https://www.w3.org/TR/webauthn/#biblio-webauthn-registries).

	// Since there is not an active registry yet, we'll check it against our internal
	// Supported types.

	// now := types.Now()

	tme := args.Provider.Time()

	abc := &types.Credential{
		PublicKey:       attestationObject.AuthData.AttData.CredentialPublicKey,
		RawID:           attestationObject.AuthData.AttData.CredentialID,
		Type:            "public-key",
		AttestationType: attestationObject.Format,
		AAGUID:          attestationObject.AuthData.AttData.AAGUID,
		SignCount:       attestationObject.AuthData.Counter,
		CloneWarning:    false,
		CreatedAt:       uint64(tme.Unix()),
		UpdatedAt:       uint64(tme.Unix()),
		SessionId:       nil,
		Receipt:         nil,
	}

	// But first let's make sure attestation is present. If it isn't, we don't need to handle
	// any of the following steps
	if attestationObject.Format == "none" {
		if len(attestationObject.AttStatement) != 0 {
			err := errors.New("attestation format none with attestation present")
			zerolog.Ctx(ctx).Error().Err(err).Send()
			return nil, err
		}
		return abc, nil
	}

	// formatHandler, valid := attestationRegistry[attestationObject.Format]
	// if !valid {
	// 	return nil, errors.New(fmt.Sprintf("Attestation format %s is unsupported", attestationObject.Format))
	// }

	// Step 14. Verify that attStmt is a correct attestation statement, conveying a valid attestation signature, by using
	// the attestation statement format fmt’s verification procedure given attStmt, authData and the hash of the serialized
	// client data computed in step 7.
	pk, attestationType, receipt, err := args.Provider.Attest(*attestationObject, clientDataHash[:])
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).
			Str("attestation_type", attestationType).
			Any("attestation_object", attestationObject).
			Msg("Error verifying attestation")
		return nil, err
	}

	if len(receipt) > 0 {
		rec, ok := receipt[0].([]byte)
		if !ok {
			err := errors.New("attestation receipt is not a byte array")
			zerolog.Ctx(ctx).Error().Err(err).Send()
			return nil, err
		}
		abc.Receipt = rec
	}

	abc.PublicKey = pk

	return abc, nil
}
