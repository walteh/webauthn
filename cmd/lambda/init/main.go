package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"

	"connectrpc.com/connect"
	"github.com/walteh/terrors"
	"github.com/walteh/webauthn/gen/buf/go/proto/webauthn/v1"
	"github.com/walteh/webauthn/gen/buf/go/proto/webauthn/v1/webauthnconnect"
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

type InputBody = connect.Request[webauthn.CreateChallengeRequest]

type OutputBody = connect.Response[webauthn.CreateChallengeResponse]

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

	_, hndl := webauthnconnect.NewCreateChallengeServiceHandler(abc)

	// convert http handler to lambda handler manually
	lmd := func(ctx context.Context, payload *Input) (*Output, error) {

		unb64, err := base64.RawStdEncoding.DecodeString(payload.Body)
		if err != nil {
			return nil, terrors.Wrap(err, "failed to decode body")
		}

		rwrit := &DummyHttpWriter{
			header: http.Header{},
			status: 200,
		}

		req, err := http.NewRequest("POST", "http://localhost:8080", bytes.NewReader(unb64))
		if err != nil {
			return nil, terrors.Wrap(err, "failed to create request")
		}

		hndl.ServeHTTP(rwrit, req)

		mvheaders := make(map[string][]string)
		for k, v := range rwrit.header {
			mvheaders[k] = v
		}

		return &Output{
			StatusCode:        rwrit.status,
			Body:              base64.StdEncoding.EncodeToString(rwrit.out),
			IsBase64Encoded:   true,
			MultiValueHeaders: mvheaders,
			Cookies:           []string{},
		}, nil
	}

	lambda.Start(lmd)
}

var (
	_ webauthnconnect.CreateChallengeServiceHandler = (*Handler)(nil)
)

func wrapConnect[Req any, Res any](f func(ctx context.Context, payload *Req) (*Res, error)) func(ctx context.Context, payload *connect.Request[Req]) (*connect.Request[Res], error) {
	return func(ctx context.Context, payload *connect.Request[Req]) (*connect.Request[Res], error) {
		res, err := f(ctx, payload.Msg)
		if err != nil {
			return nil, err
		}

		return &connect.Request[Res]{
			Msg: res,
		}, nil
	}
}

func (h *Handler) CreateChallenge(ctx context.Context, payload *InputBody) (*OutputBody, error) {

	sessionId := hex.BytesToHash(payload.Msg.GetSessionId())
	credentialId := hex.BytesToHash(payload.Msg.GetCredentialId())
	ceremonyType := payload.Msg.GetCeremonyType()

	if sessionId.IsZero() {
		return nil, terrors.New("session id is required")
	}

	if ceremonyType == "" {
		ceremonyType = string(types.AssertCeremony)
	}

	switch ceremonyType {
	case string(types.AssertCeremony):
	case string(types.CreateCeremony):
		break
	default:
		return nil, terrors.New("invalid ceremony type")
	}

	cha := types.NewCeremony(types.CredentialID(credentialId), sessionId, types.CeremonyType(ceremonyType))

	err := h.Storage.WriteNewCeremony(ctx, cha)
	if err != nil {
		return nil, err
	}

	return &connect.Response[webauthn.CreateChallengeResponse]{
		Msg: &webauthn.CreateChallengeResponse{
			Challenge: cha.ChallengeID,
			Ttl:       3 * 60 * 1000,
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

type DummyHttpWriter struct {
	header http.Header
	out    []byte
	status int
}

func (d *DummyHttpWriter) Header() http.Header {
	return d.header
}

func (d *DummyHttpWriter) Write(ok []byte) (int, error) {
	d.out = ok
	return len(ok), nil
}

func (d *DummyHttpWriter) WriteHeader(stat int) {
	d.status = stat
}
