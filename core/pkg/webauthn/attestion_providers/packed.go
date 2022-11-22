package attestation_providers

import (
	"bytes"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"strings"
	"time"

	"nugg-webauthn/core/pkg/hex"
	protocol "nugg-webauthn/core/pkg/webauthn"
	"nugg-webauthn/core/pkg/webauthn/metadata"

	"github.com/google/uuid"

	"nugg-webauthn/core/pkg/webauthn/webauthncose"
)

var packedAttestationKey = "packed"

func init() {
	protocol.RegisterAttestationFormat(packedAttestationKey, verifyPackedFormat)
}

// From §8.2. https://www.w3.org/TR/webauthn/#packed-attestation
// The packed attestation statement looks like:
//
//	packedStmtFormat = {
//	 	alg: COSEAlgorithmIdentifier,
//	 	sig: bytes,
//	 	x5c: [ attestnCert: bytes, * (caCert: bytes) ]
//	 } OR
//	 {
//	 	alg: COSEAlgorithmIdentifier, (-260 for ED256 / -261 for ED512)
//	 	sig: bytes,
//	 	ecdaaKeyId: bytes
//	 } OR
//	 {
//	 	alg: COSEAlgorithmIdentifier
//	 	sig: bytes,
//	 }
func verifyPackedFormat(att protocol.AttestationObject, clientDataHash []byte) (hex.Hash, string, []interface{}, error) {
	// Step 1. Verify that attStmt is valid CBOR conforming to the syntax defined
	// above and perform CBOR decoding on it to extract the contained fields.

	// Get the alg value - A COSEAlgorithmIdentifier containing the identifier of the algorithm
	// used to generate the attestation signature.

	alg, present := att.AttStatement["alg"].(int64)
	if !present {
		return nil, packedAttestationKey, nil, protocol.ErrAttestationFormat.WithMessage("Error retreiving alg value")
	}

	// Get the sig value - A byte string containing the attestation signature.
	sig, present := att.AttStatement["sig"].([]byte)
	if !present {
		return nil, packedAttestationKey, nil, protocol.ErrAttestationFormat.WithMessage("Error retreiving sig value")
	}

	// Step 2. If x5c is present, this indicates that the attestation type is not ECDAA.
	x5c, x509present := att.AttStatement["x5c"].([]interface{})
	if x509present {
		// Handle Basic Attestation steps for the x509 Certificate
		return handleBasicAttestation(sig, clientDataHash, att.RawAuthData, att.AuthData.AttData.AAGUID, alg, x5c)
	}

	// Step 3. If ecdaaKeyId is present, then the attestation type is ECDAA.
	// Also make sure the we did not have an x509 then
	ecdaaKeyID, ecdaaKeyPresent := att.AttStatement["ecdaaKeyId"].([]byte)
	if ecdaaKeyPresent {
		// Handle ECDAA Attestation steps for the x509 Certificate
		return handleECDAAAttesation(sig, clientDataHash, ecdaaKeyID)
	}

	// Step 4. If neither x5c nor ecdaaKeyId is present, self attestation is in use.
	return handleSelfAttestation(alg, att.AuthData.AttData.CredentialPublicKey, att.RawAuthData, clientDataHash, sig)
}

// Handle the attestation steps laid out in
func handleBasicAttestation(signature, clientDataHash, authData, aaguid []byte, alg int64, x5c []interface{}) (hex.Hash, string, []interface{}, error) {
	// Step 2.1. Verify that sig is a valid signature over the concatenation of authenticatorData
	// and clientDataHash using the attestation public key in attestnCert with the algorithm specified in alg.
	attestationType := "Packed (Basic)"

	for _, c := range x5c {
		cb, cv := c.([]byte)
		if !cv {
			return nil, attestationType, x5c, protocol.ErrAttestation.WithMessage("Error getting certificate from x5c cert chain")
		}
		ct, err := x509.ParseCertificate(cb)
		if err != nil {
			return nil, attestationType, x5c, protocol.ErrAttestationFormat.WithMessage(fmt.Sprintf("Error parsing certificate from ASN.1 data: %+v", err))
		}
		if ct.NotBefore.After(time.Now()) || ct.NotAfter.Before(time.Now()) {
			return nil, attestationType, x5c, protocol.ErrAttestationFormat.WithMessage("Cert in chain not time valid")
		}
	}

	attCertBytes, valid := x5c[0].([]byte)
	if !valid {
		return nil, attestationType, x5c, protocol.ErrAttestation.WithMessage("Error getting certificate from x5c cert chain")
	}

	signatureData := append(authData, clientDataHash...)

	attCert, err := x509.ParseCertificate(attCertBytes)
	if err != nil {
		return nil, attestationType, x5c, protocol.ErrAttestationFormat.WithMessage(fmt.Sprintf("Error parsing certificate from ASN.1 data: %+v", err))
	}

	coseAlg := webauthncose.COSEAlgorithmIdentifier(alg)
	sigAlg := webauthncose.SigAlgFromCOSEAlg(coseAlg)
	err = attCert.CheckSignature(x509.SignatureAlgorithm(sigAlg), signatureData, signature)
	if err != nil {
		return nil, attestationType, x5c, protocol.ErrInvalidAttestation.WithMessage(fmt.Sprintf("Signature validation error: %+v\n", err))
	}

	// Step 2.2 Verify that attestnCert meets the requirements in §8.2.1 Packed attestation statement certificate requirements.
	// §8.2.1 can be found here https://www.w3.org/TR/webauthn/#packed-attestation-cert-requirements

	// Step 2.2.1 (from §8.2.1) Version MUST be set to 3 (which is indicated by an ASN.1 INTEGER with value 2).
	if attCert.Version != 3 {
		return nil, attestationType, x5c, protocol.ErrAttestationCertificate.WithMessage("Attestation Certificate is incorrect version")
	}

	// Step 2.2.2 (from §8.2.1) Subject field MUST be set to:

	// 	Subject-C
	// 	ISO 3166 code specifying the country where the Authenticator vendor is incorporated (PrintableString)

	//  TODO: Find a good, useable, country code library. For now, check stringy-ness
	subjectString := strings.Join(attCert.Subject.Country, "")
	if subjectString == "" {
		return nil, attestationType, x5c, protocol.ErrAttestationCertificate.WithMessage("Attestation Certificate Country Code is invalid")
	}

	// 	Subject-O
	// 	Legal name of the Authenticator vendor (UTF8String)
	subjectString = strings.Join(attCert.Subject.Organization, "")
	if subjectString == "" {
		return nil, attestationType, x5c, protocol.ErrAttestationCertificate.WithMessage("Attestation Certificate Organization is invalid")
	}

	// 	Subject-OU
	// 	Literal string “Authenticator Attestation” (UTF8String)
	subjectString = strings.Join(attCert.Subject.OrganizationalUnit, " ")
	if subjectString != "Authenticator Attestation" {
		// TODO: Implement a return error when I'm more certain this is general practice
	}

	// 	Subject-CN
	//  A UTF8String of the vendor’s choosing
	subjectString = attCert.Subject.CommonName
	if subjectString == "" {
		return nil, attestationType, x5c, protocol.ErrAttestationCertificate.WithMessage("Attestation Certificate Common Name not set")
	}
	// TODO: And then what

	// Step 2.2.3 (from §8.2.1) If the related attestation root certificate is used for multiple authenticator models,
	// the Extension OID 1.3.6.1.4.1.45724.1.1.4 (id-fido-gen-ce-aaguid) MUST be present, containing the
	// AAGUID as a 16-byte OCTET STRING. The extension MUST NOT be marked as critical.

	idFido := asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 45724, 1, 1, 4}
	var foundAAGUID []byte
	for _, extension := range attCert.Extensions {
		if extension.Id.Equal(idFido) {
			if extension.Critical {
				return nil, attestationType, x5c, protocol.ErrInvalidAttestation.WithMessage("Attestation certificate FIDO extension marked as critical")
			}
			foundAAGUID = extension.Value
		}
	}

	// We validate the AAGUID as mentioned above
	// This is not well defined in§8.2.1 but mentioned in step 2.3: we validate the AAGUID if it is present within the certificate
	// and make sure it matches the auth data AAGUID
	// Note that an X.509 Extension encodes the DER-encoding of the value in an OCTET STRING. Thus, the
	// AAGUID MUST be wrapped in two OCTET STRINGS to be valid.
	if len(foundAAGUID) > 0 {
		unMarshalledAAGUID := []byte{}
		asn1.Unmarshal(foundAAGUID, &unMarshalledAAGUID)
		if !bytes.Equal(aaguid, unMarshalledAAGUID) {
			return nil, attestationType, x5c, protocol.ErrInvalidAttestation.WithMessage("Certificate AAGUID does not match Auth Data certificate")
		}
	}
	uuid, err := uuid.FromBytes(aaguid)

	if meta, ok := metadata.Metadata[uuid]; ok {
		for _, s := range meta.StatusReports {
			if metadata.IsUndesiredAuthenticatorStatus(metadata.AuthenticatorStatus(s.Status)) {
				return nil, attestationType, x5c, protocol.ErrInvalidAttestation.WithMessage("Authenticator with undesirable status encountered")
			}
		}

		if attCert.Subject.CommonName != attCert.Issuer.CommonName {
			var hasBasicFull = false
			for _, a := range meta.MetadataStatement.AttestationTypes {
				if metadata.AuthenticatorAttestationType(a) == metadata.AuthenticatorAttestationType(metadata.BasicFull) {
					hasBasicFull = true
				}
			}
			if !hasBasicFull {
				return nil, attestationType, x5c, protocol.ErrInvalidAttestation.WithMessage("Attestation with full attestation from authentictor that does not support full attestation")
			}
		}
	} else {
		if metadata.Conformance {
			return nil, attestationType, x5c, protocol.ErrInvalidAttestation.WithMessage("AAGUID not found in metadata during conformance testing")
		}
	}

	// Step 2.2.4 The Basic Constraints extension MUST have the CA component set to false.
	if attCert.IsCA {
		return nil, attestationType, x5c, protocol.ErrInvalidAttestation.WithMessage("Attestation certificate's Basic Constraints marked as CA")
	}

	// Note for 2.2.5 An Authority Information Access (AIA) extension with entry id-ad-ocsp and a CRL
	// Distribution Point extension [RFC5280](https://www.w3.org/TR/webauthn/#biblio-rfc5280) are
	// both OPTIONAL as the status of many attestation certificates is available through authenticator
	// metadata services. See, for example, the FIDO Metadata Service
	// [FIDOMetadataService] (https://www.w3.org/TR/webauthn/#biblio-fidometadataservice)

	// Step 2.4 If successful, return attestation type Basic and attestation trust path x5c.
	// We don't handle trust paths yet but we're done
	return nil, attestationType, x5c, nil
}

func handleECDAAAttesation(signature, clientDataHash, ecdaaKeyID []byte) (hex.Hash, string, []interface{}, error) {
	return nil, "Packed (ECDAA)", nil, protocol.ErrNotSpecImplemented
}

func handleSelfAttestation(alg int64, pubKey, authData, clientDataHash, signature []byte) (hex.Hash, string, []interface{}, error) {
	attestationType := "Packed (Self)"
	// §4.1 Validate that alg matches the algorithm of the credentialPublicKey in authenticatorData.

	// §4.2 Verify that sig is a valid signature over the concatenation of authenticatorData and
	// clientDataHash using the credential public key with alg.
	verificationData := append(authData, clientDataHash...)

	key, err := webauthncose.ParsePublicKey(pubKey)
	if err != nil {
		return nil, attestationType, nil, protocol.ErrAttestationFormat.WithMessage(fmt.Sprintf("Error parsing the public key: %+v\n", err))
	}

	switch key.(type) {
	case webauthncose.OKPPublicKeyData:
		k := key.(webauthncose.OKPPublicKeyData)
		err := verifyKeyAlgorithm(k.Algorithm, alg)
		if err != nil {
			return nil, attestationType, nil, err
		}
	case webauthncose.EC2PublicKeyData:
		k := key.(webauthncose.EC2PublicKeyData)
		err := verifyKeyAlgorithm(k.Algorithm, alg)
		if err != nil {
			return nil, attestationType, nil, err
		}
	case webauthncose.RSAPublicKeyData:
		k := key.(webauthncose.RSAPublicKeyData)
		err := verifyKeyAlgorithm(k.Algorithm, alg)
		if err != nil {
			return nil, attestationType, nil, err
		}
	default:
		return nil, attestationType, nil, protocol.ErrInvalidAttestation.WithMessage("Error verifying the public key data")
	}

	valid, err := webauthncose.VerifySignature(key, verificationData, signature)
	if !valid && err == nil {
		return nil, attestationType, nil, protocol.ErrInvalidAttestation.WithMessage("Unabled to verify signature")
	}

	return hex.BytesToHash(pubKey), attestationType, nil, err
}

func verifyKeyAlgorithm(keyAlgorithm, attestedAlgorithm int64) error {
	if keyAlgorithm != attestedAlgorithm {
		return protocol.ErrInvalidAttestation.WithMessage("Public key algorithm does not equal att statement algorithm")
	}
	return nil
}