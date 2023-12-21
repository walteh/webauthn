package passkey_assert_test

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/walteh/terrors"
	"github.com/walteh/webauthn/app/passkey_assert"
	"github.com/walteh/webauthn/gen/mockery"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/types"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestObject struct {
	name                string
	input               *passkey_assert.PasskeyAssertionInput
	want                *passkey_assert.PasskeyAssertionOutput
	existingCeremony    *types.Ceremony
	existingCredentials *types.Credential
	wantErr             bool
}

func TestHandler_Invoke(t *testing.T) {

	tests := []TestObject{
		{
			name: "A",
			input: &passkey_assert.PasskeyAssertionInput{
				SessionID:            hex.MustBase64ToHash("hOCe588daJZuGA6PeJTczQ"),
				CredentialID:         hex.HexToHash("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
				UTF8ClientDataJSON:   "{\"type\":\"webauthn.get\",\"challenge\":\"4S4RWs9FUrJWi1XpPL05OQ\",\"origin\":\"https://nugg.xyz\"}",
				RawAuthenticatorData: hex.HexToHash("0xa9b9abf7fc46b13564b49d5cf85bcbf371f9cb630e0d6b354bc60b51e065da481d00000000"),
				RawSignature:         hex.MustBase64ToHash("MEUCIA8DhDsdF8XcwD6T9X0R1C68oeFw-gNgy1lHMYGxi_WHAiEA6-JXpbLdY39d6fK9oDRDpLtDAv7DplSl7p-Nm_NiFJc"),
			},
			want: &passkey_assert.PasskeyAssertionOutput{
				AccessToken: "OpenIdToken",
			},
			existingCredentials: &types.Credential{
				CreatedAt: 1669414368,
				SessionId: hex.HexToHash("0x"),
				// AAGUID:          hex.HexToHash("0x617070617474657374646576656c6f70"),
				AAGUID:          hex.HexToHash("0x00000000000000000000000000000000"), // none attestation provider
				CloneWarning:    false,
				PublicKey:       hex.HexToHash("0xa501020326200121582030dfb831ebb382bcbd45ac6cb1745222b7d81ad8d44ab33e20d2bda632b5692a225820f6496d03d357717d7669a7af490c8706fef052c0819a02bdca4b92bd42459a00"),
				AttestationType: string(types.NotFidoAttestationType),
				Receipt:         hex.HexToHash("0x"),
				SignCount:       0,
				UpdatedAt:       1669414368,
				RawID:           types.CredentialID(hex.HexToHash("0x7053ed09000cfafdd6e1d98d929796f9c07c466b")),
				Type:            types.PublicKeyCredentialType,
			},

			existingCeremony: &types.Ceremony{
				ChallengeID:  types.CeremonyID(hex.HexToHash("0xe12e115acf4552b2568b55e93cbd3939")),
				SessionID:    hex.HexToHash("0xe12e115acf4552b2568b55e93cbd3939"),
				CredentialID: types.CredentialID(hex.HexToHash("0x7053ed09000cfafdd6e1d98d929796f9c07c466b")),
				CeremonyType: types.CreateCeremony,
				CreatedAt:    1668984054,
				Ttl:          1668984354,
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

			stgp.EXPECT().GetExisting(mock.Anything, tt.existingCeremony.ChallengeID, tt.existingCredentials.RawID).Return(tt.existingCeremony, tt.existingCredentials, nil)
			stgp.EXPECT().IncrementExistingCredential(mock.Anything, tt.existingCeremony.ChallengeID, tt.existingCredentials).Return(nil).Maybe()

			rpp.EXPECT().RPID().Return("nugg.xyz")
			rpp.EXPECT().RPOrigin().Return("https://nugg.xyz")

			cpp.EXPECT().AccessTokenForUserID(mock.Anything, tt.input.SessionID.Hex()).Return("OpenIdToken", nil)

			got, err := passkey_assert.Assert(ctx, stgp, rpp, cpp, tt.input)
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
