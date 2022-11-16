package main

import (
	"context"
	"encoding/json"
	"nugg-auth/core/pkg/applepublickey"
	"nugg-auth/core/pkg/cognito"
	"nugg-auth/core/pkg/cwebauthn"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/safeid"
	"nugg-auth/core/pkg/secretsmanager"
	"nugg-auth/core/pkg/signinwithapple"
	"nugg-auth/core/pkg/webauthn/protocol"
	"nugg-auth/core/pkg/webauthn/webauthn"

	"os"
	"time"

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
	Cognito         *cognito.Client
	SignInWithApple *signinwithapple.Client
	ApplePublicKey  *applepublickey.Client
	SecretsManager  *secretsmanager.Client
	Config          config.Config
	Logger          zerolog.Logger
	WebAuthn        *webauthn.WebAuthn
	counter         int
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
		RPDisplayName:         "nugg.xyz",
		RPID:                  "nugg.xyz",
		RPOrigin:              "https://nugg.xyz",
		AttestationPreference: protocol.PreferDirectAttestation,
	})
	if err != nil {
		return
	}

	abc := &Handler{
		Id:              ksuid.New().String(),
		Ctx:             ctx,
		Dynamo:          dynamo.NewClient(cfg, env.DynamoChallengeTableName()),
		Cognito:         cognito.NewClient(cfg, env.AppleIdentityPoolId()),
		SignInWithApple: signinwithapple.NewClient(env.AppleTokenEndpoint(), env.AppleTeamID(), env.AppleServiceName(), env.SignInWithApplePrivateKeyID()),
		ApplePublicKey:  applepublickey.NewClient(env.ApplePublicKeyEndpoint()),
		SecretsManager:  secretsmanager.NewClient(ctx, cfg, env.SignInWithApplePrivateKeyName()),
		Config:          cfg,
		WebAuthn:        web,
		Logger:          zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter:         0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	out, err := json.Marshal(payload)
	if err != nil {
		return Output{}, err
	}

	h.Logger.Info().Msg(string(out))

	inv := h.NewInvocation(h.Logger)

	h1 := payload.Headers["x-nugg-signinwithapple-identity-token"]
	c1 := payload.Headers["x-nugg-signinwithapple-registration-code"]
	u1 := payload.Headers["x-nugg-signinwithapple-username"]

	if h1 == "" {
		return inv.Error(nil, 400, "Missing x-nugg-signinwithapple-identity-token")
	}

	if c1 == "" {
		return inv.Error(nil, 400, "Missing x-nugg-signinwithapple-registration-code")
	}

	if u1 == "" {
		return inv.Error(nil, 400, "Missing x-nugg-signinwithapple-username")
	}

	publickey, err := h.ApplePublicKey.Refresh(ctx)
	if err != nil {
		return inv.Error(err, 502, "Failed to refresh public key")
	}

	tkn, err := publickey.ParseToken(h1)
	if err != nil {
		return inv.Error(err, 401, "Failed to parse token")
	}

	if !tkn.Valid {
		return inv.Error(err, 401, "Invalid token")
	}

	sub, err := tkn.GetUniqueID()
	if err != nil {
		return inv.Error(err, 400, "Failed to get sub")
	}

	privateKey, err := h.SecretsManager.Refresh(ctx)
	if err != nil {
		return inv.Error(err, 502, "Failed to refresh private key")
	}

	creds, err := h.Cognito.GetIdentityId(h.Ctx, h1)
	if err != nil {
		return inv.Error(err, 502, "Failed to get identity id")
	}

	res, err := h.SignInWithApple.ValidateRegistrationCode(ctx, privateKey, c1)
	if err != nil {
		if signinwithapple.IsInvalidGrant(err) {
			return inv.Error(err, 401, "Apple rejected the registration code, likely because it is expired")
		}
		return inv.Error(err, 502, "Failed to validate registration code")
	}

	newId := safeid.Make()

	// create a new user
	user := cwebauthn.NewUser([]byte(sub), u1)

	options, sessionData, err := h.WebAuthn.BeginRegistration(user)
	if err != nil {
		return inv.Error(err, 500, "failed to begin registration")
	}

	opts, err := json.Marshal(options)
	if err != nil {
		return inv.Error(err, 500, "Failed to marshal options")
	}

	newUser := dynamo.NewUser(newId.String(), u1, sub, creds, res)

	err = h.Dynamo.SaveChallenge(h.Ctx, sessionData, newUser)
	if err != nil {
		if dynamo.IsConditionalCheckFailed(err) {
			return inv.Error(err, 409, "User already exists")
		}
		return inv.Error(err, 502, "Failed to generate user")
	}

	return inv.Success(200, map[string]string{
		"Content-Type": "application/json",
	}, string(opts))
}
