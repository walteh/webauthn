package providers

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/nuggxyz/webauthn/pkg/errors"
	"github.com/nuggxyz/webauthn/pkg/hex"
	"github.com/nuggxyz/webauthn/pkg/webauthn/types"
	"github.com/nuggxyz/webauthn/pkg/webauthn/webauthncose"
)

type AppleAttestationProvider struct{}

func NewAppleAttestationProvider() *AppleAttestationProvider {
	return &AppleAttestationProvider{}
}

func (me *AppleAttestationProvider) ID() string {
	return "apple"
}

// From §8.8. https://www.w3.org/TR/webauthn-2/#sctn-apple-anonymous-attestation
// The apple attestation statement looks like:
// $$attStmtType //= (
//
//	fmt: "apple",
//	attStmt: appleStmtFormat
//
// )
//
//	appleStmtFormat = {
//			x5c: [ credCert: bytes, * (caCert: bytes) ]
//	  }
func (me *AppleAttestationProvider) Handler(att types.AttestationObject, clientDataHash []byte) (hex.Hash, string, []interface{}, error) {

	// VerifyAttestation verifies the attestation signature on the authenticator data

	// 7. Verify that the authenticator data’s counter field equals 0.
	if att.AuthData.Counter != 0 {
		return nil, "", nil, errors.ErrVerification.WithMessage(fmt.Sprintf("Counter was not 0, but %d\n", att.AuthData.Counter))
	}

	// 8. Verify that the authenticator data’s aaguid field is either appattestdevelop if operating in the development environment,
	// or appattest followed by seven 0x00 bytes if operating in the production environment.
	aaguid := make([]byte, 16)
	copy(aaguid, []byte("appattestdevelop"))
	if !bytes.Equal(att.AuthData.AttData.AAGUID, aaguid) {
		return nil, "", nil, errors.ErrVerification.WithMessage("AAGUID was not appattestdevelop\n")
	}

	// Step 1. Verify that attStmt is valid CBOR conforming to the syntax defined
	// above and perform CBOR decoding on it to extract the contained fields.

	// If x5c is not present, return an error
	x5c, x509present := att.AttStatement["x5c"].([]interface{})
	if !x509present {
		// Handle Basic Attestation steps for the x509 Certificate
		return nil, "", nil, errors.ErrAttestationFormat.WithMessage("Error retreiving x5c value")
	}

	credCertBytes, valid := x5c[0].([]byte)
	if !valid {
		return nil, "", nil, errors.ErrAttestation.WithMessage("Error getting certificate from x5c cert chain")
	}

	credCert, err := x509.ParseCertificate(credCertBytes)
	if err != nil {
		return nil, "", nil, errors.ErrAttestationFormat.WithMessage(fmt.Sprintf("Error parsing certificate from ASN.1 data: %+v", err))
	}

	// Step 2. Concatenate authenticatorData and clientDataHash to form nonceToHash.
	nonceToHash := append(att.RawAuthData, clientDataHash...)

	// Step 3. Perform SHA-256 hash of nonceToHash to produce nonce.
	nonce := sha256.Sum256(nonceToHash)

	// Step 4. Verify that nonce equals the value of the extension with OID 1.2.840.113635.100.8.2 in credCert.
	var attExtBytes []byte
	for _, ext := range credCert.Extensions {
		if ext.Id.Equal([]int{1, 2, 840, 113635, 100, 8, 2}) {
			attExtBytes = ext.Value
		}
	}
	if len(attExtBytes) == 0 {
		return nil, "", nil, errors.ErrAttestationFormat.WithMessage("Attestation certificate extensions missing 1.2.840.113635.100.8.2")
	}

	decoded := AppleAnonymousAttestation{}
	_, err = asn1.Unmarshal([]byte(attExtBytes), &decoded)
	if err != nil {
		return nil, "", nil, errors.ErrAttestationFormat.WithMessage("Unable to parse apple attestation certificate extensions")
	}

	if !bytes.Equal(decoded.Nonce, nonce[:]) || err != nil {
		return nil, "", nil, errors.ErrInvalidAttestation.WithMessage("Attestation certificate does not contain expected nonce")
	}

	// Step 5. Verify that the credential public key equals the Subject Public Key of credCert.
	// TODO: Probably move this part to webauthncose.go
	pubKey, err := webauthncose.ParsePublicKey(att.AuthData.AttData.CredentialPublicKey)
	if err != nil {
		return nil, "", nil, errors.ErrInvalidAttestation.WithMessage(fmt.Sprintf("Error parsing public key: %+v\n", err))
	}
	credPK := pubKey.(webauthncose.EC2PublicKeyData)
	subjectPK := credCert.PublicKey.(*ecdsa.PublicKey)
	credPKInfo := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     big.NewInt(0).SetBytes(credPK.XCoord),
		Y:     big.NewInt(0).SetBytes(credPK.YCoord),
	}
	if !credPKInfo.Equal(subjectPK) {
		return nil, "", nil, errors.ErrInvalidAttestation.WithMessage("Certificate public key does not match public key in authData")
	}

	// Step 6. If successful, return implementation-specific values representing attestation type Anonymization CA and attestation trust path x5c.
	return att.AuthData.AttData.CredentialPublicKey, "", x5c, nil
}

// Apple has not yet publish schema for the extension(as of JULY 2021.)
type AppleAnonymousAttestation struct {
	Nonce []byte `asn1:"tag:1,explicit"`
}
