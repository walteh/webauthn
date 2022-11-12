package main

import (
	"context"
	"log"
	"nugg-auth/apple/pkg/dynamo"
	"nugg-auth/apple/pkg/env"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	Ctx     context.Context
	Dynamo  *dynamo.Client
	Config  config.Config
	counter int
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
		Ctx:     ctx,
		Dynamo:  dynamo.NewClient(cfg, env.DynamoChallengeTableName()),
		Config:  cfg,
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {
	h.counter++

	log.Printf("counter: %d", h.counter)

	if payload.Headers["x-nugg-challenge-state"] == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       "Missing required headers",
		}, nil
	}

	challenge, err := h.Dynamo.GenerateChallenge(
		h.Ctx,
		payload.Headers["x-nugg-challenge-state"],
		time.Now().Add(time.Minute*5),
	)

	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       "Error generating challenge",
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 204,
		Headers:    map[string]string{"x-nugg-challenge": challenge},
	}, nil
}
