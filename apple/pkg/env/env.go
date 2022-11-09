package env

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Environment struct {
	// evnironment variables
	SignInWithApplePrivateKeyName string
	SignInWithApplePrivateKeyID   string
	ChallengeTableName            string
	AppleIdentityPoolId           string
	AppleJwtPublicKeyEndpoint     *url.URL
	AppleTeamID                   string
	AppleServiceName              string

	// aws config
	AwsConfig aws.Config
}

func osGet(key string) (value string, err error) {
	if value = os.Getenv(key); value == "" {
		return "", fmt.Errorf("env variable " + key + " is empty")
	}
	return value, nil
}

func NewEnv(ctx context.Context) (env Environment, err error) {

	if val, err := osGet("APPLE_JWT_PUBLIC_KEY_ENDPOINT"); err != nil {
		return env, err
	} else {
		// parse url
		env.AppleJwtPublicKeyEndpoint, err = url.Parse(val)
		if err != nil {
			return env, err
		}

	}

	if env.SignInWithApplePrivateKeyName, err = osGet("SIGN_IN_WITH_APPLE_PRIVATE_KEY_NAME"); err != nil {
		return env, err
	}

	if env.SignInWithApplePrivateKeyID = os.Getenv("SIGN_IN_WITH_APPLE_PRIVATE_KEY_ID"); err != nil {
		return env, err
	}

	if env.AppleTeamID, err = osGet("APPLE_TEAM_ID"); err != nil {
		return env, err
	}

	if env.AppleServiceName, err = osGet("APPLE_SERVICE_NAME"); err != nil {
		return env, err
	}

	if env.ChallengeTableName, err = osGet("CHALLENGE_TABLE_NAME"); err != nil {
		return env, err
	}

	if env.AppleIdentityPoolId, err = osGet("APPLE_IDENTITY_POOL_ID"); err != nil {
		return env, err
	}

	if env.AwsConfig, err = config.LoadDefaultConfig(ctx); err != nil {
		return env, err
	}

	return
}
