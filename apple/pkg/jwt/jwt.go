package jwt

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type Client struct {
	keyfunc jwt.Keyfunc
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
func (client *Client) Verify(ctx context.Context, token string) (*Jwt, error) {
	r, err := jwt.Parse(token, client.keyfunc)
	if err != nil {
		return nil, err
	}
	return &Jwt{r}, nil
}
