package main

import (
	"nugg-auth/core/pkg/cognito"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/safeid"
	"nugg-auth/core/pkg/webauthn/protocol"

	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
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
	Cognito cognito.Client
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
		Dynamo:  dynamo.NewClient(cfg, env.DynamoUsersTableName(), env.DynamoCeremoniesTableName(), env.DynamoCredentialsTableName()),
		Config:  cfg,
		Cognito: cognito.NewClient(cfg, env.AppleIdentityPoolId(), env.CognitoDeveloperProviderName()),
		Logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	inv := h.NewInvocation(ctx, h.Logger)

	attestation := payload.Headers["x-nugg-webauthn-creation"]

	if attestation == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-creation")
	}

	creation, err := protocol.ParseWebauthnCreation(attestation)
	if err != nil {
		return inv.Error(err, 400, "failed to parse webauthn creation")
	}

	parsedResponse, err := protocol.ParseCredentialCreationResponseHeader(creation, "public-key")
	if err != nil {
		return inv.Error(err, 400, "failed to parse attestation")
	}

	cerem := protocol.NewUnsafeGettableCeremony(parsedResponse.Response.CollectedClientData.Challenge)

	err = h.Dynamo.TransactGet(ctx, cerem)
	if err != nil {
		return inv.Error(err, 500, "failed to get ceremony")
	}

	cred, invalidErr := parsedResponse.Verify(protocol.Challenge(cerem.ChallengeID), cerem.SessionID, false, "nugg.xyz", "https://nugg.xyz")
	if invalidErr != nil {
		return inv.Error(invalidErr, 400, "invalid attestation")
	}

	nuggid := safeid.Make()

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

	userput, err := h.Dynamo.NewUserPut(nuggid.String())
	if err != nil {
		return inv.Error(err, 500, "failed to create user put")
	}

	credput, err := h.Dynamo.BuildPut(cred)
	if err != nil {
		return inv.Error(err, 500, "failed to create credential put")
	}

	ceremput, err := h.Dynamo.BuildPut(cerem)
	if err != nil {
		return inv.Error(err, 500, "failed to create ceremony put")
	}

	err = h.Dynamo.TransactWrite(ctx,
		types.TransactWriteItem{Put: userput},
		types.TransactWriteItem{Put: credput},
		types.TransactWriteItem{Put: ceremput},
	)
	if err != nil {
		return inv.Error(err, 500, "failed to write to dynamo")
	}

	result := <-chaner
	stale = true

	if result == nil {
		return inv.Error(chanerr, 500, "failed to get dev creds")
	}

	return inv.Success(204, map[string]string{"x-nugg-access-token": *result.Token}, "")
}
