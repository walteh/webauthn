package main

import (
	"context"
	"nugg-auth/core/pkg/applepublickey"
	"nugg-auth/core/pkg/cognito"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/secretsmanager"
	"nugg-auth/core/pkg/signinwithapple"
	"nugg-auth/core/pkg/webauthn/webauthn"

	"os"
	"time"

	"nugg-auth/core/pkg/webauthn/protocol"

	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	Id              string
	Ctx             context.Context
	Dynamo          *dynamo.Client
	Cognito         cognito.Client
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
		Id:              ksuid.New().String(),
		Ctx:             ctx,
		Dynamo:          dynamo.NewClient(cfg, env.DynamoUsersTableName(), env.DynamoCeremoniesTableName(), ""),
		Cognito:         cognito.NewClient(cfg, env.AppleIdentityPoolId(), env.CognitoDeveloperProviderName()),
		SignInWithApple: signinwithapple.NewClient(env.AppleTokenEndpoint(), env.AppleTeamID(), env.AppleServiceName(), env.SignInWithApplePrivateKeyID()),
		ApplePublicKey:  applepublickey.NewClient(env.ApplePublicKeyEndpoint()),
		SecretsManager:  secretsmanager.NewClient(ctx, cfg, env.SignInWithApplePrivateKeyName()),
		WebAuthn:        web,
		Config:          cfg,
		Logger:          zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter:         0,
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
		Headers: map[string]string{
			"Content-Length": "0",
		},
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
	credentialId := payload.Headers["x-nugg-webauthn-credential-id"]
	appleId := payload.Headers["x-nugg-webauthn-apple-id"]

	if attestation == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-attestation")
	}

	if clientdata == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-clientdata")
	}

	if credentialId == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-credential-id")
	}

	if appleId == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-apple-id")
	}

	// parsedResponse, err := protocol.ParseCredentialCreation(clientdata, attestation, credentialId, "public-key")
	// if err != nil {
	// 	return inv.Error(err, 500, "failed to parse attestation")
	// }

	// if parsedResponse.Response.AttestationObject.Format != "apple" {
	// 	return inv.Error(nil, 400, "invalid format")
	// }

	// ceremony, err := h.Dynamo.LoadCeremony(ctx, parsedResponse.Response.CollectedClientData.Challenge)
	// if err != nil {
	// 	if err == dynamo.ErrNotFound {
	// 		return inv.Error(err, 404, "ceremony not found")
	// 	}
	// 	return inv.Error(err, 500, "failed to load ceremony")
	// }

	// credential, err := h.WebAuthn.CreateCredential(appleId, *ceremony.SessionData, parsedResponse)
	// if err != nil {
	// 	return inv.Error(err, 500, "failed to create credential")
	// }

	// user, err := h.Dynamo.LoadUser(ctx, ceremony.UserId)
	// if err != nil {
	// 	if err == dynamo.ErrNotFound {
	// 		return inv.Error(err, 404, "user not found")
	// 	}
	// 	return inv.Error(err, 500, "failed to load user")
	// }

	// user.AppleAuthData.AddAppleWebAuthnCredentials(credential)

	// options, sessionData, err := h.WebAuthn.BeginLogin(user.CreateAppleWebAuthnUser())
	// if err != nil {
	// 	return inv.Error(err, 500, "failed to begin login")
	// }

	// opts, err := json.Marshal(options)
	// if err != nil {
	// 	return inv.Error(err, 500, "Failed to marshal options")
	// }

	// cer := dynamo.NewCeremony(user.Id, sessionData)

	// err = h.Dynamo.SaveFirstUserLogin(ctx, user, cer)
	// if err != nil {
	// 	return inv.Error(err, 500, "failed to save new user")
	// }

	return inv.Success(200, map[string]string{"Content-Type": "application/json"}, string(""))
}
