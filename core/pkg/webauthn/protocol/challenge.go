package protocol

import (
	"encoding/base64"
	"nugg-auth/core/pkg/safeid"
)

// ChallengeLength - Length of bytes to generate for a challenge
const ChallengeLength = safeid.EncodedSize

// Challenge that should be signed and returned by the authenticator
type Challenge []byte

// Create a new challenge to be sent to the authenticator. The spec recommends using
// at least 16 bytes with 100 bits of entropy. We use 32 bytes.
func CreateChallenge() (Challenge, error) {
	return safeid.Make().Bytes(), nil
}

// Verify that the challenge is valid
func (c Challenge) String() string {
	return base64.RawURLEncoding.EncodeToString(c)
}

// func (c Challenge) SafeID() (safeid.SafeID, error) {
// 	return safeid.ParseBytesStrict(c)
// }
