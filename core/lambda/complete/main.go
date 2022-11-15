package main

import (
	"context"
	"encoding/json"
	"nugg-auth/core/pkg/applepublickey"
	"nugg-auth/core/pkg/cognito"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/safeid"
	"nugg-auth/core/pkg/secretsmanager"
	"nugg-auth/core/pkg/signinwithapple"
	"nugg-auth/core/pkg/webauthn/webauthn"

	"os"
	"time"

	"nugg-auth/core/pkg/webauthn/protocol"

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
	DynamoChallenge *dynamo.Client
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

	web, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "nugg.xyz",
		RPID:          "nugg.xyz",
		RPOrigin:      "https://auth.nugg.xyz",
		// AuthenticatorSelection: protocol.AuthenticatorSelection{
		// 	AuthenticatorAttachment: protocol.AuthenticatorAttachment("apple"),
		// 	UserVerification:        protocol.VerificationRequired,
		// 	ResidentKey:             protocol.ResidentKeyRequirementRequired,
		// 	RequireResidentKey:      protocol.ResidentKeyRequired(),
		// },
		AttestationPreference: protocol.PreferDirectAttestation,
	})

	if err != nil {
		return
	}

	abc := &Handler{
		Id:              safeid.Make().String(),
		Ctx:             ctx,
		DynamoUser:      dynamo.NewClient(cfg, env.DynamoUserTableName()),
		DynamoChallenge: dynamo.NewClient(cfg, env.DynamoChallengeTableName()),
		Cognito:         cognito.NewClient(cfg, env.AppleIdentityPoolId()),
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

	if attestation == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-attestation")
	}

	if clientdata == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-clientdata")
	}

	if credentialId == "" {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-credential-id")
	}

	session, user, wanu, err := h.DynamoChallenge.LoadChallenge(ctx, clientdata, "webauthn.create", "https://nugg.xyz")
	if err != nil {
		if err == dynamo.ErrNotFound {
			return inv.Error(err, 404, "challenge not found")
		}
		return inv.Error(err, 500, "failed to load challenge")
	}

	// attestationReader := strings.NewReader(attestation)

	parsedResponse, err := protocol.ParseCredentialCreation(clientdata, attestation, credentialId, "public-key")
	if err != nil {
		return inv.Error(err, 500, "failed to parse attestation")
	}

	credential, err := h.WebAuthn.CreateCredential(wanu, *session, parsedResponse)
	if err != nil {
		return inv.Error(err, 500, "failed to create credential")
	}

	wanu.AddCredential(credential)

	options, sessionData, err := h.WebAuthn.BeginLogin(wanu)
	if err != nil {
		return inv.Error(err, 500, "failed to begin login")
	}

	opts, err := json.Marshal(options)
	if err != nil {
		return inv.Error(err, 500, "Failed to marshal options")
	}

	err = h.DynamoUser.SaveNewUser(ctx, user, sessionData)
	if err != nil {
		return inv.Error(err, 500, "failed to save new user")
	}

	return inv.Success(200, map[string]string{
		"Content-Type": "application/json",
	}, string(opts))
}
