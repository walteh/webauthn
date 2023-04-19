package env

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = time.StampMicro
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC | log.Lmicroseconds)
}

type Environment struct {
	AwsConfig aws.Config
}

func osGet(key string) (value string, err error) {
	if value = os.Getenv(key); value == "" {
		return "", fmt.Errorf("env variable " + key + " is empty")
	}
	return value, nil
}

func osMustGet(key string) (value string) {
	if value, err := osGet(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

func osMustGetUrl(key string) (value *url.URL) {
	if value, err := osGetUrl(key); err != nil {
		panic(err)
	} else {
		return value
	}
}
func osGetUrl(key string) (value *url.URL, err error) {
	if value, err = url.Parse(os.Getenv(key)); err != nil {
		return nil, err
	}
	return value, nil
}

func SignInWithApplePrivateKeyID() string { return osMustGet("SIGNIN_WITH_APPLE_PRIVATE_KEY_ID") }

func AppleTeamID() string { return osMustGet("APPLE_TEAM_ID") }

func AppleServiceName() string { return osMustGet("APPLE_SERVICE_NAME") }

func ApplePublicKeyEndpoint() *url.URL { return osMustGetUrl("APPLE_PUBLICKEY_ENDPOINT") }

func AppleTokenEndpoint() *url.URL { return osMustGetUrl("APPLE_TOKEN_ENDPOINT") }

func AppleIdentityPoolId() string { return osMustGet("COGNITO_IDENTITY_POOL_ID") }

func DynamoUsersTableName() string { return osMustGet("DYNAMO_USERS_TABLE_NAME") }

func DynamoCeremoniesTableName() string { return osMustGet("DYNAMO_CEREMONIES_TABLE_NAME") }

func DynamoCredentialsTableName() string { return osMustGet("DYNAMO_CREDENTIALS_TABLE_NAME") }

func SignInWithApplePrivateKeyName() string { return osMustGet("SM_SIGNINWITHAPPLE_PRIVATEKEY_NAME") }

func CognitoDeveloperProviderName() string { return osMustGet("COGNITO_DEVELOPER_PROVIDER_NAME") }

func (e Environment) GetAwsConfig() aws.Config { return e.AwsConfig }

func NewEnv(ctx context.Context) (env Environment, err error) {

	return
}