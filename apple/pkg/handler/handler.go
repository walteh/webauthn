package handler

import (
	"context"
	"log"
	"nugg-auth/apple/pkg/applepublickey"
	"nugg-auth/apple/pkg/cognito"
	"nugg-auth/apple/pkg/dynamo"
	"nugg-auth/apple/pkg/env"
	"nugg-auth/apple/pkg/secretsmanager"
	"nugg-auth/apple/pkg/signinwithapple"

	"github.com/aws/aws-sdk-go-v2/config"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC | log.Lmicroseconds)
}

type LambdaHander[I interface{}, O interface{}] struct {
	Ctx             context.Context
	Dynamo          *dynamo.Client
	Cognito         *cognito.Client
	SignInWithApple *signinwithapple.Client
	ApplePublicKey  *applepublickey.Client
	SecretsManager  *secretsmanager.Client
	Config          config.Config
}

type Request struct {
	AppleJwtToken string
	Operation     string
}

// type Service interface {
// 	ParseRequest(handler LambdaHander, event map[string]interface{}) (Request, error)
// 	FormatResponse(handler LambdaHander, isAuthorized bool, result map[string]interface{}, err error) (interface{}, error)
// }

func NewHandler[I interface{}, O interface{}](ctx context.Context, invoker func(handler *LambdaHander[I, O], input I) (O, error)) func(ctx context.Context, input I) (O, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil
	}

	abc := &LambdaHander[I, O]{
		Ctx:             ctx,
		Dynamo:          dynamo.NewClient(cfg, env.DynamoChallengeTableName()),
		Cognito:         cognito.NewClient(cfg, env.AppleIdentityPoolId()),
		SignInWithApple: signinwithapple.NewClient(env.AppleTokenEndpoint(), env.AppleTeamID(), env.AppleServiceName(), env.SignInWithApplePrivateKeyID()),
		ApplePublicKey:  applepublickey.NewClient(env.ApplePublicKeyEndpoint()),
		SecretsManager:  secretsmanager.NewClient(ctx, cfg, env.SignInWithApplePrivateKeyName()),
		Config:          cfg,
	}

	return func(ctx context.Context, input I) (O, error) {
		return invoker(abc, input)
	}
}

func NewDefaultHandler[I interface{}, O interface{}](invoker func(handler *LambdaHander[I, O], input I) (O, error)) func(ctx context.Context, input I) (O, error) {
	ctx := context.Background()

	return NewHandler(ctx, invoker)
}

type Stripped struct {
	Dynamo          bool
	Cognito         bool
	SignInWithApple bool
	ApplePublicKey  bool
	SecretsManager  bool
}

func NewStrippedHandler[I interface{}, O interface{}](str Stripped, invoker func(handler *LambdaHander[I, O], input I) (O, error)) func(ctx context.Context, input I) (O, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil
	}

	abc := &LambdaHander[I, O]{
		Ctx:    ctx,
		Config: cfg,
	}

	if str.Dynamo {
		abc.Dynamo = dynamo.NewClient(cfg, env.DynamoChallengeTableName())
	}

	if str.Cognito {
		abc.Cognito = cognito.NewClient(cfg, env.AppleIdentityPoolId())
	}

	if str.SignInWithApple {
		abc.SignInWithApple = signinwithapple.NewClient(env.AppleTokenEndpoint(), env.AppleTeamID(), env.AppleServiceName(), env.SignInWithApplePrivateKeyID())
	}

	if str.ApplePublicKey {
		abc.ApplePublicKey = applepublickey.NewClient(env.ApplePublicKeyEndpoint())
	}

	if str.SignInWithApple || str.SecretsManager {
		abc.SecretsManager = secretsmanager.NewClient(ctx, cfg, env.SignInWithApplePrivateKeyName())
	}

	return func(ctx context.Context, input I) (O, error) {
		return invoker(abc, input)
	}
}

// handlers

// 0 auth challenge /auth/challenge

// 1 apigw register /apple/auth/register

/// simple sign in check
// 2 apigw simple signinwithapple authorizer
// 3 appsync simple signinwithapple/passkey authorizer

// 4 apigw appattest authorizer

// 5 apigw

// func GetService(event map[string]interface{}) Service {
// 	if event["authorizationToken"] == nil && event["requestContext"] == nil {
// 		return NoopHandler{}
// 	}

// 	if event["version"] != nil {
// 		return ApiGatewayV2AuthorizerService{}
// 	}

// 	return AppSyncHandler{}

// }

// func (handler LambdaHander) Run(ctx context.Context, event map[string]interface{}) (interface{}, error) {

// 	log.Println("event", event)
// 	service := GetService(event)

// 	request, err := service.ParseRequest(handler, event)
// 	if err != nil {
// 		return service.FormatResponse(handler, false, nil, err)
// 	}

// 	if request.AppleJwtToken == "" {
// 		return service.FormatResponse(handler, false, nil, fmt.Errorf("missing apple jwt token"))
// 	}

// 	publickey, err := handler.applepublickey.Refresh(ctx)
// 	if err != nil {
// 		return service.FormatResponse(handler, false, nil, err)
// 	}

// 	tkn, err := publickey.ParseToken(request.AppleJwtToken)
// 	if err != nil {
// 		return service.FormatResponse(handler, false, nil, err)
// 	}

// 	if !tkn.Valid {
// 		return service.FormatResponse(handler, false, nil, fmt.Errorf("invalid apple jwt token"))
// 	}

// 	sub, err := tkn.GetUniqueID()
// 	if err != nil {
// 		return service.FormatResponse(handler, false, nil, err)
// 	}

// 	creds, err := handler.cognito.GetIdentityId(handler.ctx, request.AppleJwtToken)
// 	if err != nil {
// 		return service.FormatResponse(handler, false, nil, err)
// 	}

// 	return service.FormatResponse(handler, true, map[string]interface{}{
// 		"sub":      sub,
// 		"identity": creds,
// 	}, nil)

// }

// challenge, err := handler.dynamo.GenerateChallenge(ctx, userId, 60)
// if err != nil {
// 	return events.APIGatewayV2HTTPResponse{
// 		StatusCode: 500,
// 	}, err
// }

// log.Println("Challenge: " + challenge)
