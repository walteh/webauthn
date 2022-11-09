package signinwithapple

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

/*
GenerateClientSecret generates the client secret used to make requests to the validation server.
The secret expires after 6 months
signingKey - Private key from Apple obtained by going to the keys section of the developer section
teamID - Your 10-character Team ID
clientID - Your Services ID, e.g. com.aaronparecki.services
keyID - Find the 10-char Key ID value from the portal
*/
func (config *ClientConfig) GenerateClientSecret(secret string) (string, error) {
	block, _ := pem.Decode([]byte(secret))
	if block == nil {
		return "", errors.New("empty block after decoding")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Create the Claims
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		Issuer:    config.TeamID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.AddDate(0, 6, 0)), // 6 months
		Audience:  jwt.ClaimStrings{"https://appleid.apple.com"},
		Subject:   config.ClientID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["alg"] = "ES256"
	token.Header["kid"] = config.KeyID

	return token.SignedString(privKey)
}
