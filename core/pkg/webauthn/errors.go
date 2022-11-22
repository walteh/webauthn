package protocol

import errors "nugg-webauthn/core/pkg/errors"

var (
	ErrBadRequest             = errors.NewError(0x32).WithType("invalid_request").WithMessage("Error reading the request data")
	ErrChallengeMismatch      = errors.NewError(0x33).WithType("challenge_mismatch").WithMessage("Stored challenge and received challenge do not match")
	ErrOriginMismatch         = errors.NewError(0x33).WithType("origin_mismatch").WithMessage("Stored origin and received origin do not match")
	ErrParsingData            = errors.NewError(0x33).WithType("parse_error").WithMessage("Error parsing the authenticator response")
	ErrAuthData               = errors.NewError(0x33).WithType("auth_data").WithMessage("Error verifying the authenticator data")
	ErrVerification           = errors.NewError(0x33).WithType("verification_error").WithMessage("Error validating the authenticator response")
	ErrAttestation            = errors.NewError(0x33).WithType("attesation_error").WithMessage("Error validating the attestation data provided")
	ErrInvalidAttestation     = errors.NewError(0x33).WithType("invalid_attestation").WithMessage("Invalid attestation data")
	ErrAttestationFormat      = errors.NewError(0x33).WithType("invalid_attestation").WithMessage("Invalid attestation format")
	ErrAttestationCertificate = errors.NewError(0x33).WithType("invalid_certificate").WithMessage("Invalid attestation certificate")
	ErrAssertionSignature     = errors.NewError(0x33).WithType("invalid_signature").WithMessage("Assertion Signature against auth data and client hash is not valid")
	ErrUnsupportedKey         = errors.NewError(0x33).WithType("invalid_key_type").WithMessage("Unsupported Public Key Type")
	ErrUnsupportedAlgorithm   = errors.NewError(0x33).WithType("unsupported_key_algorithm").WithMessage("Unsupported public key algorithm")
	ErrNotSpecImplemented     = errors.NewError(0x33).WithType("spec_unimplemented").WithMessage("This field is not yet supported by the WebAuthn spec")
	ErrNotImplemented         = errors.NewError(0x33).WithType("not_implemented").WithMessage("This field is not yet supported by this library")
)
