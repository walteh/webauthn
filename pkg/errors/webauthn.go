package errors

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	ErrBadRequest             = errors.Wrap(errors.New("invalid_request"), "Error reading the request data")
	ErrChallengeMismatch      = errors.Wrap(errors.New("challenge_mismatch"), "Stored challenge and received challenge do not match")
	ErrOriginMismatch         = errors.Wrap(errors.New("origin_mismatch"), "Stored origin and received origin do not match")
	ErrParsingData            = errors.Wrap(errors.New("parse_error"), "Error parsing the authenticator response")
	ErrAuthData               = errors.Wrap(errors.New("auth_data"), "Error verifying the authenticator data")
	ErrVerification           = errors.Wrap(errors.New("verification_error"), "Error validating the authenticator response")
	ErrAttestation            = errors.Wrap(errors.New("attesation_error"), "Error validating the attestation data provided")
	ErrInvalidAttestation     = errors.Wrap(errors.New("invalid_attestation"), "Invalid attestation data")
	ErrAttestationFormat      = errors.Wrap(errors.New("invalid_attestation"), "Invalid attestation format")
	ErrAttestationCertificate = errors.Wrap(errors.New("invalid_certificate"), "Invalid attestation certificate")
	ErrAssertionSignature     = errors.Wrap(errors.New("invalid_signature"), "Assertion Signature against auth data and client hash is not valid")
	ErrUnsupportedKey         = errors.Wrap(errors.New("invalid_key_type"), "Unsupported Public Key Type")
	ErrUnsupportedAlgorithm   = errors.Wrap(errors.New("unsupported_key_algorithm"), "Unsupported public key algorithm")
	ErrNotSpecImplemented     = errors.Wrap(errors.New("spec_unimplemented"), "This field is not yet supported by the WebAuthn spec")
	ErrNotImplemented         = errors.Wrap(errors.New("not_implemented"), "This field is not yet supported by this library")
	Err0x66CborDecode         = errors.Wrap(errors.New("cbor_decode"), "Error decoding CBOR data")
	Err0x67InvalidInput       = errors.Wrap(errors.New("invalid_input"), "Invalid input")
)

func Wrap(ctx context.Context, e error, s ...string) error {
	event := zerolog.Ctx(ctx).Error().Err(e).CallerSkipFrame(1)
	for i, msg := range s {
		event = event.Str(fmt.Sprintf("extra[%d]", i), msg)
	}
	zerolog.Ctx(ctx).Error().Err(e).CallerSkipFrame(1).Msg("error")
	return e
}
