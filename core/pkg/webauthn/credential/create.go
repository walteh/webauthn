package credential

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/nuggxyz/webauthn/pkg/errors"
	"github.com/nuggxyz/webauthn/pkg/hex"
	"github.com/nuggxyz/webauthn/pkg/webauthn/authdata"
	"github.com/nuggxyz/webauthn/pkg/webauthn/clientdata"
	"github.com/nuggxyz/webauthn/pkg/webauthn/types"
	"github.com/nuggxyz/webauthn/pkg/webauthn/webauthncbor"
)

// func VerifyAttestationInput(args types.VerifyAttestationInputArgs) (*types.Credential, error) {

// 	attestationObject, err := ParseAttestationInput(args.Input)
// 	if err != nil {
// 		return nil, errors.ErrParsingData.WithMessage("Error parsing attestation object").WithRoot(err).WithCaller()
// 	}

// 	// Handles steps 3 through 6 - Verifying the Client Data against the Relying Party's stored data
// 	verifyError := clientdata.Verify(attestationObject.ClientData, args.StoredChallenge, types.CreateCeremony, args.RelyingPartyOrigin)
// 	if verifyError != nil {
// 		return nil, verifyError
// 	}

// 	// Step 7. Compute the hash of response.clientDataJSON using SHA-256.
// 	clientDataHash := sha256.Sum256([]byte(args.Input.UTF8ClientDataJSON))

// 	// Step 8. Perform CBOR decoding on the attestationObject field of the AuthenticatorAttestationResponse
// 	// structure to obtain the attestation statement format fmt, the authenticator data authData, and the
// 	// attestation statement attStmt. is handled while

// 	// We do the above step while parsing and decoding the CredentialCreationResponse
// 	// Handle steps 9 through 14 - This verifies the attestaion object and
// 	pk, verifyError := verify(args.Provider, attestationObject, args.RelyingPartyID, clientDataHash[:], args.VerifyUser, true)
// 	if verifyError != nil {
// 		return nil, verifyError
// 	}

// 	// Step 15. If validation is successful, obtain a list of acceptable trust anchors (attestation root
// 	// certificates or ECDAA-Issuer public keys) for that attestation type and attestation statement
// 	// format fmt, from a trusted source or from policy. For example, the FIDO Metadata Service provides
// 	// one way to obtain such information, using the aaguid in the attestedCredentialData in authData.
// 	// [https://fidoalliance.org/specs/fido-v2.0-id-20180227/fido-metadata-service-v2.0-id-20180227.html]

// 	// TODO: There are no valid AAGUIDs yet or trust sources supported. We could implement policy for the RP in
// 	// the future, however.

// 	// Step 16. Assess the attestation trustworthiness using outputs of the verification procedure in step 14, as follows:
// 	// - If self attestation was used, check if self attestation is acceptable under Relying Party policy.
// 	// - If ECDAA was used, verify that the identifier of the ECDAA-Issuer public key used is included in
// 	//   the set of acceptable trust anchors obtained in step 15.
// 	// - Otherwise, use the X.509 certificates returned by the verification procedure to verify that the
// 	//   attestation public key correctly chains up to an acceptable root certificate.

// 	// TODO: We're not supporting trust anchors, self-attestation policy, or acceptable root certs yet

// 	// Step 17. Check that the credentialId is not yet registered to any other user. If registration is
// 	// requested for a credential that is already registered to a different user, the Relying Party SHOULD
// 	// fail this registration ceremony, or it MAY decide to accept the registration, e.g. while deleting
// 	// the older registration.

// 	// TODO: We can't support this in the code's current form, the Relying Party would need to check for this
// 	// against their database

// 	// Step 18 If the attestation statement attStmt verified successfully and is found to be trustworthy, then
// 	// register the new credential with the account that was denoted in the options.user passed to create(), by
// 	// associating it with the credentialId and credentialPublicKey in the attestedCredentialData in authData, as
// 	// appropriate for the Relying Party's system.

// 	// Step 19. If the attestation statement attStmt successfully verified but is not trustworthy per step 16 above,
// 	// the Relying Party SHOULD fail the registration ceremony.

// 	// TODO: Not implemented for the reasons mentioned under Step 16

// 	// z := &SavedCredential{
// 	// 	AAGUID: pcc.Response.AttestationObject.AuthData.AttData.AAGUID,
// 	// 	// AttestationType: pcc.Response.AttestationObject.AuthData.,
// 	// 	RawID:           pcc.Response.AttestationObject.AuthData.AttData.CredentialID,
// 	// 	SignCount:       pcc.Response.AttestationObject.AuthData.Counter,
// 	// 	PublicKey:       pk,
// 	// 	Type:            "public-key",
// 	// 	AttestationType: (pcc.Response.AttestationObject.Format),
// 	// 	Receipt:         r,
// 	// 	CloneWarning:    false,
// 	// 	CreatedAt:       uint64(time.Now().Unix()),
// 	// 	UpdatedAt:       uint64(time.Now().Unix()),
// 	// 	SessionId:       sessionId,
// 	// }

// 	return pk, nil
// }

// Parse the values returned in the authenticator response and perform attestation verification
// Step 8. This returns a fully decoded struct with the data put into a format that can be
// used to verify the user and credential that was created
func ParseAttestationInput(ccr types.AttestationInput) (*types.AttestationObject, error) {
	p := types.AttestationObject{}

	abc, err := clientdata.ParseClientData(ccr.UTF8ClientDataJSON)
	if err != nil {
		return nil, errors.ErrParsingData.WithMessage("Error parsing client data").WithRoot(err).WithCaller()
	}

	p.ClientData = abc

	err = webauthncbor.Unmarshal(ccr.AttestationObject, &p)
	if err != nil {
		log.Println("Error unmarshalling cbor attestation object", err)
		return nil, errors.ErrParsingData.WithInfo(err.Error()).WithCaller()
	}

	// p.RawAuthData = hex.Hash(p.RawAuthData)

	// Step 8. Perform CBOR decoding on the attestationObject field of the AuthenticatorAttestationResponse
	// structure to obtain the attestation statement format fmt, the authenticator data authData, and
	// the attestation statement attStmt.
	dat, err := authdata.ParseAuthenticatorData(p.RawAuthData)
	if err != nil {
		log.Println("Error unmarshalling cbor auth data", err)
		return nil, fmt.Errorf("error decoding auth data: %v", err)
	}

	p.AuthData = dat

	p.Extensions = ccr.ClientExtensions

	if !p.AuthData.Flags.HasAttestedCredentialData() {
		log.Println("Authenticator data does not contain attested credential data")
		return nil, errors.ErrAttestationFormat.WithInfo("Attestation missing attested credential data flag").WithCaller()
	}

	return &p, nil
}

// // Verifies the Client and Attestation data as laid out by §7.1. Registering a new credential
// // https://www.w3.org/TR/webauthn/#registering-a-new-credential
// Verify - Perform Steps 9 through 14 of registration verification, delegating Steps
func VerifyAttestationInput(args types.VerifyAttestationInputArgs) (*types.Credential, error) {

	attestationObject, err := ParseAttestationInput(args.Input)
	if err != nil {
		return nil, errors.ErrParsingData.WithMessage("Error parsing attestation object").WithRoot(err).WithCaller()
	}

	// Handles steps 3 through 6 - Verifying the Client Data against the Relying Party's stored data
	verifyError := clientdata.Verify(types.VerifyClientDataArgs{
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

	authDataVerificationError := authdata.VerifyAuenticatorData(types.VerifyAuenticatorDataArgs{
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

	now := types.Now()

	abc := &types.Credential{
		PublicKey:       attestationObject.AuthData.AttData.CredentialPublicKey,
		RawID:           attestationObject.AuthData.AttData.CredentialID,
		Type:            "public-key",
		AttestationType: attestationObject.Format,
		AAGUID:          attestationObject.AuthData.AttData.AAGUID,
		SignCount:       attestationObject.AuthData.Counter,
		CloneWarning:    false,
		CreatedAt:       now,
		UpdatedAt:       now,
		SessionId:       hex.Hash{},
		Receipt:         hex.Hash{},
	}

	// But first let's make sure attestation is present. If it isn't, we don't need to handle
	// any of the following steps
	if attestationObject.Format == "none" {
		if len(attestationObject.AttStatement) != 0 {
			return nil, errors.ErrAttestationFormat.WithInfo("Attestation format none with attestation present").WithCaller()
		}
		return abc, nil
	}

	// formatHandler, valid := attestationRegistry[attestationObject.Format]
	// if !valid {
	// 	return nil, errors.ErrAttestationFormat.WithInfo(fmt.Sprintf("Attestation format %s is unsupported", attestationObject.Format)).WithCaller()
	// }

	// Step 14. Verify that attStmt is a correct attestation statement, conveying a valid attestation signature, by using
	// the attestation statement format fmt’s verification procedure given attStmt, authData and the hash of the serialized
	// client data computed in step 7.
	pk, attestationType, receipt, err := args.Provider.Attest(*attestationObject, clientDataHash[:])
	if err != nil {
		return nil, err.(*errors.Error).WithInfo(attestationType).WithCaller()
	}

	if len(receipt) > 0 {
		rec, ok := receipt[0].([]byte)
		if !ok {
			return nil, errors.ErrAttestationFormat.WithInfo("Attestation receipt is not a byte array").WithCaller()
		}
		abc.Receipt = rec
	}

	abc.PublicKey = pk

	return abc, nil
}
