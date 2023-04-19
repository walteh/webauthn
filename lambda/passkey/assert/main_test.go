package main

import (
	"context"
	"fmt"
	"os"

	"git.nugg.xyz/webauthn/pkg/cognito"
	"git.nugg.xyz/webauthn/pkg/hex"
	"git.nugg.xyz/webauthn/pkg/webauthn/types"

	"reflect"
	"testing"

	dynamodb_mock "git.nugg.xyz/go-sdk/dynamo/mock"
	"git.nugg.xyz/go-sdk/invocation"
	"git.nugg.xyz/go-sdk/mock"
	"git.nugg.xyz/go-sdk/x"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

func TestMain(m *testing.M) {
	code := mock.ContainerTestMain(func() int {
		return m.Run()
	}) //

	os.Exit(code)
}

func DummyHandler(t *testing.T) *Handler[x.DynamoDBAPIProvisioner] {

	h := invocation.NewDefaultTestHandler[Input, Output]()

	_, dynamoClient := dynamodb_mock.NewMockAPIFromTerraform(t, h.Ctx(), *h.Opts().AwsConfig(), "../../../terraform/dynamo.tf")

	prov, _ := buildHandler(h, dynamoClient, cognito.NewMockClient())

	return prov
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

	tests := []struct {
		name    string
		args    Input
		want    Output
		puts    []dtypes.Put
		wantErr bool
	}{
		{
			name: "A",
			args: Input{
				Headers: map[string]string{
					"x-nugg-hex-authenticator-data": "0xa9b9abf7fc46b13564b49d5cf85bcbf371f9cb630e0d6b354bc60b51e065da481d00000000",
					"x-nugg-hex-credential-id":      "0x7053ed09000cfafdd6e1d98d929796f9c07c466b",
					"x-nugg-hex-signature":          hex.MustBase64ToHash("MEUCIA8DhDsdF8XcwD6T9X0R1C68oeFw-gNgy1lHMYGxi_WHAiEA6-JXpbLdY39d6fK9oDRDpLtDAv7DplSl7p-Nm_NiFJc").Hex(),
					"x-nugg-hex-user-id":            hex.MustBase64ToHash("hOCe588daJZuGA6PeJTczQ").Hex(),
					"x-nugg-utf-client-data-json":   "{\"type\":\"webauthn.get\",\"challenge\":\"4S4RWs9FUrJWi1XpPL05OQ\",\"origin\":\"https://nugg.xyz\"}",
					"x-nugg-utf-credential-type":    "public-key"},
			},
			want: Output{
				StatusCode: 204,
				Headers: map[string]string{
					"Content-Length":          "0",
					"x-nugg-utf-access-token": "OpenIdToken",
				},
			},
			puts: []dtypes.Put{
				{
					Item: map[string]dtypes.AttributeValue{
						"challenge_id":  types.S("0xe12e115acf4552b2568b55e93cbd3939"),
						"session_id":    types.S("0xe12e115acf4552b2568b55e93cbd3939"),
						"credential_id": types.S("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
						"ceremony_type": types.S("webauthn.create"),
						"created_at":    types.N("1668984054"),
						"ttl":           types.N("1668984354"),
					},
				},
				{
					Item: map[string]dtypes.AttributeValue{
						"created_at":       types.N("1669414368"),
						"session_id":       types.S("0x"),
						"aaguid":           types.S("0x617070617474657374646576656c6f70"),
						"clone_warning":    types.BOOL(false),
						"public_key":       types.S("0xa501020326200121582030dfb831ebb382bcbd45ac6cb1745222b7d81ad8d44ab33e20d2bda632b5692a225820f6496d03d357717d7669a7af490c8706fef052c0819a02bdca4b92bd42459a00"),
						"attestation_type": types.S("none"),
						"receipt":          types.S("0x"),
						"sign_count":       types.N("0"),
						"updated_at":       types.N("1669414368"),
						"credential_id":    types.S("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
						"credential_type":  types.S("public-key"),
					},
				},
			},

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hnd := DummyHandler(t)

			for _, put := range tt.puts {
				_, err := hnd.Dynamo.PutItem(hnd.Ctx(), &dynamodb.PutItemInput{
					TableName: aws.String("nugg-credentials"),
					Item:      put.Item,
				})
				if err != nil {
					t.Errorf("Handler.Invoke() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			got, err := hnd.Invoke(context.Background(), tt.args)
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
