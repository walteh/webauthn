package main

import (
	"github.com/walteh/webauthn/app/devicecheck_assert"
	"github.com/walteh/webauthn/app/passkey_attest"
	"github.com/walteh/webauthn/pkg/accesstoken"
	"github.com/walteh/webauthn/pkg/accesstoken/cognito"

	"github.com/walteh/webauthn/pkg/aws/invocation"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/storage"

	"context"
	"os"

	"github.com/rs/xid"
	"github.com/rs/zerolog"

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
	Dynamo  storage.Provider
	Config  config.Config
	logger  zerolog.Logger
	Cognito accesstoken.Provider
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
		Id:      xid.New().String(),
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

	if res, err := devicecheck_assert.Assert(ctx, h.Dynamo, devicecheck_assert.DeviceCheckAssertionInput{
		RawAssertionObject:   dc,
		ClientDataToValidate: credentialID,
	}); err != nil || !res.OK {
		return inv.Error(err, res.SuggestedStatusCode, "devicecheck assertion failed")
	}

	out, err := passkey_attest.Attest(ctx, h.Dynamo, h.Cognito, passkey_attest.PasskeyAttestationInput{
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
