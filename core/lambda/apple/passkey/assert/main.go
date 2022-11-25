package main

import (
	"context"
	"log"
	"nugg-webauthn/core/pkg/cognito"
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/env"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/invocation"
	"nugg-webauthn/core/pkg/webauthn/assertion"
	"nugg-webauthn/core/pkg/webauthn/clientdata"
	"nugg-webauthn/core/pkg/webauthn/extensions"
	"nugg-webauthn/core/pkg/webauthn/types"

	"os"

	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Handler struct {
	Id      string
	Dynamo  *dynamo.Client
	Config  config.Config
	logger  zerolog.Logger
	Cognito cognito.Client
	Ctx     context.Context
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

	if err != nil {
		return
	}

	abc := &Handler{
		Id:      ksuid.New().String(),
		Dynamo:  dynamo.NewClient(cfg, env.DynamoUsersTableName(), env.DynamoCeremoniesTableName(), env.DynamoCredentialsTableName()),
		Config:  cfg,
		Cognito: cognito.NewClient(cfg, env.AppleIdentityPoolId(), env.CognitoDeveloperProviderName()),
		logger:  zerolog.New(os.Stdout).With().Caller().Timestamp().Logger(),
		counter: 0,
	}

	lambda.Start(abc.Invoke)
}

func (h *Handler) Invoke(ctx context.Context, input Input) (Output, error) {

	inv, ctx := invocation.NewInvocation(ctx, h, input)

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

	abc := types.AssertionInput{
		CredentialID:         credentialId,
		Signature:            signature,
		UserID:               userId,
		RawClientDataJSON:    clientDataJson,
		RawAuthenticatorData: authenticatorData,
		Type:                 credentialType,
	}

	cd, err := clientdata.ParseClientData(abc.RawClientDataJSON)
	if err != nil {
		return inv.Error(err, 400, "failed to parse client data")
	}

	// r1 := protocol.DecodeCredentialAssertionResponse(abc)

	// parsedResponse, err := assertio(*r1)
	// if err != nil {
	// 	return inv.Error(err, 400, "failed to parse attestation")
	// }

	cred := types.NewUnsafeGettableCredential(abc.CredentialID)

	cerem := types.NewUnsafeGettableCeremony(cd.Challenge)

	err = h.Dynamo.TransactGet(ctx, cred, cerem)
	if err != nil {
		return inv.Error(err, 500, "failed to send transact get")
	}

	if cred.RawID.Hex() != cerem.CredentialID.Hex() {
		log.Println(cred.RawID.Hex(), cerem.CredentialID.Hex())
		return inv.Error(nil, 400, "invalid credential id")
	}

	chaner := make(chan *cognitoidentity.GetOpenIdTokenForDeveloperIdentityOutput, 1)
	defer close(chaner)
	stale := false
	defer func() { stale = true }()
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

	// Handle steps 4 through 16
	validError := assertion.VerifyAssertionInput(types.VerifyAssertionInputArgs{
		Input:               abc,
		StoredChallenge:     cerem.ChallengeID,
		RelyingPartyID:      env.RPID(),
		RelyingPartyOrigin:  env.RPOrigin(),
		AttestationType:     "none",
		VerifyUser:          false,
		CredentialPublicKey: cred.PublicKey,
		Extensions:          extensions.ClientInputs{},
	})

	if validError != nil {
		return inv.Error(validError, 400, "failed to verify assertion")
	}

	credentialUpdate, err := cred.UpdateIncreasingCounter(h.Dynamo.MustCredentialTableName())
	if err != nil {
		return inv.Error(err, 500, "failed to create apple pass key")
	}

	err = h.Dynamo.TransactWrite(ctx, *credentialUpdate)
	if err != nil {
		return inv.Error(err, 500, "failed to update apple pass key")
	}

	result := <-chaner
	stale = true

	if result == nil {
		return inv.Error(chanerr, 500, "failed to get dev creds")
	}

	return inv.Success(204, map[string]string{
		"x-nugg-utf-access-token": *result.Token,
	}, "")
}
