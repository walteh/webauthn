package challenge

import (
	"crypto/rand"
	"io"
	"sync"

	"github.com/walteh/webauthn/pkg/hex"
)

// ChallengeLength - Length of bytes to generate for a challenge
const ChallengeLength = 16

// Challenge that should be signed and returned by the authenticator
// type Challenge hex.Hash

var rander = rand.Reader

var mu = &sync.Mutex{}

// Create a new challenge to be sent to the authenticator. The spec recommends using
// at least 16 bytes with 100 bits of entropy. We use 32 bytes.
func CreateChallenge() (hex.Hash, error) {
	mu.Lock()
	challenge, err := io.ReadAll(io.LimitReader(rander, ChallengeLength))
	mu.Unlock()
	if err != nil {
		return nil, err
	}

	x := hex.BytesToHash(challenge).Sha256()

	return x, nil
}
