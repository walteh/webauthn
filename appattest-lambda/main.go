package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LambdaHander struct {
	ctx        context.Context
	httpClient *http.Client
}

func NewHandler(ctx context.Context, httpClient *http.Client) *LambdaHander {
	return &LambdaHander{
		ctx:        ctx,
		httpClient: httpClient,
	}
}

func (handler LambdaHander) Run(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (string, error) {
	return "", nil
}

func main() {
	handler := NewHandler(context.Background(), &http.Client{})
	lambda.Start(handler.Run)
}
