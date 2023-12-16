package passkey_attest_test

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/walteh/terrors"
	"github.com/walteh/webauthn/app/passkey_attest"
	"github.com/walteh/webauthn/gen/mockery"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/types"

	"testing"
)

type TestObject struct {
	name              string
	input             *passkey_attest.PasskeyAttestationInput
	want              *passkey_attest.PasskeyAttestationOutput
	existingCeremony  *types.Ceremony
	endingCredentials *types.Credential
	wantErr           bool
}

func TestAttest(t *testing.T) {

	tests := []TestObject{
		{
			name: "A",
			input: &passkey_attest.PasskeyAttestationInput{
				RawAttestationObject: hex.HexToHash("0xa363666d74646e6f6e656761747453746d74a06861757468446174615898a9b9abf7fc46b13564b49d5cf85bcbf371f9cb630e0d6b354bc60b51e065da485d000000000000000000000000000000000000000000147053ed09000cfafdd6e1d98d929796f9c07c466ba501020326200121582030dfb831ebb382bcbd45ac6cb1745222b7d81ad8d44ab33e20d2bda632b5692a225820f6496d03d357717d7669a7af490c8706fef052c0819a02bdca4b92bd42459a00"),
				UTF8ClientDataJSON:   `{"challenge":"pVr2PUG_le6lde9wxeImHA","origin":"https://nugg.xyz","type":"webauthn.create"}`,
				RawCredentialID:      hex.HexToHash("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
				RawSessionID:         hex.HexToHash("0xe12e115acf4552b2568b55e93cbd3939"),
			},
			want: &passkey_attest.PasskeyAttestationOutput{
				AccessToken: "OpenIdToken",
			},
			existingCeremony: &types.Ceremony{
				ChallengeID:  types.CeremonyID(hex.MustBase64ToHash("pVr2PUG_le6lde9wxeImHA")),
				SessionID:    hex.HexToHash("0xe12e115acf4552b2568b55e93cbd3939"),
				CredentialID: types.CredentialID(hex.HexToHash("0x7053ed09000cfafdd6e1d98d929796f9c07c466b")),
				CeremonyType: types.CreateCeremony,
				CreatedAt:    1668984054,
				Ttl:          1668984354,
			},
			// checks: []dtypes.Update{
			// 	{
			// 		Key: map[string]dtypes.AttributeValue{
			// 			"credential_id": types.S("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
			// 		},
			// 		TableName: aws.String("credential"),
			// 		ExpressionAttributeValues: map[string]dtypes.AttributeValue{
			// 			// "created_at":       types.N("1669414368"),
			// 			"session_id":       types.S("0x"),
			// 			"aaguid":           types.S("0x00000000000000000000000000000000"),
			// 			"clone_warning":    types.BOOL(false),
			// 			"public_key":       types.S("0xa501020326200121582030dfb831ebb382bcbd45ac6cb1745222b7d81ad8d44ab33e20d2bda632b5692a225820f6496d03d357717d7669a7af490c8706fef052c0819a02bdca4b92bd42459a00"),
			// 			"attestation_type": types.S("none"),
			// 			"receipt":          types.S("0x"),
			// 			"sign_count":       types.N("0"),
			// 			// "updated_at":       types.N("1669414368"),
			// 			"credential_id":   types.S("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
			// 			"credential_type": types.S("public-key"),
			// 		},
			// 	},
			// 	{
			// 		Key: map[string]dtypes.AttributeValue{
			// 			"challenge_id": types.S(hex.MustBase64ToHash("pVr2PUG_le6lde9wxeImHA").Hex()),
			// 		},
			// 		TableName: aws.String("ceremony"),
			// 		ExpressionAttributeValues: map[string]dtypes.AttributeValue{
			// 			"challenge_id": nil,
			// 		},
			// 	},
			// },
			endingCredentials: &types.Credential{
				// CreatedAt: 1669414368,
				SessionId: hex.HexToHash("0x"),
				// AAGUID:          hex.HexToHash("0x617070617474657374646576656c6f70"),
				AAGUID:          hex.HexToHash("0x00000000000000000000000000000000"),
				CloneWarning:    false,
				PublicKey:       hex.HexToHash("0xa501020326200121582030dfb831ebb382bcbd45ac6cb1745222b7d81ad8d44ab33e20d2bda632b5692a225820f6496d03d357717d7669a7af490c8706fef052c0819a02bdca4b92bd42459a00"),
				AttestationType: string(types.NotFidoAttestationType),
				Receipt:         hex.HexToHash("0x"),
				SignCount:       0,
				// UpdatedAt:       1669414368,
				RawID: types.CredentialID(hex.HexToHash("0x7053ed09000cfafdd6e1d98d929796f9c07c466b")),
				Type:  types.PublicKeyCredentialType,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()

			ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Logger().WithContext(ctx)

			stgp := mockery.NewMockProvider_storage(t)
			rpp := mockery.NewMockProvider_relyingparty(t)
			cpp := mockery.NewMockProvider_accesstoken(t)

			// stgp.EXPECT().GetExisting(ctx, tt.existingCeremony.ChallengeID, types.CredentialID(nil)).Return(tt.existingCeremony, nil, nil)
			stgp.EXPECT().WriteNewCredential(mock.Anything, tt.existingCeremony.ChallengeID, mock.MatchedBy(func(cred *types.Credential) bool {
				// this is just a hack to get a better error message
				return assert.Equal(t, tt.endingCredentials, cred)
			})).Return(nil)

			rpp.EXPECT().RPID().Return("nugg.xyz")
			rpp.EXPECT().RPOrigin().Return("https://nugg.xyz")

			cpp.EXPECT().AccessTokenForUserID(mock.Anything, tt.existingCeremony.SessionID.Hex()).Return("OpenIdToken", nil)

			got, err := passkey_attest.Attest(ctx, stgp, rpp, cpp, tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				if err != nil {
					fmt.Println(terrors.ExtractErrorDetail(err))
					require.NoError(t, err)
				}
			}

			assert.Equal(t, tt.want, got)

			rpp.AssertExpectations(t)
			stgp.AssertExpectations(t)
			cpp.AssertExpectations(t)

		})
	}
}
