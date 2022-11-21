package main

import (
	"context"
	"log"
	"nugg-auth/core/pkg/cognito"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/hex"
	"nugg-auth/core/pkg/webauthn/protocol"

	"os"
	"time"

	"github.com/k0kubun/pp"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	Id      string
	Dynamo  *dynamo.Client
	Config  config.Config
	Logger  zerolog.Logger
	Cognito cognito.Client
	Ctx     context.Context
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

	h.cancel()

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

	h.cancel()

	return output, nil
}

func main() {

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	abc := &Handler{
		Id:      ksuid.New().String(),
		Dynamo:  dynamo.NewClient(cfg, env.DynamoUsersTableName(), env.DynamoCeremoniesTableName(), env.DynamoCredentialsTableName()),
		Config:  cfg,
		Cognito: cognito.NewClient(cfg, env.AppleIdentityPoolId(), env.CognitoDeveloperProviderName()),
		Logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, input Input) (Output, error) {

	inv := h.NewInvocation(ctx, h.Logger)

	assertion := hex.HexToHash(input.Headers["x-nugg-hex-assertion"])

	if assertion.IsZero() {
		return inv.Error(nil, 400, "invalid assertion")
	}

	args, err := protocol.ParseCredentialAssertionResponsePayload(assertion)
	if err != nil {
		return inv.Error(nil, 400, "unable to decode headers")
	}

	r1 := protocol.DecodeCredentialAssertionResponse(args)

	parsedResponse, err := protocol.ParseCredentialAssertionResponse(*r1)
	if err != nil {
		return inv.Error(err, 400, "failed to parse attestation")
	}

	cred := protocol.NewUnsafeGettableCredential(parsedResponse.RawID)

	cerem := protocol.NewUnsafeGettableCeremony(parsedResponse.Response.CollectedClientData.Challenge)

	err = h.Dynamo.TransactGet(inv.Ctx, cred, cerem)
	if err != nil {
		return inv.Error(err, 500, "failed to send transact get")
	}

	pp.Println(cred, cerem)

	if cred.RawID.Hex() != cerem.CredentialID.Hex() {
		log.Println(cred.RawID.Hex(), cerem.CredentialID.Hex())
		return inv.Error(nil, 400, "invalid credential id")
	}

	chaner := make(chan *cognitoidentity.GetOpenIdTokenForDeveloperIdentityOutput, 1)
	defer close(chaner)
	stale := false
	var chanerr error

	go func() {
		go func() {
			<-inv.Ctx.Done()
			if !stale {
				chaner <- nil
			}
			stale = true

		}()

		z, err := h.Cognito.GetDevCreds(ctx, cerem.CredentialID)
		if !stale {
			if err != nil {
				chanerr = err
				chaner <- nil
			} else {
				chaner <- z
			}
		}
	}()

	// Handle steps 4 through 16
	validError := parsedResponse.Verify(cerem.ChallengeID, env.RPID(), env.RPOrigin(), "none", false, cred.PublicKey, nil)
	if validError != nil {
		return inv.Error(validError, 400, "failed to verify assertion")
	}

	credentialUpdate, err := cred.Update(h.Dynamo.MustCredentialTableName(), parsedResponse.Response.AuthenticatorData.Counter)
	if err != nil {
		return inv.Error(err, 500, "failed to create apple pass key")
	}

	err = h.Dynamo.TransactWrite(ctx, *credentialUpdate)
	if err != nil {
		return inv.Error(err, 500, "failed to update apple pass key")
	}

	result := <-chaner
	stale = true

	if result == nil {
		return inv.Error(chanerr, 500, "failed to get dev creds")
	}

	return inv.Success(204, map[string]string{
		"x-nugg-access-token": *result.Token,
	}, "")
}
