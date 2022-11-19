package protocol

import (
	"crypto/rand"
	"crypto/sha256"
	"io"
	"testing"
)

type MockHashDeterministicReader struct {
	initial []byte
	last    []byte
}

func NewDeterministic(str string) io.Reader {
	return &MockHashDeterministicReader{
		initial: []byte(str),
		last:    []byte(str),
	}
}

func (f *MockHashDeterministicReader) CalculateDeterministicHash(hashes int) []byte {
	rns := NewDeterministic(string(f.initial))
	var dummy []byte
	for i := 0; i < hashes; i++ {
		rns.Read(dummy)
	}

	return rns.(*MockHashDeterministicReader).last
}

func (m *MockHashDeterministicReader) Read(p []byte) (n int, err error) {

	hash := sha256.Sum256(m.last)
	m.last = hash[:][:ChallengeLength]
	copy(p, m.last)

	return len(m.last), nil
}

func MockSetRander(t *testing.T, str string) *MockHashDeterministicReader {
	t.Helper()
	t.Cleanup(func() {
		rander = rand.Reader
	})

	rander = NewDeterministic(str)

	return rander.(*MockHashDeterministicReader)
}
