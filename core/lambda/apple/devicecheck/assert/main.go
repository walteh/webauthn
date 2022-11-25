package main

import (
	"log"
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/env"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/webauthn/assertion"
	"nugg-webauthn/core/pkg/webauthn/clientdata"
	"nugg-webauthn/core/pkg/webauthn/extensions"
	"nugg-webauthn/core/pkg/webauthn/types"

	"context"
	"os"
	"time"

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

	assert := hex.HexToHash(input.Headers["x-nugg-hex-request-assertion"])

	body := hex.HexToHash(input.Body)

	if assert.IsZero() || body.IsZero() {
		return inv.Error(nil, 400, "missing required headers")
	}

	parsed, err := assertion.ParseAssertionResponse(assert)
	if err != nil {
		return inv.Error(err, 400, "failed to parse assertion")
	}

	payload := types.AssertionInput{
		UserID:               parsed.SessionID,
		CredentialID:         parsed.CredentialID,
		RawClientDataJSON:    parsed.UTF8ClientDataJSON,
		RawAuthenticatorData: nil,
		Signature:            parsed.AssertionObject,
		Type:                 types.NotFidoAttestationType,
	}

	cd, err := clientdata.ParseClientData(parsed.UTF8ClientDataJSON)
	if err != nil {
		return inv.Error(err, 400, "failed to parse client data")
	}

	cred := types.NewUnsafeGettableCredential(parsed.CredentialID)
	cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	err = h.Dynamo.TransactGet(ctx, cred, cerem)
	if err != nil {
		return inv.Error(err, 500, "failed to send transact get")
	}

	if cred.RawID.Hex() != cerem.CredentialID.Hex() {
		log.Println(cred.RawID.Hex(), cerem.CredentialID.Hex())
		return inv.Error(nil, 400, "invalid credential id")
	}

	if !cerem.ChallengeID.Equals(cd.Challenge) {
		return inv.Error(nil, 400, "invalid challenge id")
	}

	// Handle steps 4 through 16
	validError := assertion.VerifyAssertionInput(types.VerifyAssertionInputArgs{
		Input:               payload,
		StoredChallenge:     cerem.ChallengeID,
		RelyingPartyID:      "4497QJSAD3.xyz.nugg.app",
		RelyingPartyOrigin:  env.RPOrigin(),
		AttestationType:     "none",
		VerifyUser:          false,
		CredentialPublicKey: cred.PublicKey,
		Extensions:          extensions.ClientInputs{},
	})
	if validError != nil {
		return inv.Error(validError, 400, "failed to verify assertion")
	}

	credentialUpdate, err := cred.UpdateIncreasingCounter(h.Dynamo.MustCredentialTableName())
	if err != nil {
		return inv.Error(err, 500, "failed to create apple pass key")
	}

	err = h.Dynamo.TransactWrite(inv.Ctx, *credentialUpdate)
	if err != nil {
		return inv.Error(err, 500, err.Error())
	}

	return inv.Success(204, map[string]string{}, "")

}
