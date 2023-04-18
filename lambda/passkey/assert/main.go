package main

import (
	"context"

	"git.nugg.xyz/go-sdk/invocation"
	"git.nugg.xyz/go-sdk/x"
	"git.nugg.xyz/webauthn/pkg/cognito"
	"git.nugg.xyz/webauthn/pkg/env"
	"git.nugg.xyz/webauthn/pkg/hex"
	"git.nugg.xyz/webauthn/pkg/passkey"

	"github.com/aws/aws-lambda-go/events"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	*invocation.Handler[Input, Output]

	Dynamo  x.DynamoDBAPI
	Cognito cognito.Client
}

func NewHandler() (*Handler, error) {

	ctx := context.Background()

	handler := invocation.NewHandler[Input, Output](ctx)

	dbc := handler.Opts().NewDynamoDBClient()

	api := x.NewDynamoDBAPI(dbc, "")

	abc := &Handler{
		Dynamo:  api,
		Cognito: cognito.NewClient(*handler.Opts().AwsConfig(), env.AppleIdentityPoolId(), env.CognitoDeveloperProviderName()),
	}

	return abc, nil
}

func (h *Handler) Invoke(ctx context.Context, input Input) (Output, error) {

	inv, ctx := h.NewInvocation(ctx, input)

	authenticatorData := hex.HexToHash(input.Headers["x-nugg-hex-authenticator-data"])
	credentialId := hex.HexToHash(input.Headers["x-nugg-hex-credential-id"])
	signature := hex.HexToHash(input.Headers["x-nugg-hex-signature"])
	userId := hex.HexToHash(input.Headers["x-nugg-hex-user-id"])

	clientDataJson := input.Headers["x-nugg-utf-client-data-json"]
	credentialType := input.Headers["x-nugg-utf-credential-type"]

	// make sure all the above values exist one by one in the headers
	if len(authenticatorData) == 0 || len(credentialId) == 0 || len(signature) == 0 || len(userId) == 0 || len(clientDataJson) == 0 || len(credentialType) == 0 {
		return inv.Error(nil, 400, "missing required headers")
	}

	abc := passkey.PasskeyAssertionInput{
		RawAuthenticatorData: authenticatorData,
		CredentialID:         credentialId,
		RawSignature:         signature,
		UTF8ClientDataJSON:   clientDataJson,
		SessionID:            userId,
	}

	res, err := passkey.Assert(ctx, h.Dynamo, h.Cognito, abc)
	if err != nil {
		return inv.Error(err, res.SuggestedStatusCode, "failed to assert passkey")
	}

	return inv.Success(Output{
		Headers: map[string]string{
			"x-nugg-utf-access-token": res.AccessToken,
		},
		StatusCode: 204,
	})
}
