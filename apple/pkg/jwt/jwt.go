package jwt

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type Client interface {
	// GetItem gets an item from the table
	BuildKeyFunc() (jwt.Keyfunc, error)
}

type Jwt struct {
	*jwt.Token
}

func (me *Jwt) MapClaims() (jwt.MapClaims, error) {
	if val, ok := me.Claims.(jwt.MapClaims); ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("unable to parse claims to map")
	}
}

func (me *Jwt) Sub() (string, error) {
	return me.getClaimValue("sub")
}

func (me *Jwt) getClaimValue(str string) (string, error) {
	claims, err := me.MapClaims()
	if err != nil {
		return "", err
	}
	if val, ok := claims[str].(string); ok {
		if val == "" {
			return "", fmt.Errorf("claim %s is empty", str)
		}
		return val, nil
	} else {
		return "", fmt.Errorf("unable to find %s", str)
	}
}

// Client is a wrapper for the DynamoDB client
// decode jwt token
func VerifyJwt(ctx context.Context, token string, client Client) (*Jwt, error) {

	keyfunc, err := client.BuildKeyFunc()
	if err != nil {
		return nil, err
	}

	r, err := jwt.Parse(token, keyfunc)
	if err != nil {
		return nil, err
	}
	return &Jwt{r}, nil
}
