package main

import (
	"context"

	"github.com/walteh/webauthn/pkg/dynamo"
	"github.com/walteh/webauthn/pkg/env"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/invocation"
	"github.com/walteh/webauthn/pkg/webauthn/types"

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
		Dynamo:  dynamo.NewClient(cfg, "", env.DynamoCeremoniesTableName(), ""),
		Config:  cfg,
		logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	inv, ctx := h.NewInvocation(ctx, payload)

	sessionId := hex.HexToHash(payload.Headers["x-nugg-hex-session-id"])
	ceremonyType := payload.Headers["x-nugg-utf-ceremony-type"]
	credentialId := hex.HexToHash(payload.Headers["x-nugg-hex-credential-id"])

	if sessionId.IsZero() {
		return inv.Error(nil, 400, "missing x-nugg-hex-sessionId header")
	}

	if ceremonyType == "" {
		ceremonyType = string(types.AssertCeremony)
	}

	switch ceremonyType {
	case string(types.AssertCeremony):
	case string(types.CreateCeremony):
		break
	default:
		return inv.Error(nil, 400, "invalid x-nugg-utf-ceremony-type header")
	}

	cha := types.NewCeremony(credentialId, sessionId, types.CeremonyType(ceremonyType))

	cer, err := dynamo.MakePut(h.Dynamo.MustCeremonyTableName(), cha)
	if err != nil {
		return inv.Error(err, 500, "failed to create ceremony")
	}

	err = h.Dynamo.TransactWrite(ctx, *cer)
	if err != nil {
		return inv.Error(err, 500, "Failed to save ceremony")
	}

	return inv.Success(Output{
		StatusCode: 204,
		Headers: map[string]string{
			"x-nugg-hex-challenge": cha.ChallengeID.Hex(),
		},
	})
}
