package jwt

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

type Client struct {
	keyfunc jwt.Keyfunc
}

// Client is a wrapper for the DynamoDB client
// decode jwt token
func (client *Client) Decode(ctx context.Context, token string) (*jwt.Token, error) {
	r, err := jwt.Parse(token, client.keyfunc)
	if err != nil {
		return nil, err
	}
	return r, nil
}
