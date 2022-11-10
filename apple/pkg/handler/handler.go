package handler

import (
	"context"
	"fmt"
	"log"
	"nugg-auth/apple/pkg/applepublickey"
	"nugg-auth/apple/pkg/cognito"
	"nugg-auth/apple/pkg/dynamo"
	"nugg-auth/apple/pkg/env"
	"nugg-auth/apple/pkg/secretsmanager"
	"nugg-auth/apple/pkg/signinwithapple"
)

type LambdaHander struct {
	ctx             context.Context
	dynamo          *dynamo.Client
	cognito         *cognito.Client
	signinwithapple *signinwithapple.Client
	applepublickey  *applepublickey.Client
	secretsmanager  *secretsmanager.Client
	env             env.Environment
}

type Request struct {
	AppleJwtToken string
	Operation     string
}

type Service interface {
	ParseRequest(handler LambdaHander, event map[string]interface{}) (Request, error)
	FormatResponse(handler LambdaHander, isAuthorized bool, result map[string]interface{}, err error) (interface{}, error)
}

func NewHandler(ctx context.Context, env env.Environment) (*LambdaHander, error) {
	return &LambdaHander{
		ctx:             ctx,
		dynamo:          dynamo.NewClient(env.AwsConfig, env.ChallengeTableName),
		cognito:         cognito.NewClient(env.AwsConfig, env.AppleIdentityPoolId),
		signinwithapple: signinwithapple.NewClient(env.AppleTokenEndpoint, env.AppleTeamID, env.AppleServiceName, env.SignInWithApplePrivateKeyID),
		applepublickey:  applepublickey.NewClient(env.ApplePublicKeyEndpoint),
		secretsmanager:  secretsmanager.NewClient(ctx, env.AwsConfig, env.SignInWithApplePrivateKeyName),
		env:             env,
	}, nil
}

// handlers

// 0 auth challenge /auth/challenge

// 1 apigw register /apple/auth/register

/// simple sign in check
// 2 apigw simple signinwithapple authorizer
// 3 appsync simple signinwithapple/passkey authorizer

// 4 apigw appattest authorizer

// 5 apigw

func GetService(event map[string]interface{}) Service {
	if event["authorizationToken"] == nil && event["requestContext"] == nil {
		return NoopHandler{}
	}

	if event["version"] != nil {
		return ApiGatewayV2AuthorizerService{}
	}

	return AppSyncHandler{}

}

func (handler LambdaHander) Run(ctx context.Context, event map[string]interface{}) (interface{}, error) {

	log.Println("event", event)
	service := GetService(event)

	request, err := service.ParseRequest(handler, event)
	if err != nil {
		return service.FormatResponse(handler, false, nil, err)
	}

	if request.AppleJwtToken == "" {
		return service.FormatResponse(handler, false, nil, fmt.Errorf("missing apple jwt token"))
	}

	publickey, err := handler.applepublickey.Refresh(ctx)
	if err != nil {
		return service.FormatResponse(handler, false, nil, err)
	}

	tkn, err := publickey.ParseToken(request.AppleJwtToken)
	if err != nil {
		return service.FormatResponse(handler, false, nil, err)
	}

	if !tkn.Valid {
		return service.FormatResponse(handler, false, nil, fmt.Errorf("invalid apple jwt token"))
	}

	sub, err := tkn.GetUniqueID()
	if err != nil {
		return service.FormatResponse(handler, false, nil, err)
	}

	creds, err := handler.cognito.GetIdentityId(handler.ctx, request.AppleJwtToken)
	if err != nil {
		return service.FormatResponse(handler, false, nil, err)
	}

	return service.FormatResponse(handler, true, map[string]interface{}{
		"sub":      sub,
		"identity": creds,
	}, nil)

}

// challenge, err := handler.dynamo.GenerateChallenge(ctx, userId, 60)
// if err != nil {
// 	return events.APIGatewayV2HTTPResponse{
// 		StatusCode: 500,
// 	}, err
// }

// log.Println("Challenge: " + challenge)
