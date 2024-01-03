package clientdata_test

import (
	"context"
	"net/url"

	"testing"

	"github.com/rs/zerolog"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/challenge"
	"github.com/walteh/webauthn/pkg/webauthn/clientdata"
	"github.com/walteh/webauthn/pkg/webauthn/types"
)

func setupCollectedClientData(challenge []byte) types.CollectedClientData {
	ccd := types.CollectedClientData{
		Type:   types.CreateCeremony,
		Origin: "example.com",
	}

	ccd.Challenge = types.CeremonyID(hex.BytesToHash(challenge))
	return ccd
}

func TestVerifyCollectedClientData(t *testing.T) {

	ctx := zerolog.New(zerolog.NewConsoleWriter()).Level(zerolog.TraceLevel).With().Caller().Logger().WithContext(context.Background())

	newChallenge, err := challenge.CreateChallenge()
	if err != nil {
		t.Fatalf("error creating challenge: %s", err)
	}

	ccd := setupCollectedClientData(newChallenge)
	var storedChallenge = newChallenge

	originURL, _ := url.Parse(ccd.Origin)
	err = clientdata.Verify(ctx, types.VerifyClientDataArgs{
		ClientData:         ccd,
		StoredChallenge:    types.CeremonyID(storedChallenge),
		CeremonyType:       ccd.Type,
		RelyingPartyOrigin: types.FullyQualifiedOrigin(originURL)})
	if err != nil {
		t.Fatalf("error verifying challenge: expected %#v got %#v", (ccd.Challenge), storedChallenge)
	}
}

func TestVerifyCollectedClientDataIncorrectChallenge(t *testing.T) {
	ctx := zerolog.New(zerolog.NewConsoleWriter()).Level(zerolog.TraceLevel).With().Caller().Logger().WithContext(context.Background())

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
	err = clientdata.Verify(ctx, types.VerifyClientDataArgs{
		ClientData:         ccd,
		StoredChallenge:    types.CeremonyID(storedChallenge),
		CeremonyType:       ccd.Type,
		RelyingPartyOrigin: ccd.Origin,
	})
	if err == nil {
		t.Fatalf("error expected but not received. expected %#v got %#v", (ccd.Challenge), storedChallenge)
	}
}
