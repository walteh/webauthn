package accesstoken

import "context"

type Provider interface {
	AccessTokenForUserID(ctx context.Context, userID string) (string, error)
}
