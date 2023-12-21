package main

import (
	"context"

	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/storage"
	"github.com/walteh/webauthn/pkg/storage/dynamodb"
	"github.com/walteh/webauthn/pkg/webauthn/types"

	"os"

	"github.com/rs/xid"
	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	Id      string
	Ctx     context.Context
	Storage storage.Provider
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
		Id:      xid.New().String(),
		Ctx:     ctx,
		Storage: dynamodb.NewDynamoDBStorageClient(cfg, "ceremonies", "credentials"),
		Config:  cfg,
		logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	sessionId := hex.HexToHash(payload.Headers["x-nugg-hex-session-id"])
	ceremonyType := payload.Headers["x-nugg-utf-ceremony-type"]
	credentialId := hex.HexToHash(payload.Headers["x-nugg-hex-credential-id"])

	if sessionId.IsZero() {
		return Output{
			StatusCode: 400,
		}, nil
	}

	if ceremonyType == "" {
		ceremonyType = string(types.AssertCeremony)
	}

	switch ceremonyType {
	case string(types.AssertCeremony):
	case string(types.CreateCeremony):
		break
	default:
		return Output{
			StatusCode: 400,
		}, nil
	}

	cha := types.NewCeremony(types.CredentialID(credentialId), sessionId, types.CeremonyType(ceremonyType))

	err := h.Storage.WriteNewCeremony(ctx, cha)
	if err != nil {
		return Output{
			StatusCode: 500,
		}, nil
	}

	return Output{
		StatusCode: 204,
		Headers: map[string]string{
			"x-nugg-hex-challenge": cha.ChallengeID.Ref().Hex(),
		},
	}, nil
}

// cer, err := dynamo.MakePut(h.Dynamo.MustCeremonyTableName(), cha)
// if err != nil {
// 	return inv.Error(err, 500, "failed to create ceremony")
// }

// err = h.Dynamo.TransactWrite(ctx, *cer)
// if err != nil {
// 	return inv.Error(err, 500, "Failed to save ceremony")
// }
