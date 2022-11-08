package main

import (
	"context"
	"fmt"
	"log"
	"nugg-auth/apple/pkg/env"
	"nugg-auth/apple/pkg/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	log.Println("Starting lambda handler")

	ctx := context.Background()

	env, err := env.NewEnv(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Environment loaded")

	handler, err := handler.NewHandler(ctx, env)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("Handler created")

	lambda.Start(handler.Run)
}
