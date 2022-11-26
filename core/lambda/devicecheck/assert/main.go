package main

import (
	"log"
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/env"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/invocation"
	"nugg-webauthn/core/pkg/webauthn/assertion"
	"nugg-webauthn/core/pkg/webauthn/clientdata"
	"nugg-webauthn/core/pkg/webauthn/extensions"
	"nugg-webauthn/core/pkg/webauthn/providers"
	"nugg-webauthn/core/pkg/webauthn/types"

	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
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

	assert := hex.HexToHash(input.Headers["x-nugg-hex-request-assertion"])

	var body hex.Hash
	var err error

	if input.IsBase64Encoded {
		body, err = hex.Base64ToHash(input.Body)
		if err != nil {
			return inv.Error(err, 400, "failed to parse assertion")
		}
	} else {
		body = hex.HexToHash(input.Body)
	}

	if assert.IsZero() || len(body) == 0 {
		return inv.Error(nil, 400, "missing required headers")
	}

	parsed, err := assertion.ParseFidoAssertionInput(assert)
	if err != nil {
		return inv.Error(err, 400, "failed to parse assertion")
	}

	cd, err := clientdata.ParseClientData(parsed.RawClientDataJSON)
	if err != nil {
		return inv.Error(err, 400, "failed to parse client data")
	}

	cred := types.NewUnsafeGettableCredential(parsed.CredentialID)
	cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	err = h.Dynamo.TransactGet(ctx, cred, cerem)
	if err != nil {
		return inv.Error(err, 500, "failed to send transact get")
	}

	if cred.RawID.Hex() != cerem.CredentialID.Hex() {
		log.Println(cred.RawID.Hex(), cerem.CredentialID.Hex())
		return inv.Error(nil, 400, "invalid credential id")
	}

	if !cerem.ChallengeID.Equals(cd.Challenge) {
		return inv.Error(nil, 400, "invalid challenge id")
	}

	attestationProvider := providers.NewAppAttestSandbox()

	if attestationProvider.ID() != cred.AttestationType {
		return inv.Error(nil, 400, "invalid attestation provider")
	}

	// Handle steps 4 through 16
	if validError := assertion.VerifyAssertionInput(types.VerifyAssertionInputArgs{
		Input:                          parsed,
		StoredChallenge:                cerem.ChallengeID,
		RelyingPartyID:                 "4497QJSAD3.xyz.nugg.app",
		RelyingPartyOrigin:             env.RPOrigin(),
		AAGUID:                         cred.AAGUID,
		CredentialAttestationType:      types.FidoAttestationType,
		AttestationProvider:            attestationProvider,
		VerifyUser:                     false,
		CredentialPublicKey:            cred.PublicKey,
		Extensions:                     extensions.ClientInputs{},
		DataSignedByClient:             append(body, cerem.ChallengeID...),
		UseSavedAttestedCredentialData: true,
	}); validError != nil {
		return inv.Error(validError, 400, "failed to verify assertion")
	}

	credentialUpdate, err := cred.UpdateIncreasingCounter(h.Dynamo.MustCredentialTableName())
	if err != nil {
		return inv.Error(err, 500, "failed to create apple pass key")
	}

	err = h.Dynamo.TransactWrite(ctx, *credentialUpdate)
	if err != nil {
		return inv.Error(err, 500, err.Error())
	}

	return inv.Success(204, map[string]string{}, "")

}
