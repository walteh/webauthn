package main

import (
	"context"
	"nugg-crypto/go/pkg/env"
	"nugg-crypto/go/pkg/handler"

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
