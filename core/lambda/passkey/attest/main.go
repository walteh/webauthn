package main

import (
	"nugg-webauthn/core/pkg/cognito"
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/env"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/invocation"
	"nugg-webauthn/core/pkg/webauthn/clientdata"
	"nugg-webauthn/core/pkg/webauthn/credential"
	"nugg-webauthn/core/pkg/webauthn/handlers/devicecheck"
	"nugg-webauthn/core/pkg/webauthn/providers"
	"nugg-webauthn/core/pkg/webauthn/types"

	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	Id      string
	Ctx     context.Context
	Dynamo  *dynamo.Client
	Config  config.Config
	logger  zerolog.Logger
	Cognito cognito.Client
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
		Dynamo:  dynamo.NewClient(cfg, env.DynamoUsersTableName(), env.DynamoCeremoniesTableName(), env.DynamoCredentialsTableName()),
		Config:  cfg,
		Cognito: cognito.NewClient(cfg, env.AppleIdentityPoolId(), env.CognitoDeveloperProviderName()),
		logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, payload Input) (Output, error) {

	inv, ctx := invocation.NewInvocation(ctx, h, payload)

	attestation := hex.HexToHash(payload.Headers["x-nugg-hex-attestation-object"])
	dc := hex.HexToHash(payload.Headers["x-nugg-hex-devicecheck-assertion-object"])
	clientData := payload.Headers["x-nugg-utf-client-data-json"]
	credentialID := hex.HexToHash(payload.Headers["x-nugg-hex-credential-id"])

	if attestation.IsZero() || clientData == "" || credentialID.IsZero() {
		return inv.Error(nil, 400, "missing header x-nugg-webauthn-creation")
	}

	code, ok, err := devicecheck.Assert(ctx, h.Dynamo, dc, credentialID)
	if err != nil || !ok {
		return inv.Error(err, code, "devicecheck assertion failed")
	}

	parsedResponse := types.AttestationInput{
		AttestationObject:  attestation,
		UTF8ClientDataJSON: clientData,
		CredentialID:       credentialID,
		ClientExtensions:   nil,
	}

	cd, err := clientdata.ParseClientData(parsedResponse.UTF8ClientDataJSON)
	if err != nil {
		return inv.Error(err, 400, "invalid client data")
	}

	cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	err = h.Dynamo.TransactGet(ctx, cerem)
	if err != nil {
		return inv.Error(err, 500, "failed to get ceremony")
	}

	cred, invalidErr := credential.VerifyAttestationInput(types.VerifyAttestationInputArgs{
		Provider:           providers.NewNoneAttestationProvider(),
		Input:              parsedResponse,
		StoredChallenge:    cerem.ChallengeID,
		SessionId:          cerem.SessionID,
		VerifyUser:         false,
		RelyingPartyID:     "nugg.xyz",
		RelyingPartyOrigin: "https://nugg.xyz",
	})

	if invalidErr != nil {
		return inv.Error(invalidErr, 400, "invalid attestation")
	}

	chaner := make(chan *cognitoidentity.GetOpenIdTokenForDeveloperIdentityOutput, 1)
	defer close(chaner)
	stale := false
	var chanerr error

	go func() {
		go func() {
			<-ctx.Done()
			if !stale {
				chaner <- nil
			}
			stale = true

		}()

		z, err := h.Cognito.GetDevCreds(ctx, cerem.CredentialID)
		if !stale {
			if err != nil {
				chanerr = err
				chaner <- nil
			} else {
				chaner <- z
			}
		}
	}()

	credput, err := h.Dynamo.BuildPut(cred)
	if err != nil {
		return inv.Error(err, 500, "failed to create credential put")
	}

	ceremput, err := h.Dynamo.BuildPut(cerem)
	if err != nil {
		return inv.Error(err, 500, "failed to create ceremony put")
	}

	err = h.Dynamo.TransactWrite(ctx,
		dtypes.TransactWriteItem{Put: credput},
		dtypes.TransactWriteItem{Put: ceremput},
	)
	if err != nil {
		return inv.Error(err, 500, "failed to write to dynamo")
	}

	result := <-chaner
	stale = true

	if result == nil {
		return inv.Error(chanerr, 500, "failed to get dev creds")
	}

	return inv.Success(204, map[string]string{"x-nugg-utf-access-token": *result.Token}, "")
}
