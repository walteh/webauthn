package signinwithapple

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateClientSecret(t *testing.T) {

	tests := []struct {
		name       string
		signingKey string
		wantSecret bool
		wantErr    bool
	}{
		{
			name:       "bad key",
			signingKey: "bad_key",
			wantSecret: false,
			wantErr:    true,
		},
		{
			name:       "bad key wrong format",
			signingKey: testWrongKey,
			wantSecret: false,
			wantErr:    true,
		},
		{
			name:       "good key",
			signingKey: testGoodKey,
			wantSecret: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			config := &ClientConfig{
				"1234567890", "com.example.app", "0987654321",
			}

			got, err := config.GenerateClientSecret(tt.signingKey)
			if !tt.wantErr {
				assert.NoError(t, err, "expected no error but got %s", err)
			}
			if tt.wantSecret {
				assert.NotEmpty(t, got, "wanted a secret string returned but got none")

				token, _, err := jwt.NewParser().ParseUnverified(got, jwt.MapClaims{})
				assert.NoError(t, err, "error while decoding the secret")

				r, b := token.Claims.(jwt.MapClaims)
				assert.True(t, b, "invalid issuer")
				assert.True(t, r.VerifyIssuer("1234567890", true))

				assert.Equal(t, "com.example.app", token.Claims.(jwt.MapClaims)["sub"])

				assert.Equal(t, string("ES256"), token.Method.Alg(), "invalid algorithm")
			}
		})
	}
}

const (
	testGoodKey = `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQg+94fs23vSrhBIXNz
OdeRb7+FJkIsVrnTSf7eIYKdf4mgCgYIKoZIzj0DAQehRANCAATyBS3eRgOJ53OQ
LFhGSJw4aiqju7muVwoIWFxCcFJasRwyGcbs0C7vt3xKV/DRJvID4UljaI53wETq
RxlkNCeV
-----END PRIVATE KEY-----` // A revoked key that can be used for testing

	testWrongKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCjcGqTkOq0CR3rTx0ZSQSIdTrDrFAYl29611xN8aVgMQIWtDB/
lD0W5TpKPuU9iaiG/sSn/VYt6EzN7Sr332jj7cyl2WrrHI6ujRswNy4HojMuqtfa
b5FFDpRmCuvl35fge18OvoQTJELhhJ1EvJ5KUeZiuJ3u3YyMnxxXzLuKbQIDAQAB
AoGAPrNDz7TKtaLBvaIuMaMXgBopHyQd3jFKbT/tg2Fu5kYm3PrnmCoQfZYXFKCo
ZUFIS/G1FBVWWGpD/MQ9tbYZkKpwuH+t2rGndMnLXiTC296/s9uix7gsjnT4Naci
5N6EN9pVUBwQmGrYUTHFc58ThtelSiPARX7LSU2ibtJSv8ECQQDWBRrrAYmbCUN7
ra0DFT6SppaDtvvuKtb+mUeKbg0B8U4y4wCIK5GH8EyQSwUWcXnNBO05rlUPbifs
DLv/u82lAkEAw39sTJ0KmJJyaChqvqAJ8guulKlgucQJ0Et9ppZyet9iVwNKX/aW
9UlwGBMQdafQ36nd1QMEA8AbAw4D+hw/KQJBANJbHDUGQtk2hrSmZNoV5HXB9Uiq
7v4N71k5ER8XwgM5yVGs2tX8dMM3RhnBEtQXXs9LW1uJZSOQcv7JGXNnhN0CQBZe
nzrJAWxh3XtznHtBfsHWelyCYRIAj4rpCHCmaGUM6IjCVKFUawOYKp5mmAyObkUZ
f8ue87emJLEdynC1CLkCQHduNjP1hemAGWrd6v8BHhE3kKtcK6KHsPvJR5dOfzbd
HAqVePERhISfN6cwZt5p8B3/JUwSR8el66DF7Jm57BM=
-----END RSA PRIVATE KEY-----` // Wrong format - this is PKCS1
)
