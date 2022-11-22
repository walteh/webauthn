package main

import (
	"context"
	"fmt"
	"log"
	"nugg-webauthn/core/pkg/applepublickey"
	"nugg-webauthn/core/pkg/cognito"
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/env"
	"nugg-webauthn/core/pkg/secretsmanager"
	"nugg-webauthn/core/pkg/signinwithapple"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Input = events.APIGatewayV2CustomAuthorizerV2Request
type Output = events.APIGatewayV2CustomAuthorizerSimpleResponse

type Handler struct {
	Ctx             context.Context
	Dynamo          *dynamo.Client
	Cognito         cognito.Client
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
		Dynamo:          dynamo.NewClient(cfg, "", env.DynamoCeremoniesTableName(), ""),
		Cognito:         cognito.NewClient(cfg, env.AppleIdentityPoolId(), env.CognitoDeveloperProviderName()),
		SignInWithApple: signinwithapple.NewClient(env.AppleTokenEndpoint(), env.AppleTeamID(), env.AppleServiceName(), env.SignInWithApplePrivateKeyID()),
		ApplePublicKey:  applepublickey.NewClient(env.ApplePublicKeyEndpoint()),
		SecretsManager:  secretsmanager.NewClient(ctx, cfg, env.SignInWithApplePrivateKeyName()),
		Config:          cfg,
		counter:         0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Error(err error, message string) (Output, error) {

	return Output{
		IsAuthorized: false,
		Context:      map[string]interface{}{"error": err.Error(), "message": message},
	}, nil
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	h.counter++

	if payload.Headers["Authorization"] == "" {
		return h.Error(nil, "Missing required headers")
	}

	publickey, err := h.ApplePublicKey.Refresh(ctx)
	if err != nil {
		return h.Error(err, "Failed to refresh public key")
	}

	tkn, err := publickey.ParseToken(payload.Headers["Authorization"])
	if err != nil {
		return h.Error(err, "Failed to parse token")
	}

	if !tkn.Valid {
		return h.Error(fmt.Errorf("Unathorized"), "Invalid token")
	}

	sub, err := tkn.GetUniqueID()
	if err != nil {
		return h.Error(err, "Failed to get sub")
	}

	creds, err := h.Cognito.GetIdentityId(h.Ctx, payload.Headers["Authorization"])
	if err != nil {
		return h.Error(err, "Failed to get identity id")
	}

	return Output{
		IsAuthorized: true,
		Context: map[string]interface{}{
			"sub":    sub,
			"creds":  creds,
			"claims": tkn.Claims,
		},
	}, nil
}