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
	// AppsyncSchema             *graphql.Schema
	ChallengeTableName        string
	AppleIdentityPoolName     string
	AppleJwtPublicKeyEndpoint *url.URL

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

	// if val, err := osGet("APPSYNC_SCHEMA"); err != nil {
	// 	return env, err
	// } else {
	// 	// parse schema
	// 	env.AppsyncSchema = graphql.MustParseSchema(val, nil)
	// 	if err != nil {
	// 		return env, err
	// 	}
	// }

	if env.ChallengeTableName, err = osGet("CHALLENGE_TABLE_NAME"); err != nil {
		return env, err
	}

	if env.AppleIdentityPoolName, err = osGet("APPLE_IDENTITY_POOL_NAME"); err != nil {
		return env, err
	}

	if env.AwsConfig, err = config.LoadDefaultConfig(ctx); err != nil {
		return env, err
	}

	return
}
