package main

import (
	"context"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/webauthn/protocol"
	"nugg-auth/core/pkg/webauthn/webauthn"

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
	Id       string
	Ctx      context.Context
	Dynamo   *dynamo.Client
	Config   config.Config
	Logger   zerolog.Logger
	WebAuthn *webauthn.WebAuthn
	counter  int
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

	web, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "nugg.xyz",
		RPID:          "nugg.xyz",
		RPOrigin:      "https://nugg.xyz",
		// passkeys do not support attestation as they can move between devices
		// https://developer.apple.com/forums/thread/713195
		AttestationPreference: protocol.PreferNoAttestation,
	})
	if err != nil {
		return
	}

	abc := &Handler{
		Id:       ksuid.New().String(),
		Ctx:      ctx,
		Dynamo:   dynamo.NewClient(cfg, "", env.DynamoCeremonyTableName(), ""),
		Config:   cfg,
		WebAuthn: web,
		Logger:   zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter:  0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	inv := h.NewInvocation(h.Logger)

	signature := payload.Headers["x-nugg-webauthn-signature"]
	clientdata := payload.Headers["x-nugg-webauthn-clientdata"]
	credentialId := payload.Headers["x-nugg-webauthn-credential-id"]
	userId := payload.Headers["x-nugg-webauthn-user-id"]

	if signature == "" || clientdata == "" || credentialId == "" {
		return inv.Error(nil, 400, "missing headers")
	}

	parsedResponse, err := protocol.ParseCredentialRequest(credentialId, signature, clientdata, userId, "public-key")
	if err != nil {
		return inv.Error(err, 400, "failed to parse attestation")
	}

	res, err := h.Dynamo.TransactGet(ctx,
		types.TransactGetItem{Get: h.Dynamo.NewCeremonyGet(parsedResponse.Response.CollectedClientData.Challenge)},
		types.TransactGetItem{Get: h.Dynamo.NewCredentialGet(parsedResponse.ParsedCredential.ID)},
	)
	if err != nil {
		return inv.Error(err, 500, "failed to get ceremony")
	}

	cer, err := h.Dynamo.FindCeremonyInGetResult(res)
	if err != nil {
		return inv.Error(err, 500, "failed to find ceremony")
	}

	userid, creds, err := h.Dynamo.FindApplePassKeyInGetResult(res)
	if err != nil {
		return inv.Error(err, 500, "failed to find ceremony")
	}

	dacred, err := h.WebAuthn.ValidateLogin(userid, []webauthn.Credential{*creds}, *cer.SessionData, parsedResponse)
	if err != nil {
		return inv.Error(err, 500, "failed to begin registration")
	}

	apl, err := h.Dynamo.NewApplePassKeyCredentialUpdate(userid, dacred)
	if err != nil {
		return inv.Error(err, 500, "failed to create apple pass key")
	}

	err = h.Dynamo.TransactWrite(ctx, types.TransactWriteItem{Update: apl})
	if err != nil {
		return inv.Error(err, 500, "failed to update apple pass key")
	}

	return inv.Success(204, map[string]string{}, "")
}
