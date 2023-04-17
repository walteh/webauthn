package errors

import "git.nugg.xyz/go-sdk/errors"

var (
	ErrBadRequest             = errors.New("invalid_request").WithMessage("Error reading the request data")
	ErrChallengeMismatch      = errors.New("challenge_mismatch").WithMessage("Stored challenge and received challenge do not match")
	ErrOriginMismatch         = errors.New("origin_mismatch").WithMessage("Stored origin and received origin do not match")
	ErrParsingData            = errors.New("parse_error").WithMessage("Error parsing the authenticator response")
	ErrAuthData               = errors.New("auth_data").WithMessage("Error verifying the authenticator data")
	ErrVerification           = errors.New("verification_error").WithMessage("Error validating the authenticator response")
	ErrAttestation            = errors.New("attesation_error").WithMessage("Error validating the attestation data provided")
	ErrInvalidAttestation     = errors.New("invalid_attestation").WithMessage("Invalid attestation data")
	ErrAttestationFormat      = errors.New("invalid_attestation").WithMessage("Invalid attestation format")
	ErrAttestationCertificate = errors.New("invalid_certificate").WithMessage("Invalid attestation certificate")
	ErrAssertionSignature     = errors.New("invalid_signature").WithMessage("Assertion Signature against auth data and client hash is not valid")
	ErrUnsupportedKey         = errors.New("invalid_key_type").WithMessage("Unsupported Public Key Type")
	ErrUnsupportedAlgorithm   = errors.New("unsupported_key_algorithm").WithMessage("Unsupported public key algorithm")
	ErrNotSpecImplemented     = errors.New("spec_unimplemented").WithMessage("This field is not yet supported by the WebAuthn spec")
	ErrNotImplemented         = errors.New("not_implemented").WithMessage("This field is not yet supported by this library")
	Err0x66CborDecode         = errors.New("cbor_decode").WithMessage("Error decoding CBOR data")
	Err0x67InvalidInput       = errors.New("invalid_input").WithMessage("Invalid input")
)
