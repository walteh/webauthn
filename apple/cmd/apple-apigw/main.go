package main

import (
	"nugg-auth/apple/pkg/handler"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.NewDefaultHandler(Run))
}

type Input = events.APIGatewayV2CustomAuthorizerV2Request
type Output = events.APIGatewayV2CustomAuthorizerSimpleResponse

func Run(handler *handler.LambdaHander[Input, Output], event Input) (Output, error) {

	return Output{}, nil

}
