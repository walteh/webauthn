package main

import (
	"git.nugg.xyz/go-sdk/invocation"
	"git.nugg.xyz/webauthn/pkg/cognito"
	"git.nugg.xyz/webauthn/pkg/dynamo"
	"git.nugg.xyz/webauthn/pkg/env"
	"git.nugg.xyz/webauthn/pkg/hex"
	"git.nugg.xyz/webauthn/pkg/webauthn/handlers/devicecheck"
	"git.nugg.xyz/webauthn/pkg/webauthn/handlers/passkey"

	"context"
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
	*invocation.Handler[Input, Output]
	Id      string
	Ctx     context.Context
	Dynamo  *dynamo.Client
	Config  config.Config
	logger  zerolog.Logger
	Cognito cognito.Client
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

	abc := &Handler{
		Id:      ksuid.New().String(),
		Ctx:     ctx,
		Dynamo:  dynamo.NewClient(cfg, env.DynamoUsersTableName(), env.DynamoCeremoniesTableName(), env.DynamoCredentialsTableName()),
		Config:  cfg,
		Cognito: cognito.NewClient(cfg, env.AppleIdentityPoolId(), env.CognitoDeveloperProviderName()),
		logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	inv, ctx := h.NewInvocation(ctx, payload)

	attestation := hex.HexToHash(payload.Headers["x-nugg-hex-attestation-object"])
	dc := hex.HexToHash(payload.Headers["x-nugg-hex-request-assertion"])
	clientData := payload.Headers["x-nugg-utf-client-data-json"]
	credentialID := hex.HexToHash(payload.Headers["x-nugg-hex-credential-id"])

	if attestation.IsZero() || clientData == "" || credentialID.IsZero() || dc.IsZero() {
		return inv.Error(nil, 400, "missing required headers")
	}

	if res, err := devicecheck.Assert(ctx, h.Dynamo, devicecheck.DeviceCheckAssertionInput{
		RawAssertionObject:   dc,
		ClientDataToValidate: credentialID,
	}); err != nil || !res.OK {
		return inv.Error(err, res.SuggestedStatusCode, "devicecheck assertion failed")
	}

	out, err := passkey.Attest(ctx, h.Dynamo, h.Cognito, passkey.PasskeyAttestationInput{
		RawAttestationObject: attestation,
		UTF8ClientDataJSON:   clientData,
		RawCredentialID:      credentialID,
	})
	if err != nil {
		return inv.Error(err, out.SuggestedStatusCode, "passkey attestation failed")
	}

	return inv.Success(Output{
		Headers: map[string]string{
			"x-nugg-utf-access-token": out.AccessToken,
		},
		StatusCode: 204,
	})
}
