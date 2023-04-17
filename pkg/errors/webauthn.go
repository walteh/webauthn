package errors

import "github.com/nuggxyz/golang/errors"

var (
	ErrBadRequest             = errors.NewError(0x32).WithType("invalid_request").WithInfo("Error reading the request data")
	ErrChallengeMismatch      = errors.NewError(0x33).WithType("challenge_mismatch").WithInfo("Stored challenge and received challenge do not match")
	ErrOriginMismatch         = errors.NewError(0x33).WithType("origin_mismatch").WithInfo("Stored origin and received origin do not match")
	ErrParsingData            = errors.NewError(0x33).WithType("parse_error").WithInfo("Error parsing the authenticator response")
	ErrAuthData               = errors.NewError(0x33).WithType("auth_data").WithInfo("Error verifying the authenticator data")
	ErrVerification           = errors.NewError(0x33).WithType("verification_error").WithInfo("Error validating the authenticator response")
	ErrAttestation            = errors.NewError(0x33).WithType("attesation_error").WithInfo("Error validating the attestation data provided")
	ErrInvalidAttestation     = errors.NewError(0x33).WithType("invalid_attestation").WithInfo("Invalid attestation data")
	ErrAttestationFormat      = errors.NewError(0x33).WithType("invalid_attestation").WithInfo("Invalid attestation format")
	ErrAttestationCertificate = errors.NewError(0x33).WithType("invalid_certificate").WithInfo("Invalid attestation certificate")
	ErrAssertionSignature     = errors.NewError(0x33).WithType("invalid_signature").WithInfo("Assertion Signature against auth data and client hash is not valid")
	ErrUnsupportedKey         = errors.NewError(0x33).WithType("invalid_key_type").WithInfo("Unsupported Public Key Type")
	ErrUnsupportedAlgorithm   = errors.NewError(0x33).WithType("unsupported_key_algorithm").WithInfo("Unsupported public key algorithm")
	ErrNotSpecImplemented     = errors.NewError(0x33).WithType("spec_unimplemented").WithInfo("This field is not yet supported by the WebAuthn spec")
	ErrNotImplemented         = errors.NewError(0x33).WithType("not_implemented").WithInfo("This field is not yet supported by this library")
	Err0x66CborDecode         = errors.NewError(0x66).WithType("cbor_decode").WithInfo("Error decoding CBOR data")
	Err0x67InvalidInput       = errors.NewError(0x67).WithType("invalid_input").WithInfo("Invalid input")
)
