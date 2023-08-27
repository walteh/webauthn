package main

import (
	"context"

	passkey "github.com/walteh/webauthn/app/passkey_assert"
	"github.com/walteh/webauthn/pkg/cognito"
	"github.com/walteh/webauthn/pkg/env"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/indexable"
	"github.com/walteh/webauthn/pkg/invocation"

	"github.com/aws/aws-lambda-go/events"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler[D indexable.DynamoDBAPI] struct {
	*invocation.Handler[Input, Output]

	Dynamo  D
	Cognito cognito.Client
}

func buildHandler[D indexable.DynamoDBAPI](h *invocation.Handler[Input, Output], d D, c cognito.Client) (*Handler[D], error) {

	abc := &Handler[D]{
		Handler: h,
		Dynamo:  d,
		Cognito: c,
	}

	return abc, nil
}

func NewHandler() (*Handler[indexable.DynamoDBAPI], error) {

	ctx := context.Background()

	handler := invocation.NewHandler[Input, Output](ctx)

	dbc := handler.Opts().NewDynamoDBClient()

	api := indexable.NewDynamoDBAPI(dbc, "")

	cog := cognito.NewClient(*handler.Opts().AwsConfig(), env.AppleIdentityPoolId(), env.CognitoDeveloperProviderName())

	abc, _ := buildHandler(handler, api, cog)

	return abc, nil
}

func (h *Handler[D]) Invoke(ctx context.Context, input Input) (Output, error) {

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
