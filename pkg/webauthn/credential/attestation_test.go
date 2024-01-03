package credential_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/credential"
	"github.com/walteh/webauthn/pkg/webauthn/providers"
	"github.com/walteh/webauthn/pkg/webauthn/types"
)

func TestAttestationVerify(t *testing.T) {
	for i := range testAttestationOptions {
		t.Run(fmt.Sprintf("Running test %d", i), func(t *testing.T) {

			ctx := context.Background()

			arg := testAttestationResponses[i]

			opt := testAttestationOptions[i]

			args := types.VerifyAttestationInputArgs{
				Provider:           providers.NewNoneAttestationProvider(),
				Input:              arg,
				StoredChallenge:    opt.challenge,
				SessionId:          hex.Hash{0, 2, 3},
				VerifyUser:         false,
				RelyingPartyID:     opt.relyingPartyId,
				RelyingPartyOrigin: opt.relyingPartyName,
			}

			// Test Base Verification
			_, err := credential.VerifyAttestationInput(ctx, args)
			require.NoError(t, err)
		})
	}
}

func TestPackedAttestationVerification(t *testing.T) {
	t.Run("Testing Self Packed", func(t *testing.T) {
		arg := testAttestationResponses[0]

		ctx := context.Background()

		abc, err := credential.ParseAttestationInput(ctx, arg)
		if err != nil {
			t.Fatal(err)
		}

		// Test Packed Verification
		// Unpack args
		clientDataHash := sha256.Sum256([]byte(arg.UTF8ClientDataJSON))

		packed := providers.NewPackedAttestationProvider()

		_, _, _, err = packed.Attest(*abc, clientDataHash[:])
		require.NoError(t, err)
	})
}

var testAttestationResponses = []types.AttestationInput{
	{
		CredentialID:       types.CredentialID(hex.MustBase64ToHash("AOx6vFGGITtlwjhqFFvAkJmBzSzfwE1dBa1fVR_Ltq5L35FJRNdgkXe84v3-0TEVNCSp")),
		CredentialType:     "public-key",
		UTF8ClientDataJSON: "{\"challenge\":\"rWiex8xDOPfiCgyFu4BLW6vVOmXKgPwHrlMCgEs9SBA\",\"origin\":\"http://localhost:9005\",\"type\":\"webauthn.create\"}",
		AttestationObject:  hex.MustBase64ToHash("o2NmbXRmcGFja2VkZ2F0dFN0bXSiY2FsZyZjc2lnWEcwRQIhAJgdgw5x8JzE4JfR6x1RBO8eCHNE8eW_L1VTV03zpyL5AiBv8eUzua3XSS3bPYC7m8eXzJhcaRyeGe7UcuqIrDSvC2hhdXRoRGF0YVi3SZYN5YgOjGh0NBcPZHZgW4_krrmihjLHmVzzuoMdl2NFXJE5zK3OAAI1vMYKZIsLJfHwVQMAMwDserxRhiE7ZcI4ahRbwJCZgc0s38BNXQWtX1Ufy7auS9-RSUTXYJF3vOL9_tExFTQkqaUBAgMmIAEhWCCm9OYidwiIoH9SwVQqUAnH8Gj5ZJ2_qr8gjbg41q4M1SJYIA07XKpHSgS1mE7R1MjotVIQqyHi9WAxGwHQsCteVK2V"),
	},
	{
		CredentialID:       types.CredentialID(hex.MustBase64ToHash("FOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmg")),
		CredentialType:     "public-key",
		UTF8ClientDataJSON: "{\"challenge\":\"-Ri5NZTzJ8b6mvW3TVScLotEoALfgBa2Bn4YSaIObHc\",\"origin\":\"https://webauthn.io\",\"type\":\"webauthn.create\"}",
		AttestationObject:  hex.MustBase64ToHash("o2NmbXRoZmlkby11MmZnYXR0U3RtdKJjc2lnWEYwRAIgfyIhwZj-fkEVyT1GOK8chDHJR2chXBLSRg6bTCjODmwCIHH6GXI_BQrcR-GHg5JfazKVQdezp6_QWIFfT4ltTCO2Y3g1Y4FZAlMwggJPMIIBN6ADAgECAgQSNtF_MA0GCSqGSIb3DQEBCwUAMC4xLDAqBgNVBAMTI1l1YmljbyBVMkYgUm9vdCBDQSBTZXJpYWwgNDU3MjAwNjMxMCAXDTE0MDgwMTAwMDAwMFoYDzIwNTAwOTA0MDAwMDAwWjAxMS8wLQYDVQQDDCZZdWJpY28gVTJGIEVFIFNlcmlhbCAyMzkyNTczNDEwMzI0MTA4NzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABNNlqR5emeDVtDnA2a-7h_QFjkfdErFE7bFNKzP401wVE-QNefD5maviNnGVk4HJ3CsHhYuCrGNHYgTM9zTWriGjOzA5MCIGCSsGAQQBgsQKAgQVMS4zLjYuMS40LjEuNDE0ODIuMS41MBMGCysGAQQBguUcAgEBBAQDAgUgMA0GCSqGSIb3DQEBCwUAA4IBAQAiG5uzsnIk8T6-oyLwNR6vRklmo29yaYV8jiP55QW1UnXdTkEiPn8mEQkUac-Sn6UmPmzHdoGySG2q9B-xz6voVQjxP2dQ9sgbKd5gG15yCLv6ZHblZKkdfWSrUkrQTrtaziGLFSbxcfh83vUjmOhDLFC5vxV4GXq2674yq9F2kzg4nCS4yXrO4_G8YWR2yvQvE2ffKSjQJlXGO5080Ktptplv5XN4i5lS-AKrT5QRVbEJ3B4g7G0lQhdYV-6r4ZtHil8mF4YNMZ0-RaYPxAaYNWkFYdzOZCaIdQbXRZefgGfbMUiAC2gwWN7fiPHV9eu82NYypGU32OijG9BjhGt_aGF1dGhEYXRhWMR0puqSE8mcL3SyJJKzIM9AJiqUwalQoDl_KSULYIQe8EEAAAAAAAAAAAAAAAAAAAAAAAAAAABAFOxcmsqPLNCHtyILvbNkrtHMdKAeqSJXYZDbeFd0kc5Enm8Kl6a0Jp0szgLilDw1S4CjZhe9Z2611EUGbjyEmqUBAgMmIAEhWCD_ap3Q9zU8OsGe967t48vyRxqn8NfFTk307mC1WsH2ISJYIIcqAuW3MxhU0uDtaSX8-Ftf_zeNJLdCOEjZJGHsrLxH"),
	},
	{
		CredentialID:       types.CredentialID(hex.MustBase64ToHash("6Jry73M_WVWDoXLsGxRsBVVHpPWDpNy1ETGXUEvJLdTAn5Ew6nDGU6W8iO3ZkcLEqr-CBwvx0p2WAxzt8RiwQQ")),
		CredentialType:     "public-key",
		UTF8ClientDataJSON: "{\"challenge\":\"sVt4ScceMzqFSnfAq8hgLzblvo3fa4_aFVEcIESHIJ0\",\"origin\":\"https://webauthn.io\",\"type\":\"webauthn.create\"}",
		AttestationObject:  hex.MustBase64ToHash("o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YVjEdKbqkhPJnC90siSSsyDPQCYqlMGpUKA5fyklC2CEHvBBAAAAAAAAAAAAAAAAAAAAAAAAAAAAQOia8u9zP1lVg6Fy7BsUbAVVR6T1g6TctRExl1BLyS3UwJ-RMOpwxlOlvIjt2ZHCxKq_ggcL8dKdlgMc7fEYsEGlAQIDJiABIVgg--n_QvZithDycYmnifk6vMHiwBP6kugn2PlsnvkrcSgiWCBAlBYm2B-rMtQlp5MxGTLoGDHoktxb0p364Hy2BH9U2Q"),
	},
}

// 1 and 3 will pass with packed too, not just none
var testAttestationOptions = []struct {
	challenge        types.CeremonyID
	relyingPartyId   string
	relyingPartyName string
}{

	{
		challenge:        types.CeremonyID(hex.MustBase64ToHash("rWiex8xDOPfiCgyFu4BLW6vVOmXKgPwHrlMCgEs9SBA=")),
		relyingPartyId:   "localhost",
		relyingPartyName: "http://localhost:9005",
	},
	{
		challenge:        types.CeremonyID(hex.MustBase64ToHash("+Ri5NZTzJ8b6mvW3TVScLotEoALfgBa2Bn4YSaIObHc=")),
		relyingPartyId:   "webauthn.io",
		relyingPartyName: "https://webauthn.io",
	},

	{
		challenge:        types.CeremonyID(hex.MustBase64ToHash("sVt4ScceMzqFSnfAq8hgLzblvo3fa4/aFVEcIESHIJ0=")),
		relyingPartyId:   "webauthn.io",
		relyingPartyName: "https://webauthn.io",
	},
}
