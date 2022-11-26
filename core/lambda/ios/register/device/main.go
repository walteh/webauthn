package main

import (
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/env"
	"nugg-webauthn/core/pkg/errors"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/invocation"
	"nugg-webauthn/core/pkg/webauthn/handlers/devicecheck"

	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	Id      string
	Ctx     context.Context
	Dynamo  *dynamo.Client
	Config  config.Config
	logger  zerolog.Logger
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
		Dynamo:  dynamo.NewClient(cfg, "", env.DynamoCeremoniesTableName(), env.DynamoCredentialsTableName()),
		Config:  cfg,
		logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, input Input) (Output, error) {

	inv, ctx := invocation.NewInvocation(ctx, h, input)

	attestation := hex.HexToHash(input.Body)
	attestationKey := hex.HexToHash(input.Headers["x-nugg-hex-attestation-key"])
	clientDataJson := input.Headers["x-nugg-utf-client-data-json"]
	sessionId := hex.HexToHash(input.Headers["x-nugg-hex-session-id"])

	if attestation.IsZero() || clientDataJson == "" {
		return inv.Error(nil, 400, "missing required headers")
	}

	res, err := devicecheck.Attest(ctx, h.Dynamo, devicecheck.DeviceCheckAttestationInput{
		RawAttestationObject: attestation,
		UTF8ClientDataJSON:   clientDataJson,
		RawCredentialID:      attestationKey,
		RawSessionID:         sessionId,
	})
	if err != nil || !res.OK {
		if err == nil {
			err = errors.ErrVerification
		}
		return inv.Error(err, res.SuggestedStatusCode, "failed to attest devicecheck")
	}

	return inv.Success(204, map[string]string{}, "")

}
