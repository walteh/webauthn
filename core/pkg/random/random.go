package random

import (
	"crypto/rand"
	"sync"

	"github.com/oklog/ulid/v2"
	"github.com/segmentio/ksuid"
)

// var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// func Sequence(seed string) string {

// 	rand.Seed(time.Now().UnixNano())

// 	b := make([]rune, 64)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}

// 	keccak := sha256.New()

// 	keccak.Write([]byte(string(b)))

// 	return base64.RawStdEncoding.EncodeToString(keccak.Sum([]byte(seed)))
// }

func KSUID() ksuid.KSUID {
	return ksuid.New()
}

var randMutex = sync.Mutex{}

var rander = rand.Reader

func ULID() (ksuid ulid.ULID) {

	randMutex.Lock()

	ulid := ulid.MustNew(ulid.Now(), rander)

	randMutex.Unlock()

	return ulid

}
