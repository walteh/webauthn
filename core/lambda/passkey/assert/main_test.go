package main

import (
	"context"
	"fmt"
	"nugg-webauthn/core/pkg/cognito"
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/hex"

	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

func DummyHandler(t *testing.T) *Handler {
	dynamoClient := dynamo.NewMockClient(t)

	err := dynamoClient.TransactWrite(context.Background(),
		types.TransactWriteItem{
			Put: &types.Put{
				Item: map[string]types.AttributeValue{
					"user_id": &types.AttributeValueMemberS{
						Value: "01GJ64VW4WB50CYYN2KDN94WKH",
					},
					"created_at": &types.AttributeValueMemberS{
						Value: "2022-11-18 07:22:54.668188 -0600 CST m=+0.032937709",
					},
					"updated_at": &types.AttributeValueMemberS{
						Value: "2022-11-18 07:22:54.668281 -0600 CST m=+0.033030751",
					},
				},
				TableName: &dynamoClient.UserTableName,
			},
		},
		types.TransactWriteItem{
			Put: &types.Put{
				Item: map[string]types.AttributeValue{
					"receipt": &types.AttributeValueMemberS{
						Value: "0x",
					},
					"aaguid": &types.AttributeValueMemberS{
						Value: "0x",
					},
					"created_at": &types.AttributeValueMemberN{
						Value: "1669001824",
					},
					"credential_id": &types.AttributeValueMemberS{
						Value: "0x7053ed09000cfafdd6e1d98d929796f9c07c466b",
					},
					"credential_type": &types.AttributeValueMemberS{
						Value: "public-key",
					},
					"public_key": &types.AttributeValueMemberS{
						Value: "0xa501020326200121582030dfb831ebb382bcbd45ac6cb1745222b7d81ad8d44ab33e20d2bda632b5692a225820f6496d03d357717d7669a7af490c8706fef052c0819a02bdca4b92bd42459a00",
					},
					"updated_at": &types.AttributeValueMemberN{
						Value: "1669001824",
					},
					"session_id": &types.AttributeValueMemberS{
						Value: "0x",
					},
					"attestation_type": &types.AttributeValueMemberS{
						Value: "none",
					},
					"sign_count": &types.AttributeValueMemberN{
						Value: "0",
					},
					"clone_warning": &types.AttributeValueMemberBOOL{
						Value: false,
					},
				},
				TableName: &dynamoClient.CredentialTableName,
			},
		}, types.TransactWriteItem{
			Put: &types.Put{
				Item: map[string]types.AttributeValue{
					"challenge_id": &types.AttributeValueMemberS{
						Value: "0xe12e115acf4552b2568b55e93cbd3939",
					},
					"session_id": &types.AttributeValueMemberS{
						Value: "0xe12e115acf4552b2568b55e93cbd3939",
					},
					"credential_id": &types.AttributeValueMemberS{
						Value: "0x7053ed09000cfafdd6e1d98d929796f9c07c466b",
					},
					"ceremony_type": &types.AttributeValueMemberS{
						Value: "",
					},
					"created_at": &types.AttributeValueMemberN{
						Value: "0",
					},
					"ttl": &types.AttributeValueMemberN{
						Value: "0",
					},
				},
				TableName: &dynamoClient.CeremonyTableName,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	return &Handler{
		Id:      "test",
		Ctx:     context.Background(),
		Dynamo:  dynamoClient,
		Config:  nil,
		Cognito: cognito.NewMockClient(),
		logger:  zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Timestamp().Logger(),
		counter: 0,
	}
}

var arg = fmt.Sprintf(
	`{"rawAuthenticatorData":"%s","signature":"%s","userID":"%s","rawClientDataJSON":"%s","credentialID":"%s","credentialType":"public-key"}`,
	hex.MustBase64ToHash("qbmr9_xGsTVktJ1c-FvL83H5y2MODWs1S8YLUeBl2kgdAAAAAA").Hex(),
	hex.MustBase64ToHash("MEUCIA8DhDsdF8XcwD6T9X0R1C68oeFw-gNgy1lHMYGxi_WHAiEA6-JXpbLdY39d6fK9oDRDpLtDAv7DplSl7p-Nm_NiFJc").Hex(),
	hex.MustBase64ToHash("hOCe588daJZuGA6PeJTczQ").Hex(),
	`{\"type\":\"webauthn.get\",\"challenge\":\"4S4RWs9FUrJWi1XpPL05OQ\",\"origin\":\"https://nugg.xyz\"}`,
	hex.MustBase64ToHash("cFPtCQAM-v3W4dmNkpeW-cB8Rms").Hex(),
)

func TestHandler_Invoke(t *testing.T) {

	Handler := DummyHandler(t)

	tests := []struct {
		name    string
		args    Input
		want    Output
		wantErr bool
	}{
		{
			name: "A",
			args: Input{
				Headers: map[string]string{
					"x-nugg-hex-assertion": hex.Hash(arg).Hex(),
				},
			},
			want: Output{
				StatusCode: 204,
				Headers: map[string]string{
					"Content-Length":          "0",
					"x-nugg-utf-access-token": "OpenIdToken",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := Handler.Invoke(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.Invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.Invoke() = %v, want %v", got, tt.want)
			}
		})
	}
}
