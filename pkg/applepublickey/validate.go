package applepublickey

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/walteh/terrors"
)

type SafeJwtToken struct {
	*jwt.Token
	// The unique identifier for the user.
}

func (client *PublicKeyResponse) BuildKeyFunc() (jwt.Keyfunc, error) {

	return func(t *jwt.Token) (any, error) {

		// check the signing method
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, terrors.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		// check the kid
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, terrors.Errorf("kid not found in token header")
		}

		// get the public key
		key := client.GetPublicKey(kid)
		if key == nil {
			return nil, terrors.Errorf("key not found for kid: %s", kid)
		}

		return key.RSA(), nil

	}, nil

}

// getClaims decodes the id_token response and returns the JWT claims to identify the user
func (r *PublicKeyResponse) ParseToken(token string) (*SafeJwtToken, error) {
	keyfunc, err := r.BuildKeyFunc()
	if err != nil {
		return nil, err
	}

	j, err := jwt.Parse(token, keyfunc)
	if err != nil {
		return nil, err
	}

	return &SafeJwtToken{j}, nil
}

// Get decodes the id_token response and returns the unique subject ID to identify the user
func (r *SafeJwtToken) GetUniqueID() (string, error) {
	if val, ok := (r.Claims).(jwt.MapClaims)["sub"].(string); ok {
		return val, nil
	} else {
		return "", terrors.Errorf("token does not contain a sub claim")
	}

}

// GetEmail decodes the id_token response and returns the email address of the user
func (r *SafeJwtToken) GetEmail() (email string, emailVerified bool, isPrivate bool, err error) {
	var ok bool

	if email, ok = (r.Claims).(jwt.MapClaims)["email"].(string); !ok {
		return "", false, false, terrors.Errorf("could not get email from token")
	}

	if emailVerified, ok = (r.Claims).(jwt.MapClaims)["email_verified"].(bool); !ok {
		return email, false, false, terrors.Errorf("could not get email from token")
	}

	if isPrivate, ok = (r.Claims).(jwt.MapClaims)["is_private_email"].(bool); !ok {
		return email, emailVerified, false, terrors.Errorf("could not get email from token")
	}

	return email, emailVerified, isPrivate, nil
}
