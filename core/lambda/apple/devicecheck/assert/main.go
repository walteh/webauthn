package main

import (
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/env"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/webauthn/protocol"

	"context"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	Logger  zerolog.Logger
	counter int
}

func init() {
	zerolog.TimeFieldFormat = time.StampMicro
}

type Invocation struct {
	zerolog.Logger
	Start  time.Time
	Ctx    context.Context
	cancel context.CancelFunc
}

func (h *Handler) NewInvocation(ctx context.Context, logger zerolog.Logger) *Invocation {
	ctx, cnl := context.WithCancel(ctx)

	h.counter++
	return &Invocation{
		cancel: cnl,
		Ctx:    ctx,
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
		Dynamo:  dynamo.NewClient(cfg, "", env.DynamoCeremoniesTableName(), env.DynamoCredentialsTableName()),
		Config:  cfg,
		Logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, input Input) (Output, error) {

	inv := h.NewInvocation(ctx, h.Logger)

	attestation := hex.HexToHash(input.Headers["x-nugg-hex-attestation"])
	clientData := input.Headers["x-nugg-utf-client-data-json"]
	payload := hex.HexToHash(input.Headers["x-nugg-hex-payload"])

	if attestation.IsZero() || clientData == "" {
		return inv.Error(nil, 400, "missing required headers")
	}

	p, err := protocol.FormatAttestationInput(clientData, attestation).Parse()
	if err != nil {
		return inv.Error(err, 400, err.Error())
	}

	cer := protocol.NewUnsafeGettableCeremony(p.CollectedClientData.Challenge)

	err = h.Dynamo.TransactGet(ctx, cer)
	if err != nil {
		return inv.Error(err, 400, err.Error())
	}

	if !cer.ChallengeID.Equals(p.CollectedClientData.Challenge) {
		return inv.Error(nil, 400, "invalid credential id")
	}

	pk, err := p.AttestationObject.Verify("4497QJSAD3.xyz.nugg.app", payload.Sha256(), false, false)
	if err != nil {
		return inv.Error(err, 400, err.Error())
	}

	putter, err := dynamo.MakePut(h.Dynamo.MustCredentialTableName(), pk)
	if err != nil {
		return inv.Error(err, 500, err.Error())
	}

	err = h.Dynamo.TransactWrite(inv.Ctx, types.TransactWriteItem{Put: putter})
	if err != nil {
		return inv.Error(err, 500, err.Error())
	}

	return inv.Success(204, map[string]string{}, "")

}
