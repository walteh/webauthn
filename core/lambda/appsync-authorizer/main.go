package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"nugg-auth/core/pkg/applepublickey"
	"nugg-auth/core/pkg/cognito"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/env"
	"nugg-auth/core/pkg/secretsmanager"
	"nugg-auth/core/pkg/signinwithapple"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Input = events.AppSyncLambdaAuthorizerRequest
type Output = events.AppSyncLambdaAuthorizerResponse

type Handler struct {
	Ctx             context.Context
	Dynamo          *dynamo.Client
	Cognito         *cognito.Client
	SignInWithApple *signinwithapple.Client
	ApplePublicKey  *applepublickey.Client
	SecretsManager  *secretsmanager.Client
	Config          config.Config
	counter         int
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC | log.Lmicroseconds)
}

func main() {

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}

	abc := &Handler{
		Ctx:             ctx,
		Dynamo:          dynamo.NewClient(cfg, env.DynamoChallengeTableName()),
		Cognito:         cognito.NewClient(cfg, env.AppleIdentityPoolId()),
		SignInWithApple: signinwithapple.NewClient(env.AppleTokenEndpoint(), env.AppleTeamID(), env.AppleServiceName(), env.SignInWithApplePrivateKeyID()),
		ApplePublicKey:  applepublickey.NewClient(env.ApplePublicKeyEndpoint()),
		SecretsManager:  secretsmanager.NewClient(ctx, cfg, env.SignInWithApplePrivateKeyName()),
		Config:          cfg,
		counter:         0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Error(err error, message string) (Output, error) {
	log.Printf("Error: %s %s", message, err)
	return Output{
		IsAuthorized:    false,
		ResolverContext: map[string]interface{}{"error": err.Error(), "message": message},
	}, nil
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	h.counter++

	log.Printf("Invoke: %d %s", h.counter, payload)

	if payload.AuthorizationToken == "" {
		return h.Error(nil, "Missing required headers")
	}

	var authToken struct {
		Authorization string `json:"x-nugg-authorization"`
	}

	// un base64 decode the token
	base64Token, err := base64.StdEncoding.DecodeString(payload.AuthorizationToken)
	if err != nil {
		return h.Error(err, "Failed to decode")
	}

	if err := json.Unmarshal(base64Token, &authToken); err != nil {
		return h.Error(err, "Invalid authorization token")
	}

	publickey, err := h.ApplePublicKey.Refresh(ctx)
	if err != nil {
		return h.Error(err, "Failed to refresh public key")
	}

	tkn, err := publickey.ParseToken(authToken.Authorization)
	if err != nil {
		return h.Error(err, "Failed to parse token")
	}

	if !tkn.Valid {
		return h.Error(err, "Invalid token")
	}

	sub, err := tkn.GetUniqueID()
	if err != nil {
		return h.Error(err, "Failed to get sub")
	}

	creds, err := h.Cognito.GetIdentityId(h.Ctx, authToken.Authorization)
	if err != nil {
		return h.Error(err, "Failed to get identity id")
	}

	return Output{
		IsAuthorized: true,
		ResolverContext: map[string]interface{}{
			"sub":   sub,
			"creds": creds,
		},
	}, nil
}
