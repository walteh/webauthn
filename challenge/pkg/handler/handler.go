package handler

import (
	"context"
	"encoding/json"
	"nugg-auth/challenge/pkg/dynamo"
	"nugg-auth/challenge/pkg/env"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

// Path: challenge/pkg/handler/apigwauthorizer.go

type LambdaHander struct {
	ctx    context.Context
	dynamo *dynamo.Client
	env    env.Environment
}

func NewHandler(ctx context.Context, env env.Environment) (*LambdaHander, error) {
	return &LambdaHander{
		ctx:    ctx,
		dynamo: dynamo.NewClient(env.AwsConfig, env.ChallengeTableName),
		env:    env,
	}, nil
}

type Request struct {
	events.APIGatewayV2HTTPRequest
}

func (h LambdaHander) Invoke(ctx context.Context, payload []byte) ([]byte, error) {

	var event events.APIGatewayV2HTTPRequest

	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, err
	}

	res, err := h.invoke(ctx, event)
	if err != nil {
		return nil, err
	}

	return json.Marshal(res)
}

func (h LambdaHander) invoke(ctx context.Context, payload events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	if payload.Headers["x-nugg-challenge-state"] == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       "Missing required headers",
		}, nil
	}

	// get challenge from dynamo
	challenge, err := h.dynamo.GenerateChallenge(ctx, payload.Headers["x-nugg-challenge-state"], time.Now().Add(time.Minute*5).Unix())
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 204,
		Headers: map[string]string{
			"x-nugg-challenge": challenge,
		},
	}, nil
}
