package main

import (
	"nugg-auth/apple/pkg/handler"

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
