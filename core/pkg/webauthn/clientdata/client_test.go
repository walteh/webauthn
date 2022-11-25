package clientdata

import (
	"net/url"

	"testing"

	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/webauthn/challenge"
	"nugg-webauthn/core/pkg/webauthn/types"
)

func setupCollectedClientData(challenge []byte) types.CollectedClientData {
	ccd := types.CollectedClientData{
		Type:   types.CreateCeremony,
		Origin: "example.com",
	}

	ccd.Challenge = hex.BytesToHash(challenge)
	return ccd
}

func TestVerifyCollectedClientData(t *testing.T) {
	newChallenge, err := challenge.CreateChallenge()
	if err != nil {
		t.Fatalf("error creating challenge: %s", err)
	}

	ccd := setupCollectedClientData(newChallenge)
	var storedChallenge = newChallenge

	originURL, _ := url.Parse(ccd.Origin)
	err = Verify(types.VerifyClientDataArgs{
		ClientData:         ccd,
		StoredChallenge:    storedChallenge,
		CeremonyType:       ccd.Type,
		RelyingPartyOrigin: types.FullyQualifiedOrigin(originURL)})
	if err != nil {
		t.Fatalf("error verifying challenge: expected %#v got %#v", (ccd.Challenge), storedChallenge)
	}
}

func TestVerifyCollectedClientDataIncorrectChallenge(t *testing.T) {
	newChallenge, err := challenge.CreateChallenge()
	if err != nil {
		t.Fatalf("error creating challenge: %s", err)
	}
	ccd := setupCollectedClientData(newChallenge)
	bogusChallenge, err := challenge.CreateChallenge()
	if err != nil {
		t.Fatalf("error creating challenge: %s", err)
	}
	storedChallenge := (bogusChallenge)
	err = Verify(types.VerifyClientDataArgs{
		ClientData:         ccd,
		StoredChallenge:    storedChallenge,
		CeremonyType:       ccd.Type,
		RelyingPartyOrigin: ccd.Origin,
	})
	if err == nil {
		t.Fatalf("error expected but not received. expected %#v got %#v", (ccd.Challenge), storedChallenge)
	}
}
