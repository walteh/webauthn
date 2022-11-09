package env

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Environment struct {
	ChallengeTableName string

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

	if env.ChallengeTableName, err = osGet("CHALLENGES_DYNAMO_TABLE_NAME"); err != nil {
		return env, err
	}

	// load aws config
	if env.AwsConfig, err = config.LoadDefaultConfig(ctx); err != nil {
		return env, err
	}

	return env, nil

}
