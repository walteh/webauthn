package passkey_assert_test

import (
	"github.com/walteh/webauthn/app/passkey_assert"
	"github.com/walteh/webauthn/pkg/accesstoken/cognito"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/types"

	"github.com/walteh/buildrc/integration/dynamodb"

	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_Invoke(t *testing.T) {

	tests := []struct {
		name    string
		input   passkey_assert.PasskeyAssertionInput
		want    passkey_assert.PasskeyAssertionOutput
		puts    []dtypes.Put
		checks  []dtypes.Update
		wantErr bool
	}{
		{
			name: "A",
			input: passkey_assert.PasskeyAssertionInput{
				SessionID:            hex.MustBase64ToHash("hOCe588daJZuGA6PeJTczQ"),
				CredentialID:         hex.HexToHash("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
				UTF8ClientDataJSON:   "{\"type\":\"webauthn.get\",\"challenge\":\"4S4RWs9FUrJWi1XpPL05OQ\",\"origin\":\"https://nugg.xyz\"}",
				RawAuthenticatorData: hex.HexToHash("0xa9b9abf7fc46b13564b49d5cf85bcbf371f9cb630e0d6b354bc60b51e065da481d00000000"),
				RawSignature:         hex.MustBase64ToHash("MEUCIA8DhDsdF8XcwD6T9X0R1C68oeFw-gNgy1lHMYGxi_WHAiEA6-JXpbLdY39d6fK9oDRDpLtDAv7DplSl7p-Nm_NiFJc"),
			},
			want: passkey_assert.PasskeyAssertionOutput{
				SuggestedStatusCode: 204,
				AccessToken:         "OpenIdToken",
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
					TableName: aws.String("ceremony"),
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
					TableName: aws.String("credential"),
				},
			},
			checks: []dtypes.Update{
				{
					Key: map[string]dtypes.AttributeValue{
						"credential_id": types.S("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
					},
					TableName: aws.String("credential"),
					ExpressionAttributeValues: map[string]dtypes.AttributeValue{
						"sign_count": types.N("1"),
					},
				},
				{
					Key: map[string]dtypes.AttributeValue{
						"challenge_id": types.S("0xe12e115acf4552b2568b55e93cbd3939"),
					},
					TableName: aws.String("ceremony"),
					ExpressionAttributeValues: map[string]dtypes.AttributeValue{
						"challenge_id": nil,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := invocation.NewDefaultTestHandler[string, string]()

			_, api := dynamodb_mock.NewMockAPIFromTerraform(t, h.Ctx(), *h.Opts().AwsConfig(), "../../terraform/dynamodb.tf")

			for _, table := range tt.puts {
				_, err := api.PutItem(h.Ctx(), &dynamodb.PutItemInput{
					Item:      table.Item,
					TableName: table.TableName,
				})
				require.NoError(t, err)
			}

			cog := cognito.NewMockClient()

			got, err := passkey_assert.Assert(h.Ctx(), api, cog, tt.input)
			require.NoError(t, err)

			assert.Equal(t, tt.want, got)

			for _, table := range tt.checks {
				r, err := api.GetItem(h.Ctx(), &dynamodb.GetItemInput{
					Key:       table.Key,
					TableName: table.TableName,
				})
				require.NoError(t, err)

				assert.Equal(t, table.ExpressionAttributeValues, r.Item)
			}
		})
	}
}
