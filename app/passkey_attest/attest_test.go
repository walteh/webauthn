package passkey_test

import (
	dynamodb_mock "git.nugg.xyz/go-sdk/dynamo/mock"
	"git.nugg.xyz/pkg/invocation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/walteh/webauthn/pkg/cognito"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/passkey"
	"github.com/walteh/webauthn/pkg/webauthn/types"

	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestAttest(t *testing.T) {

	tests := []struct {
		name    string
		input   passkey.PasskeyAttestationInput
		want    passkey.PasskeyAttestationOutput
		puts    []dtypes.Put
		checks  []dtypes.Update
		wantErr bool
	}{
		{
			name: "A",
			input: passkey.PasskeyAttestationInput{
				RawAttestationObject: hex.HexToHash("0xa363666d74646e6f6e656761747453746d74a06861757468446174615898a9b9abf7fc46b13564b49d5cf85bcbf371f9cb630e0d6b354bc60b51e065da485d000000000000000000000000000000000000000000147053ed09000cfafdd6e1d98d929796f9c07c466ba501020326200121582030dfb831ebb382bcbd45ac6cb1745222b7d81ad8d44ab33e20d2bda632b5692a225820f6496d03d357717d7669a7af490c8706fef052c0819a02bdca4b92bd42459a00"),
				UTF8ClientDataJSON:   `{"challenge":"pVr2PUG_le6lde9wxeImHA","origin":"https://nugg.xyz","type":"webauthn.create"}`,
				RawCredentialID:      hex.HexToHash("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
			},
			want: passkey.PasskeyAttestationOutput{
				SuggestedStatusCode: 204,
				AccessToken:         "OpenIdToken",
			},
			puts: []dtypes.Put{
				{
					Item: map[string]dtypes.AttributeValue{
						"challenge_id":  types.S(hex.MustBase64ToHash("pVr2PUG_le6lde9wxeImHA").Hex()),
						"session_id":    types.S("0xe12e115acf4552b2568b55e93cbd3939"),
						"credential_id": types.S("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
						"ceremony_type": types.S("webauthn.create"),
						"created_at":    types.N("1668984054"),
						"ttl":           types.N("1668984354"),
					},
					TableName: aws.String("ceremony"),
				},
			},
			checks: []dtypes.Update{
				{
					Key: map[string]dtypes.AttributeValue{
						"credential_id": types.S("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
					},
					TableName: aws.String("credential"),
					ExpressionAttributeValues: map[string]dtypes.AttributeValue{
						// "created_at":       types.N("1669414368"),
						"session_id":       types.S("0x"),
						"aaguid":           types.S("0x00000000000000000000000000000000"),
						"clone_warning":    types.BOOL(false),
						"public_key":       types.S("0xa501020326200121582030dfb831ebb382bcbd45ac6cb1745222b7d81ad8d44ab33e20d2bda632b5692a225820f6496d03d357717d7669a7af490c8706fef052c0819a02bdca4b92bd42459a00"),
						"attestation_type": types.S("none"),
						"receipt":          types.S("0x"),
						"sign_count":       types.N("0"),
						// "updated_at":       types.N("1669414368"),
						"credential_id":   types.S("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
						"credential_type": types.S("public-key"),
					},
				},
				{
					Key: map[string]dtypes.AttributeValue{
						"challenge_id": types.S(hex.MustBase64ToHash("pVr2PUG_le6lde9wxeImHA").Hex()),
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

			got, err := passkey.Attest(h.Ctx(), api, cog, tt.input)
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
