package protocol

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"sync"
)

// ChallengeLength - Length of bytes to generate for a challenge
const ChallengeLength = 16

// Challenge that should be signed and returned by the authenticator
type Challenge URLEncodedBase64

var rander = rand.Reader

var mu = &sync.Mutex{}

// Create a new challenge to be sent to the authenticator. The spec recommends using
// at least 16 bytes with 100 bits of entropy. We use 32 bytes.
func CreateChallenge() (Challenge, error) {
	mu.Lock()
	challenge, err := io.ReadAll(io.LimitReader(rander, ChallengeLength))
	mu.Unlock()
	if err != nil {
		return nil, err
	}

	return challenge, nil
}

func (c Challenge) String() string {
	return base64.RawURLEncoding.EncodeToString(c)
}

// UnmarshalJSON base64 decodes a URL-encoded value, storing the result in the
// provided byte slice.
func (dest *Challenge) UnmarshalJSON(data []byte) error {
	return (*URLEncodedBase64)(dest).UnmarshalJSON(data)
}

// MarshalJSON base64 encodes a non URL-encoded value, storing the result in the
// provided byte slice.
func (data Challenge) MarshalJSON() ([]byte, error) {
	return URLEncodedBase64(data).MarshalJSON()
}
