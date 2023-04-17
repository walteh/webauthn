package passkey

import (
	"context"

	"git.nugg.xyz/webauthn/pkg/cognito"
	"git.nugg.xyz/webauthn/pkg/dynamo"
	"git.nugg.xyz/webauthn/pkg/hex"
	"git.nugg.xyz/webauthn/pkg/webauthn/types"

	"reflect"
	"testing"

	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestAttest(t *testing.T) {

	dynamo.AttachLocalDynamoServer(t)

	tests := []struct {
		name    string
		input   PasskeyAttestationInput
		want    PasskeyAttestationOutput
		puts    []dtypes.Put
		checks  []dtypes.Update
		wantErr bool
	}{
		{
			name: "A",
			input: PasskeyAttestationInput{
				RawAttestationObject: hex.HexToHash("0xa363666d74646e6f6e656761747453746d74a06861757468446174615898a9b9abf7fc46b13564b49d5cf85bcbf371f9cb630e0d6b354bc60b51e065da485d000000000000000000000000000000000000000000147053ed09000cfafdd6e1d98d929796f9c07c466ba501020326200121582030dfb831ebb382bcbd45ac6cb1745222b7d81ad8d44ab33e20d2bda632b5692a225820f6496d03d357717d7669a7af490c8706fef052c0819a02bdca4b92bd42459a00"),
				UTF8ClientDataJSON:   `{"challenge":"pVr2PUG_le6lde9wxeImHA","origin":"https://nugg.xyz","type":"webauthn.create"}`,
				RawCredentialID:      hex.HexToHash("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
			},
			want: PasskeyAttestationOutput{
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
					TableName: dynamo.MockCeremonyTableName(),
				},
			},
			checks: []dtypes.Update{
				{
					Key: map[string]dtypes.AttributeValue{
						"credential_id": types.S("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
					},
					TableName: dynamo.MockCredentialTableName(),
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
					TableName: dynamo.MockCeremonyTableName(),
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

			client := dynamo.NewMockClient(t)

			dynamo.MockBatchPut(t, client, tt.puts...)

			cog := cognito.NewMockClient()

			got, err := Attest(context.Background(), client, cog, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.Invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.Invoke() = %v, want %v", got, tt.want)
			}

			dynamo.MockBatchCheck(t, client, tt.checks...)
		})
	}
}
