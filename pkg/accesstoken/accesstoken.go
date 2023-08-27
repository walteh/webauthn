package accesstoken

import "context"

type AccessTokenProvider interface {
	AccessTokenForUserID(ctx context.Context, userID string) (string, error)
}
