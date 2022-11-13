package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"nugg-auth/core/pkg/applepublickey"
	"nugg-auth/core/pkg/cognito"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/random"
	"nugg-auth/core/pkg/secretsmanager"
	"nugg-auth/core/pkg/signinwithapple"
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
	Dynamo          *dynamo.Client
	Cognito         *cognito.Client
	SignInWithApple *signinwithapple.Client
	ApplePublicKey  *applepublickey.Client
	SecretsManager  *secretsmanager.Client
	Config          config.Config
	Logger          zerolog.Logger
	counter         int
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

	abc := &Handler{
		Id:              random.KSUID(),
		Ctx:             ctx,
		Dynamo:          dynamo.NewClient(cfg, env.DynamoUserTableName()),
		Cognito:         cognito.NewClient(cfg, env.AppleIdentityPoolId()),
		SignInWithApple: signinwithapple.NewClient(env.AppleTokenEndpoint(), env.AppleTeamID(), env.AppleServiceName(), env.SignInWithApplePrivateKeyID()),
		ApplePublicKey:  applepublickey.NewClient(env.ApplePublicKeyEndpoint()),
		SecretsManager:  secretsmanager.NewClient(ctx, cfg, env.SignInWithApplePrivateKeyName()),
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

	inv := h.NewInvocation(h.Logger)

	h1 := payload.Headers["x-nugg-signinwithapple-identity-token"]
	c1 := payload.Headers["x-nugg-signinwithapple-registration-code"]
	a1 := payload.Headers["Authorization"]

	if h1 == "" && c1 == "" && a1 != "" {

		var authToken struct {
			IdentityToken    string `json:"x-nugg-signinwithapple-identity-token"`
			RegistrationCode string `json:"x-nugg-signinwithapple-registration-code"`
		}

		// un base64 decode the token
		base64Token, err := base64.StdEncoding.DecodeString(a1)
		if err != nil {
			return inv.Error(err, 400, "Failed to decode")
		}

		if err := json.Unmarshal(base64Token, &authToken); err != nil {
			return inv.Error(err, 400, "Invalid authorization token")
		}

		h1 = authToken.IdentityToken
		c1 = authToken.RegistrationCode

	}

	if h1 == "" {
		return inv.Error(nil, 400, "Missing x-nugg-signinwithapple-identity-token")
	}

	if c1 == "" {
		return inv.Error(nil, 400, "Missing x-nugg-signinwithapple-registration-code")
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

	err = h.Dynamo.GenerateUser(h.Ctx, sub, creds, res)
	if err != nil {
		if dynamo.IsConditionalCheckFailed(err) {
			return inv.Error(err, 409, "User already exists")
		}
		return inv.Error(err, 502, "Failed to generate user")
	}

	return inv.Success(204, map[string]string{}, "")
}
