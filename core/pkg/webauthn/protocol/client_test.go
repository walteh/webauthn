package protocol

import (
	"net/url"
	"nugg-webauthn/core/pkg/hex"
	"testing"
)

func setupCollectedClientData(challenge []byte) *CollectedClientData {
	ccd := &CollectedClientData{
		Type:   CreateCeremony,
		Origin: "example.com",
	}

	ccd.Challenge = hex.BytesToHash(challenge)
	return ccd
}

func TestVerifyCollectedClientData(t *testing.T) {
	newChallenge, err := CreateChallenge()
	if err != nil {
		t.Fatalf("error creating challenge: %s", err)
	}

	ccd := setupCollectedClientData(newChallenge)
	var storedChallenge = newChallenge

	originURL, _ := url.Parse(ccd.Origin)
	err = ccd.Verify(storedChallenge, ccd.Type, FullyQualifiedOrigin(originURL))
	if err != nil {
		t.Fatalf("error verifying challenge: expected %#v got %#v", (ccd.Challenge), storedChallenge)
	}
}

func TestVerifyCollectedClientDataIncorrectChallenge(t *testing.T) {
	newChallenge, err := CreateChallenge()
	if err != nil {
		t.Fatalf("error creating challenge: %s", err)
	}
	ccd := setupCollectedClientData(newChallenge)
	bogusChallenge, err := CreateChallenge()
	if err != nil {
		t.Fatalf("error creating challenge: %s", err)
	}
	storedChallenge := (bogusChallenge)
	err = ccd.Verify(storedChallenge, ccd.Type, ccd.Origin)
	if err == nil {
		t.Fatalf("error expected but not received. expected %#v got %#v", (ccd.Challenge), storedChallenge)
	}
}
