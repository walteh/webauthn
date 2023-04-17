package helpers

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"

	"git.nugg.xyz/webauthn/pkg/webauthn/cfssl/pkcs7"

	cferr "git.nugg.xyz/webauthn/pkg/webauthn/cfssl/errors"
)

// ParseCertificatePEM parses and returns a PEM-encoded certificate,
// can handle PEM encoded PKCS #7 structures.
func ParseCertificatePEM(certPEM []byte) (*x509.Certificate, error) {
	certPEM = bytes.TrimSpace(certPEM)
	cert, rest, err := ParseOneCertificateFromPEM(certPEM)
	if err != nil {
		// Log the actual parsing error but throw a default parse error message.
		log.Printf("Certificate parsing error: %v", err)
		return nil, cferr.New(cferr.CertificateError, cferr.ParseFailed)
	} else if cert == nil {
		return nil, cferr.New(cferr.CertificateError, cferr.DecodeFailed)
	} else if len(rest) > 0 {
		return nil, cferr.Wrap(cferr.CertificateError, cferr.ParseFailed, errors.New("the PEM file should contain only one object"))
	} else if len(cert) > 1 {
		return nil, cferr.Wrap(cferr.CertificateError, cferr.ParseFailed, errors.New("the PKCS7 object in the PEM file should contain only one certificate"))
	}
	return cert[0], nil
}

// ParseOneCertificateFromPEM attempts to parse one PEM encoded certificate object,
// either a raw x509 certificate or a PKCS #7 structure possibly containing
// multiple certificates, from the top of certsPEM, which itself may
// contain multiple PEM encoded certificate objects.
func ParseOneCertificateFromPEM(certsPEM []byte) ([]*x509.Certificate, []byte, error) {

	block, rest := pem.Decode(certsPEM)
	if block == nil {
		return nil, rest, nil
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		pkcs7data, err := pkcs7.ParsePKCS7(block.Bytes)
		if err != nil {
			return nil, rest, err
		}
		if pkcs7data.ContentInfo != "SignedData" {
			return nil, rest, errors.New("only PKCS #7 Signed Data Content Info supported for certificate parsing")
		}
		certs := pkcs7data.Content.SignedData.Certificates
		if certs == nil {
			return nil, rest, errors.New("PKCS #7 structure contains no certificates")
		}
		return certs, rest, nil
	}
	var certs = []*x509.Certificate{cert}
	return certs, rest, nil
}
