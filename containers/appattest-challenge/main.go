package main

import (
	"appattest-challenge/pkg/env"
	"appattest-challenge/pkg/handler"
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	ctx := context.Background()

	env, err := env.NewEnv(ctx)
	if err != nil {
		panic(err)
	}

	handler := handler.NewHandler(ctx, env)

	lambda.Start(handler.Run)
}
