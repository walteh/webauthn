package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/webauthn/protocol"
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

func (h *Handler) Invoke(ctx context.Context, input Input) (Output, error) {

	inv := h.NewInvocation(h.Logger)

	attestation := input.Headers["x-nugg-devicecheck-attestation"]
	clientData := input.Headers["x-nugg-devicecheck-clientdatajson"]
	payload := input.Headers["x-nugg-devicecheck-payload"]

	if attestation == "" || clientData == "" {
		return inv.Error(nil, 400, "missing required headers")
	}

	i, err := protocol.FormatAttestationInput(clientData, attestation)
	if err != nil {
		return inv.Error(err, 400, "invalid attestation")
	}

	p, err := i.Parse()
	if err != nil {
		return inv.Error(err, 400, err.Error())
	}

	b, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return inv.Error(err, 400, "invalid payload")
	}
	d := sha256.Sum256(b)

	// relyiing part for apple appattestation
	pk, err := p.AttestationObject.Verify("4497QJSAD3.xyz.nugg.app", d[:], false, false)
	if err != nil {
		return inv.Error(err, 400, err.Error())
	}

	putter, err := dynamo.MakePut(h.Dynamo.MustCredentialTableName(), pk)
	if err != nil {
		return inv.Error(err, 500, err.Error())
	}

	err = h.Dynamo.TransactWrite(ctx, types.TransactWriteItem{Put: putter})
	if err != nil {
		return inv.Error(err, 500, err.Error())
	}

	// webauthn.MakeNewCredential(&protocol.ParsedCredentialCreationData{
	// 	ParsedPublicKeyCredential: protocol.ParsedPublicKeyCredential{
	// 		RawID: rec.RawID,
	// 		Type:  rec.Type,
	// 		ParsedCredential: protocol.ParsedCredential{
	// 			ID: rec.ID,
	// 			Type: "",
	// 		},
	// 		ClientExtensionResults: ,
	// 	},
	// })

	// getter := h.Dynamo.NewCeremonyGet(string(clientData))

	// res, err := h.Dynamo.TransactGet(ctx, types.TransactGetItem{Get: getter})
	// if err != nil {
	// 	return inv.Error(err, 500, "failed to get ceremony")
	// }

	// c, err := h.Dynamo.FindDeviceCheckCeremonyInGetResult(res)
	// if err != nil {
	// 	return inv.Error(err, 500, "failed to find ceremony")
	// }

	// if c.Id != clientData {
	// 	return inv.Error(nil, 400, "ceremony mismatch")
	// }

	return inv.Success(204, map[string]string{
		"Content-Length": "0",
	}, "")

}
