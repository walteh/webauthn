package handler

import (
	"appattest-challenge/pkg/cognito"
	"appattest-challenge/pkg/dynamo"
	"appattest-challenge/pkg/env"
	"appattest-challenge/pkg/jwt"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type LambdaHander struct {
	ctx     context.Context
	dynamo  *dynamo.Client
	cognito *cognito.Client
	jwt     *jwt.Client
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
	}
}

func (handler LambdaHander) Run(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	userId := event.Headers["x-nugg-apple-userid"]

	token := event.Headers["x-nugg-apple-token"]

	if userId == "" || token == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       "Missing required headers",
		}, nil
	}

	// challenge, err := handler.dynamo.GenerateChallenge(ctx, userId, 60)
	// if err != nil {
	// 	return events.APIGatewayV2HTTPResponse{
	// 		StatusCode: 500,
	// 	}, err
	// }

	// log.Println("Challenge: " + challenge)

	creds, err := handler.cognito.GetIdentityId(ctx, token)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
		}, err
	}

	log.Println("Creds: ", creds)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"x-nugg-challenge": "a",
		},
		Body: creds,
	}, nil

}
