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
	ApplePublicKeyEndpoint        *url.URL
	AppleTokenEndpoint            *url.URL
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

	if val, err := osGet("APPLE_PUBLICKEY_ENDPOINT"); err != nil {
		return env, err
	} else {
		env.ApplePublicKeyEndpoint, err = url.Parse(val)
		if err != nil {
			return env, err
		}
	}

	if val, err := osGet("APPLE_TOKEN_ENDPOINT"); err != nil {
		return env, err
	} else {
		env.ApplePublicKeyEndpoint, err = url.Parse(val)
		if err != nil {
			return env, err
		}
	}

	if env.SignInWithApplePrivateKeyName, err = osGet("SM_SIGNINWITHAPPLE_PRIVATEKEY_NAME"); err != nil {
		return env, err
	}

	if env.SignInWithApplePrivateKeyID = os.Getenv("APPLE_KEY_ID"); err != nil {
		return env, err
	}

	if env.AppleTeamID, err = osGet("APPLE_TEAM_ID"); err != nil {
		return env, err
	}

	if env.AppleServiceName, err = osGet("APPLE_SERVICE_NAME"); err != nil {
		return env, err
	}

	if env.ChallengeTableName, err = osGet("DYNAMO_CHALLENGE_TABLE_NAME"); err != nil {
		return env, err
	}

	if env.AppleIdentityPoolId, err = osGet("COGNITO_IDENTITY_POOL_ID"); err != nil {
		return env, err
	}

	if env.AwsConfig, err = config.LoadDefaultConfig(ctx); err != nil {
		return env, err
	}

	return
}
