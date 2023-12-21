package main

import (
	"github.com/walteh/webauthn/app/devicecheck_attest"
	"github.com/walteh/webauthn/pkg/aws/invocation"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/storage"

	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	*invocation.Handler[Input, Output]
	Id      string
	Ctx     context.Context
	Dynamo  storage.Provider
	Config  config.Config
	logger  zerolog.Logger
	counter int
}

func (h Handler) ID() string {
	return h.Id
}

func (h *Handler) IncrementCounter() int {
	h.counter += 1
	return h.counter
}

func main() {

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}

	abc := &Handler{
		Id:      xid.New().String(),
		Ctx:     ctx,
		Dynamo:  dynamo.NewClient(cfg, "", env.DynamoCeremoniesTableName(), env.DynamoCredentialsTableName()),
		Config:  cfg,
		logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, input Input) (Output, error) {

	inv, ctx := h.NewInvocation(ctx, input)

	attestation := hex.HexToHash(input.Body)
	attestationKey := hex.HexToHash(input.Headers["x-nugg-hex-attestation-key"])
	clientDataJson := input.Headers["x-nugg-utf-client-data-json"]
	sessionId := hex.HexToHash(input.Headers["x-nugg-hex-session-id"])

	if attestation.IsZero() || clientDataJson == "" {
		return inv.Error(nil, 400, "missing required headers")
	}

	res, err := devicecheck_attest.Attest(ctx, h.Dynamo, devicecheck_attest.DeviceCheckAttestationInput{
		RawAttestationObject: attestation,
		UTF8ClientDataJSON:   clientDataJson,
		RawCredentialID:      attestationKey,
		RawSessionID:         sessionId,
	})
	if err != nil || !res.OK {
		if err == nil {
			err = errors.ErrVerification
		}
		return inv.Error(err, res.SuggestedStatusCode, "failed to attest devicecheck")
	}

	return inv.Success(Output{
		StatusCode: 204,
	})

}
