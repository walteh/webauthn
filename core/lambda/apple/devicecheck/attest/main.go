package main

import (
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/env"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/invocation"
	protocol "nugg-webauthn/core/pkg/webauthn"

	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	Id      string
	Ctx     context.Context
	Dynamo  *dynamo.Client
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

func (h Handler) Logger() zerolog.Logger {
	return h.logger
}

func main() {

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}

	abc := &Handler{
		Id:      ksuid.New().String(),
		Ctx:     ctx,
		Dynamo:  dynamo.NewClient(cfg, "", env.DynamoCeremoniesTableName(), env.DynamoCredentialsTableName()),
		Config:  cfg,
		logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, input Input) (Output, error) {

	inv, ctx := invocation.NewInvocation(ctx, h, input)

	attestation := hex.HexToHash(input.Headers["x-nugg-hex-attestation"])
	clientData := input.Headers["x-nugg-utf-client-data-json"]
	payload := hex.HexToHash(input.Headers["x-nugg-hex-payload"])

	if attestation.IsZero() || clientData == "" {
		return inv.Error(nil, 400, "missing required headers")
	}

	p, err := protocol.FormatAttestationInput(clientData, attestation).Parse()
	if err != nil {
		return inv.Error(err, 400, err.Error())
	}

	cer := protocol.NewUnsafeGettableCeremony(p.CollectedClientData.Challenge)

	err = h.Dynamo.TransactGet(ctx, cer)
	if err != nil {
		return inv.Error(err, 400, err.Error())
	}

	if !cer.ChallengeID.Equals(p.CollectedClientData.Challenge) {
		return inv.Error(nil, 400, "invalid credential id")
	}

	pk, err := p.AttestationObject.Verify("4497QJSAD3.xyz.nugg.app", payload.Sha256(), false, false)
	if err != nil {
		return inv.Error(err, 400, err.Error())
	}

	putter, err := dynamo.MakePut(h.Dynamo.MustCredentialTableName(), pk)
	if err != nil {
		return inv.Error(err, 500, err.Error())
	}

	err = h.Dynamo.TransactWrite(ctx, types.TransactWriteItem{Put: putter})
	if err != nil {
		return inv.Error(err, 500, err.Error())
	}

	return inv.Success(204, map[string]string{}, "")

}
