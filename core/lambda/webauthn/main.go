package main

import (
	"context"
	"nugg-auth/core/pkg/applepublickey"
	"nugg-auth/core/pkg/cognito"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/random"
	"nugg-auth/core/pkg/secretsmanager"
	"nugg-auth/core/pkg/signinwithapple"
	"nugg-auth/core/pkg/webauthn"

	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	Id              string
	Ctx             context.Context
	DynamoUser      *dynamo.Client
	DynamoCeremony  *dynamo.Client
	Cognito         *cognito.Client
	SignInWithApple *signinwithapple.Client
	ApplePublicKey  *applepublickey.Client
	SecretsManager  *secretsmanager.Client
	Config          config.Config
	Logger          zerolog.Logger
	counter         int
	WebAuthn        *webauthn.WebAuthn
}

func init() {
	zerolog.TimeFieldFormat = time.StampMicro
}

func main() {

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}

	web, err := webauthn.NewConfig()

	if err != nil {
		return
	}

	abc := &Handler{
		Id:             random.KSUID().String(),
		Ctx:            ctx,
		DynamoUser:     dynamo.NewClient(cfg, env.DynamoUserTableName()),
		DynamoCeremony: dynamo.NewClient(cfg, env.DynamoCeremonyTableName()),

		Cognito:         cognito.NewClient(cfg, env.AppleIdentityPoolId()),
		SignInWithApple: signinwithapple.NewClient(env.AppleTokenEndpoint(), env.AppleTeamID(), env.AppleServiceName(), env.SignInWithApplePrivateKeyID()),
		ApplePublicKey:  applepublickey.NewClient(env.ApplePublicKeyEndpoint()),
		SecretsManager:  secretsmanager.NewClient(ctx, cfg, env.SignInWithApplePrivateKeyName()),
		WebAuthn:        web,
		Config:          cfg,
		Logger:          zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),

		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

type Invocation struct {
	zerolog.Logger
	Start time.Time
}

func (h *Handler) NewInvocation() *Invocation {
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

	h.Logger.Info().
		Int("status_code", code).
		Str("body", output.Body).
		Dict("headers", r).
		CallerSkipFrame(1).
		TimeDiff("duration", time.Now(), h.Start).
		Msg(message)

	return output, nil
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	inv := h.NewInvocation()

	attestation := payload.Headers["x-nugg-webauthn-attestation"]
	clientdata := payload.Headers["x-nugg-webauthn-clientdata"]
	assertion := payload.Headers["x-nugg-webauthn-assertion"]
	credentialId := payload.Headers["x-nugg-webauthn-credential-id"]
	userId := payload.Headers["x-nugg-webauthn-user-id"]
	username := payload.Headers["x-nugg-webauthn-username"]

	if attestation == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-attestation")
	}

	if clientdata == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-clientdata")
	}

	if assertion == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-assertion")
	}

	if credentialId == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-credential-id")
	}

	if userId == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-user-id")
	}

	if username == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-username")
	}

	user := webauthn.NewUser([]byte(userId), username)

	options, sessionData, err := h.WebAuthn.BeginRegistration(user)
	if err != nil {
		return inv.Error(err, 500, "failed to begin registration")
	}

	err = h.DynamoCeremony.StartWebAuthnCeremony(ctx, userId, sessionData)
	if err != nil {
		return inv.Error(err, 500, "failed to start webauthn ceremony")
	}

	return inv.Success(204, map[string]string{
		"x-nugg-webauthn-challange": string(options.Response.Challenge),
		"x-nugg-webauthn-user-id":   string(options.Response.User.ID),
		"x-nugg-webauthn-username":  options.Response.User.Name,
	}, "")
}
