package main

import (
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/env"
	"nugg-webauthn/core/pkg/errors"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/invocation"
	protocol "nugg-webauthn/core/pkg/webauthn"
	attestation_providers "nugg-webauthn/core/pkg/webauthn/providers"

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

	attestation := hex.HexToHash(input.Body)
	attestationKey := hex.HexToHash(input.Headers["x-nugg-hex-attestation-key"])
	clientDataJson := input.Headers["x-nugg-utf-client-data-json"]
	sessionId := hex.HexToHash(input.Headers["x-nugg-hex-session-id"])

	if attestation.IsZero() || clientDataJson == "" {
		return inv.Error(nil, 400, "missing required headers")
	}

	p, err := protocol.FormatAttestationInput(clientDataJson, attestation).Parse()
	if err != nil {
		return inv.Error(err, 400, err.Error())
	}

	cer := protocol.NewUnsafeGettableCeremony(p.CollectedClientData.Challenge)

	err = h.Dynamo.TransactGet(ctx, cer)
	if err != nil {
		return inv.Error(err, 401, err.Error())
	}

	if !cer.ChallengeID.Equals(p.CollectedClientData.Challenge) {
		return inv.Error(nil, 401, "invalid credential id")
	}

	if p.CollectedClientData.Origin != "https://nugg.xyz" {
		return inv.Error(nil, 401, "invalid origin")
	}

	if p.CollectedClientData.Type != protocol.CreateCeremony || cer.CeremonyType != protocol.CreateCeremony {
		return inv.Error(nil, 401, "invalid ceremony type")
	}

	if !cer.SessionID.Equals(sessionId) {
		return inv.Error(nil, 401, "invalid session id")
	}

	provider := attestation_providers.NewAppAttest()

	pk, err := p.AttestationObject.Verify(provider, "4497QJSAD3.xyz.nugg.app", hex.Hash([]byte(clientDataJson)).Sha256(), false, false)
	if err != nil {
		return inv.Error(err, 401, err.Error())
	}

	if !attestationKey.Equals(pk.RawID) {

		err := errors.NewError(0x93).WithCaller().
			WithKV("attestationKey", attestationKey.Hex()).
			WithKV("pk.RawID", pk.RawID.Hex())

		return inv.Error(err, 401, "invalid public key")
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
