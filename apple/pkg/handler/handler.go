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
	Env             env.Environment
}

type Request struct {
	AppleJwtToken string
	Operation     string
}

// type Service interface {
// 	ParseRequest(handler LambdaHander, event map[string]interface{}) (Request, error)
// 	FormatResponse(handler LambdaHander, isAuthorized bool, result map[string]interface{}, err error) (interface{}, error)
// }

func NewHandler[I interface{}, O interface{}](ctx context.Context, env env.Environment, invoker func(handler *LambdaHander[I, O], input I) (O, error)) func(ctx context.Context, input I) (O, error) {
	abc := &LambdaHander[I, O]{
		Ctx:             ctx,
		Dynamo:          dynamo.NewClient(env.AwsConfig, env.ChallengeTableName),
		Cognito:         cognito.NewClient(env.AwsConfig, env.AppleIdentityPoolId),
		SignInWithApple: signinwithapple.NewClient(env.AppleTokenEndpoint, env.AppleTeamID, env.AppleServiceName, env.SignInWithApplePrivateKeyID),
		ApplePublicKey:  applepublickey.NewClient(env.ApplePublicKeyEndpoint),
		SecretsManager:  secretsmanager.NewClient(ctx, env.AwsConfig, env.SignInWithApplePrivateKeyName),
		Env:             env,
	}

	return func(ctx context.Context, input I) (O, error) {
		return invoker(abc, input)
	}
}

func NewDefaultHandler[I interface{}, O interface{}](invoker func(handler *LambdaHander[I, O], input I) (O, error)) func(ctx context.Context, input I) (O, error) {
	ctx := context.Background()

	env, err := env.NewEnv(ctx)
	if err != nil {
		panic(err)
	}

	return NewHandler(ctx, env, invoker)
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
