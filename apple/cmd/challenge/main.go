package main

import (
	"log"
	"nugg-auth/apple/pkg/handler"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Input = events.AppSyncLambdaAuthorizerRequest
type Output = events.AppSyncLambdaAuthorizerResponse
type AppsyncHandler = handler.LambdaHander[Input, Output]

func main() {
	lambda.Start(handler.NewDefaultHandler(Run))
}

func Run(handler *AppsyncHandler, event Input) (Output, error) {

	return Output{}, nil

}

type Request struct {
	events.APIGatewayV2HTTPRequest
}

func Invoke(h *AppsyncHandler, payload events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	if payload.Headers["x-nugg-challenge-state"] == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       "Missing required headers",
		}, nil
	}

	type Input = events.AppSyncLambdaAuthorizerRequest
	type Output = events.AppSyncLambdaAuthorizerResponse
	type AppsyncHandler = handler.LambdaHander[Input, Output]

	// get challenge from dynamo
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
		Headers: map[string]string{
			"x-nugg-challenge": challenge,
		},
	}, nil
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC | log.Lmicroseconds)
}
