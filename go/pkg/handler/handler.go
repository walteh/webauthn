package handler

import (
	"context"
	"fmt"
	"nugg-crypto/go/pkg/cognito"
	"nugg-crypto/go/pkg/dynamo"
	"nugg-crypto/go/pkg/env"
	"nugg-crypto/go/pkg/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type LambdaHander struct {
	ctx     context.Context
	dynamo  *dynamo.Client
	cognito *cognito.Client
	jwt     *jwt.Client
	env     env.Environment
}

type Request struct {
	AppleJwtToken string
	Operation     string
}

type Service interface {
	ParseRequest(handler LambdaHander, event interface{}) (Request, error)
	FormatResponse(handler LambdaHander, isAuthorized bool, result map[string]interface{}, err error) (interface{}, error)
}

func NewHandler(ctx context.Context, env env.Environment) *LambdaHander {

	jwtclient, err := jwt.NewAppleClient(ctx, env.AppleJwtPublicKeyEndpoint)
	if err != nil {
		panic(err)
	}

	return &LambdaHander{
		ctx:     ctx,
		dynamo:  dynamo.NewClient(env.AwsConfig, env.ChallengeTableName),
		cognito: cognito.NewClient(env.AwsConfig, env.AppleIdentityPoolName),
		jwt:     jwtclient,
		env:     env,
	}
}

func GetService(event interface{}) Service {
	switch event.(type) {
	case events.AppSyncLambdaAuthorizerRequest:
		return &AppSyncHandler{}
	case events.APIGatewayV2CustomAuthorizerV2Request:
		return &ApigwV2Handler{}
	default:
		return &NoopHandler{}
	}
}

func (handler LambdaHander) Run(ctx lambdacontext.LambdaContext, event interface{}) (interface{}, error) {

	service := GetService(event)

	request, err := service.ParseRequest(handler, event)
	if err != nil {
		return service.FormatResponse(handler, false, nil, err)
	}

	if request.AppleJwtToken == "" {
		return service.FormatResponse(handler, false, nil, fmt.Errorf("missing apple jwt token"))
	}

	tkn, err := handler.jwt.Verify(handler.ctx, request.AppleJwtToken)
	if err != nil {
		return service.FormatResponse(handler, false, nil, err)
	}

	if !tkn.Valid {
		return service.FormatResponse(handler, false, nil, fmt.Errorf("invalid apple jwt token"))
	}

	sub, err := tkn.Sub()
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
