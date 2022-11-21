package main

import (
	"context"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/hex"
	"nugg-auth/core/pkg/webauthn/protocol"

	"os"
	"time"

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
	Logger  zerolog.Logger
	counter int
}

func init() {
	zerolog.TimeFieldFormat = time.StampMicro
}

type Invocation struct {
	zerolog.Logger
	Start time.Time
}

func (h *Handler) NewInvocation(logger zerolog.Logger) *Invocation {
	h.counter++
	return &Invocation{
		Logger: h.Logger.With().Int("counter", h.counter).Str("handler", h.Id).Logger(),
		Start:  time.Now(),
	}
}

func (h *Invocation) Error(err error, code int, message string) (Output, error) {
	h.Logger.Error().Err(err).
		Int("status_code", code).
		Str("body", "").
		CallerSkipFrame(1).
		TimeDiff("duration", time.Now(), h.Start).
		Msg(message)

	return Output{
		StatusCode: code,
	}, nil
}

func (h *Invocation) Success(code int, headers map[string]string, message string) (Output, error) {

	output := Output{
		StatusCode: code,
		Headers:    headers,
	}

	if message != "" && code != 204 {
		output.Body = message
	}

	r := zerolog.Dict()
	for k, v := range output.Headers {
		r = r.Str(k, v)
	}

	if message == "" {
		message = "empty"
	}

	if code == 204 && headers["Content-Length"] == "" {
		output.Headers["Content-Length"] = "0"
	}

	h.Logger.Info().
		Int("status_code", code).
		Str("body", output.Body).
		Dict("headers", r).
		CallerSkipFrame(1).
		TimeDiff("duration", time.Now(), h.Start).
		Msg(message)

	return output, nil
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
		Logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	inv := h.NewInvocation(h.Logger)

	sessionId := hex.HexToHash(payload.Headers["x-nugg-hex-session-id"])
	credentialId := hex.HexToHash(payload.Headers["x-nugg-hex-credential-id"])

	if sessionId.IsZero() {
		return inv.Error(nil, 400, "missing x-nugg-hex-sessionId header")
	}

	cha := protocol.NewCeremony(credentialId, sessionId, protocol.CreateCeremony)

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
