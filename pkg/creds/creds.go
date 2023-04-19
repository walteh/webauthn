package creds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

type CredentialProvider interface {
	Get(ctx context.Context, challenge string, cred string) (*types.Credentials, error)
}
