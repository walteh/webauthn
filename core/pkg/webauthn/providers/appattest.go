package providers

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/webauthn/errors"
	"nugg-webauthn/core/pkg/webauthn/types"
)

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

type AppAttest struct {
	production bool
}

func NewAppAttestSandbox() *AppAttest {
	return &AppAttest{
		production: false,
	}
}

func NewAppAttestProduction() *AppAttest {
	return &AppAttest{
		production: true,
	}
}

func (me *AppAttest) ID() string {
	return "apple-appattest"
}

const APPLE_ROOT_CERT = `-----BEGIN CERTIFICATE-----
MIICITCCAaegAwIBAgIQC/O+DvHN0uD7jG5yH2IXmDAKBggqhkjOPQQDAzBSMSYw
JAYDVQQDDB1BcHBsZSBBcHAgQXR0ZXN0YXRpb24gUm9vdCBDQTETMBEGA1UECgwK
QXBwbGUgSW5jLjETMBEGA1UECAwKQ2FsaWZvcm5pYTAeFw0yMDAzMTgxODMyNTNa
Fw00NTAzMTUwMDAwMDBaMFIxJjAkBgNVBAMMHUFwcGxlIEFwcCBBdHRlc3RhdGlv
biBSb290IENBMRMwEQYDVQQKDApBcHBsZSBJbmMuMRMwEQYDVQQIDApDYWxpZm9y
bmlhMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAERTHhmLW07ATaFQIEVwTtT4dyctdh
NbJhFs/Ii2FdCgAHGbpphY3+d8qjuDngIN3WVhQUBHAoMeQ/cLiP1sOUtgjqK9au
Yen1mMEvRq9Sk3Jm5X8U62H+xTD3FE9TgS41o0IwQDAPBgNVHRMBAf8EBTADAQH/
MB0GA1UdDgQWBBSskRBTM72+aEH/pwyp5frq5eWKoTAOBgNVHQ8BAf8EBAMCAQYw
CgYIKoZIzj0EAwMDaAAwZQIwQgFGnByvsiVbpTKwSga0kP0e8EeDS4+sQmTvb7vn
53O5+FRXgeLhpJ06ysC5PrOyAjEAp5U4xDgEgllF7En3VcE3iexZZtKeYnpqtijV
oyFraWVIyd/dganmrduC1bmTBGwD
-----END CERTIFICATE-----`

func (me *AppAttest) Attest(att types.AttestationObject, clientDataHash []byte) (hex.Hash, string, []interface{}, error) {

	// 7. Verify that the authenticator data’s counter field equals 0.
	if att.AuthData.Counter != 0 {
		return nil, "", nil, errors.ErrVerification.WithMessage(fmt.Sprintf("Counter was not 0, but %d\n", att.AuthData.Counter))
	}

	// 8. Verify that the authenticator data’s aaguid field is either appattestdevelop if operating in the development environment,
	// or appattest followed by seven 0x00 bytes if operating in the production environment.
	aaguid := make([]byte, 16)
	if me.production {
		copy(aaguid, []byte("appattest"))
	} else {
		copy(aaguid, []byte("appattestdevelop"))
	}
	if !bytes.Equal(att.AuthData.AttData.AAGUID, aaguid) {
		return nil, "", nil, errors.ErrVerification.WithMessage("AAGUID was not appattestdevelop\n")
	}

	roots := x509.NewCertPool()
	intermediates := x509.NewCertPool()

	// Add Apple root Cert
	ok := roots.AppendCertsFromPEM([]byte(APPLE_ROOT_CERT))
	if !ok {
		return nil, "", nil, errors.ErrAttestationFormat.WithMessage("Error adding root certificate to pool.")
	}

	x5c, x509present := att.AttStatement["x5c"].([]interface{})
	if !x509present {
		return nil, "", nil, errors.ErrAttestationFormat.WithMessage("Error retrieving x5c value")
	}

	_, receiptPresent := att.AttStatement["receipt"].([]byte)
	if !receiptPresent {
		return nil, "", nil, errors.ErrAttestationFormat.WithMessage("Error retreiving receipt value")
	}

	for _, c := range x5c {
		cb, cv := c.([]byte)
		if !cv {
			return nil, "", nil, errors.ErrAttestationCertificate.WithMessage("Error getting certificate from x5c cert chain 1")
		}
		ct, err := x509.ParseCertificate(cb)
		if err != nil {
			return nil, "", nil, errors.ErrAttestationCertificate.WithMessage(fmt.Sprintf("Error parsing certificate from ASN.1 data: %+v", err))
		}
		if ct.IsCA {
			intermediates.AddCert(ct)
		}
	}

	credCertBytes, valid := x5c[0].([]byte)
	if !valid {
		return nil, "", nil, errors.ErrAttestationCertificate.WithMessage("Error getting certificate from x5c cert chain 2")
	}

	credCert, err := x509.ParseCertificate(credCertBytes)
	if err != nil {
		return nil, "", nil, errors.ErrAttestationCertificate.WithMessage(fmt.Sprintf("Error parsing certificate from ASN.1 data: %+v", err))
	}

	// Create verification options.
	verifyOptions := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
	}

	// 1. Verify that the x5c array contains the intermediate and leaf certificates for App Attest,
	// starting from the credential certificate stored in the first data buffer in the array (credcert).
	// Verify the validity of the certificates using Apple’s root certificate.
	_, err = credCert.Verify(verifyOptions)
	if err != nil {
		return nil, "", nil, errors.ErrAttestationCertificate.WithMessage(fmt.Sprintf("Invalid certificate %+v", err))
	}

	// 2. Create clientDataHash as the SHA256 hash of the one-time challenge sent to your app before performing the attestation,
	// and append that hash to the end of the authenticator data (authData from the decoded object).
	nonceData := append(att.RawAuthData, clientDataHash...)

	// 3. Generate a new SHA256 hash of the composite item to create nonce.
	nonce := sha256.Sum256(nonceData)

	// 4. Obtain the value of the credCert extension with OID 1.2.840.113635.100.8.2, which is a DER-encoded ASN.1 sequence.
	// Decode the sequence and extract the single octet string that it contains.
	// Verify that the string equals nonce.
	credCertOID := asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 8, 2}
	var credCertId []byte
	for _, extension := range credCert.Extensions {
		if extension.Id.Equal(credCertOID) {
			credCertId = extension.Value
		}
	}

	if len(credCertId) <= 0 {
		return nil, "", nil, errors.ErrInvalidAttestation.WithMessage("Certificate did not contain credCert extension")
	}
	var unMarshalledCredCertOctet []asn1.RawValue
	var unMarshalledCredCert asn1.RawValue
	asn1.Unmarshal(credCertId, &unMarshalledCredCertOctet)
	asn1.Unmarshal(unMarshalledCredCertOctet[0].Bytes, &unMarshalledCredCert)
	if !bytes.Equal(nonce[:], unMarshalledCredCert.Bytes) {
		return nil, "", nil, errors.ErrInvalidAttestation.WithMessage("Certificate CredCert extension does not match nonce.").WithKV("nonce", nonce[:]).WithKV("credCert", unMarshalledCredCert.Bytes)
	}

	// 5. Create the SHA256 hash of the public key in credCert, and verify that it matches the key identifier from your app.
	var publicKeyBytes []byte
	switch pub := credCert.PublicKey.(type) {
	case *ecdsa.PublicKey:
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
		pubKeyHash := sha256.Sum256(publicKeyBytes)
		if !bytes.Equal(pubKeyHash[:], att.AuthData.AttData.CredentialID) {
			return nil, "", nil, errors.ErrInvalidAttestation.WithMessage("The key id is not a valid SHA256 hash of the certificate public key.")
		}
	default:
		return nil, "", nil, errors.ErrInvalidAttestation.WithMessage("Wrong algorithm")
	}

	// Return x963-encoded public key and receipt.
	return hex.BytesToHash(publicKeyBytes), string(aaguid), []interface{}{att.AttStatement["receipt"]}, nil
}

// // Apple has not yet publish schema for the extension(as of JULY 2021.)
// type AppleAnonymousAttestation struct {
// 	Nonce []byte `asn1:"tag:1,explicit"`
// }
