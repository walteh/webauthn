package applepublickey

import "github.com/golang-jwt/jwt/v4"

func MockSafeJwtWithClaims(claims jwt.MapClaims) *SafeJwtToken {
	return &SafeJwtToken{
		Token: &jwt.Token{
			Claims: claims,
		},
	}
}
