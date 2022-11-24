package main

import (
	"context"
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/env"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/invocation"
	protocol "nugg-webauthn/core/pkg/webauthn"

	"os"

	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
		Dynamo:  dynamo.NewClient(cfg, "", env.DynamoCeremoniesTableName(), ""),
		Config:  cfg,
		logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	inv, ctx := invocation.NewInvocation(ctx, h, payload)

	sessionId := hex.HexToHash(payload.Headers["x-nugg-hex-session-id"])
	ceremonyType := payload.Headers["x-nugg-utf-ceremony-type"]
	credentialId := hex.HexToHash(payload.Headers["x-nugg-hex-credential-id"])

	if sessionId.IsZero() {
		return inv.Error(nil, 400, "missing x-nugg-hex-sessionId header")
	}

	if ceremonyType == "" {
		ceremonyType = string(protocol.AssertCeremony)
	}

	switch ceremonyType {
	case string(protocol.AssertCeremony):
	case string(protocol.CreateCeremony):
		break
	default:
		return inv.Error(nil, 400, "invalid x-nugg-utf-ceremony-type header")
	}

	cha := protocol.NewCeremony(credentialId, sessionId, protocol.CeremonyType(ceremonyType))

	cer, err := dynamo.MakePut(h.Dynamo.MustCeremonyTableName(), cha)
	if err != nil {
		return inv.Error(err, 500, "failed to create ceremony")
	}

	err = h.Dynamo.TransactWrite(ctx, types.TransactWriteItem{Put: cer})
	if err != nil {
		return inv.Error(err, 500, "Failed to save ceremony")
	}

	return inv.Success(204, map[string]string{"x-nugg-hex-challenge": cha.ChallengeID.Hex()}, "")
}
