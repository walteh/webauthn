package main

import (
	"context"

	"github.com/nuggxyz/golang/invocation"
	"github.com/nuggxyz/webauthn/pkg/cognito"
	"github.com/nuggxyz/webauthn/pkg/dynamo"
	"github.com/nuggxyz/webauthn/pkg/env"
	"github.com/nuggxyz/webauthn/pkg/hex"
	"github.com/nuggxyz/webauthn/pkg/webauthn/handlers/passkey"

	"os"

	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	Id      string
	Dynamo  *dynamo.Client
	Config  config.Config
	logger  zerolog.Logger
	Cognito cognito.Client
	Ctx     context.Context
	counter int
}

func (h Handler) ID() string {
	return h.Id
}

func (h *Handler) IncrementCounter() int {
	h.counter += 1
	return h.counter
}

func (h Handler) Logger() zerolog.Logger {
	return h.logger
}
func main() {

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	abc := &Handler{
		Id:      ksuid.New().String(),
		Dynamo:  dynamo.NewClient(cfg, env.DynamoUsersTableName(), env.DynamoCeremoniesTableName(), env.DynamoCredentialsTableName()),
		Config:  cfg,
		Cognito: cognito.NewClient(cfg, env.AppleIdentityPoolId(), env.CognitoDeveloperProviderName()),
		logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, input Input) (Output, error) {

	inv, ctx := invocation.NewInvocation(ctx, h, input)

	authenticatorData := hex.HexToHash(input.Headers["x-nugg-hex-authenticator-data"])
	credentialId := hex.HexToHash(input.Headers["x-nugg-hex-credential-id"])
	signature := hex.HexToHash(input.Headers["x-nugg-hex-signature"])
	userId := hex.HexToHash(input.Headers["x-nugg-hex-user-id"])

	clientDataJson := input.Headers["x-nugg-utf-client-data-json"]
	credentialType := input.Headers["x-nugg-utf-credential-type"]

	// make sure all the above values exist one by one in the headers
	if len(authenticatorData) == 0 || len(credentialId) == 0 || len(signature) == 0 || len(userId) == 0 || len(clientDataJson) == 0 || len(credentialType) == 0 {
		return inv.Error(nil, 400, "missing required headers")
	}

	abc := passkey.PasskeyAssertionInput{
		RawAuthenticatorData: authenticatorData,
		CredentialID:         credentialId,
		RawSignature:         signature,
		UTF8ClientDataJSON:   clientDataJson,
		SessionID:            userId,
	}

	res, err := passkey.Assert(ctx, h.Dynamo, h.Cognito, abc)
	if err != nil {
		return inv.Error(err, res.SuggestedStatusCode, "failed to assert passkey")
	}

	return inv.Success(204, map[string]string{"x-nugg-utf-access-token": res.AccessToken}, "")
}
