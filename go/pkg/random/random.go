package random

import (
	// "crypto"
	sha256 "crypto/sha256"
	"encoding/base64"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Sequence(seed string) string {

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 64)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	keccak := sha256.New()

	keccak.Write([]byte(string(b)))

	return base64.RawStdEncoding.EncodeToString(keccak.Sum([]byte(seed)))
}
