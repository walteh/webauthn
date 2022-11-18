package main

import (
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/safeid"
	"nugg-auth/core/pkg/webauthn/protocol"
	"nugg-auth/core/pkg/webauthn/webauthn"

	"context"
	"os"
	"time"

	"github.com/k0kubun/pp"
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
		Dynamo:   dynamo.NewClient(cfg, "", env.DynamoCeremonyTableName(), env.DynamoCeremonyTableName()),
		Config:   cfg,
		WebAuthn: web,
		Logger:   zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter:  0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	inv := h.NewInvocation(h.Logger)

	attestation := payload.Headers["x-nugg-webauthn-attestation"]
	clientdata := payload.Headers["x-nugg-webauthn-clientdata"]
	credentialId := payload.Headers["x-nugg-webauthn-credential-id"]

	if attestation == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-attestation")
	}

	if clientdata == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-clientdata")
	}

	if credentialId == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-credential-id")
	}

	parsedResponse, err := protocol.ParseCredentialCreation(clientdata, attestation, credentialId, "public-key")
	if err != nil {
		return inv.Error(err, 400, "failed to parse attestation")
	}

	getter := h.Dynamo.NewCeremonyGet(parsedResponse.Response.CollectedClientData.Challenge)

	res, err := h.Dynamo.TransactGet(ctx, types.TransactGetItem{Get: getter})
	if err != nil {
		return inv.Error(err, 500, "failed to get ceremony")
	}

	ceremony, err := h.Dynamo.FindCeremonyInGetResult(res)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return inv.Error(err, 404, "ceremony not found")
		}
		return inv.Error(err, 500, "failed to load ceremony")
	}

	credential, err := h.WebAuthn.CreateCredential("", *ceremony.SessionData, parsedResponse)
	if err != nil {
		return inv.Error(err, 500, "failed to create credential")
	}

	nuggid := safeid.Make()

	_, sessionData, err := h.WebAuthn.BeginLogin(nuggid.String(), []webauthn.Credential{*credential})
	if err != nil {
		return inv.Error(err, 500, "failed to begin login")
	}

	userput, err := h.Dynamo.NewUserPut(nuggid.String())
	if err != nil {
		return inv.Error(err, 500, "failed to create user put")
	}

	credput, err := h.Dynamo.NewApplePassKeyCredentialPut(nuggid.String(), sessionData.UserID, credential)
	if err != nil {
		return inv.Error(err, 500, "failed to create credential put")
	}

	cerput, err := h.Dynamo.NewCeremonyPut(sessionData)
	if err != nil {
		return inv.Error(err, 500, "failed to create ceremony put")
	}

	h.Dynamo.TransactWrite(ctx,
		types.TransactWriteItem{Put: userput},
		types.TransactWriteItem{Put: credput},
		types.TransactWriteItem{Put: cerput},
	)

	pp.Println(types.TransactWriteItem{Put: userput},
		types.TransactWriteItem{Put: credput},
		types.TransactWriteItem{Put: cerput})

	return inv.Success(204, map[string]string{"x-nugg-challenge": sessionData.Challenge.String()}, "")
}
