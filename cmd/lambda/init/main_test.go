package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/walteh/webauthn/gen/buf/go/proto/webauthn/v1"
	"github.com/walteh/webauthn/gen/mockery"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/challenge"
	"github.com/walteh/webauthn/pkg/webauthn/types"

	"github.com/rs/zerolog"
)

func DummyHandler(t *testing.T) *Handler {

	return &Handler{
		Id:      "test",
		Ctx:     context.Background(),
		Config:  nil,
		logger:  zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Timestamp().Logger(),
		counter: 0,
	}
}

func TestHandler_Invoke_UnitTest1234(t *testing.T) {

	expected := challenge.MockSetRander(t, "xsTWpSak5HWm")

	tests := []struct {
		name    string
		args    *InputBody
		want    *OutputBody
		wantErr bool
	}{
		{
			name: "A",
			args: &InputBody{
				Msg: &webauthn.CreateChallengeRequest{
					SessionId: hex.HexToHash("0xff33ff"),
				},
			},
			want: &OutputBody{
				Msg: &webauthn.CreateChallengeResponse{
					Challenge: expected.CalculateDeterministicHash(1),
					Ttl:       180000,
				},
			},
			wantErr: false,
		},
		{
			name: "B",
			args: &InputBody{
				Msg: &webauthn.CreateChallengeRequest{
					SessionId:    hex.HexToHash("0xff33ff"),
					CredentialId: hex.HexToHash("0x7053ed09000cfafdd6e1d98d929796f9c07c466b"),
				},
			},
			want: &OutputBody{
				Msg: &webauthn.CreateChallengeResponse{
					Challenge: expected.CalculateDeterministicHash(1),
					Ttl:       180000,
				},
			},

			wantErr: false,
		},
		{
			name: "C",
			args: &InputBody{
				Msg: &webauthn.CreateChallengeRequest{
					SessionId:    hex.HexToHash("0xff33ff"),
					CeremonyType: "webauthn.not-webauthn",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			expected = challenge.MockSetRander(t, "xsTWpSak5HWm")

			Handler := DummyHandler(t)

			dynamoClient := mockery.NewMockProvider_storage(t)

			Handler.Storage = dynamoClient

			dynamoClient.EXPECT().WriteNewCeremony(mock.Anything, mock.MatchedBy(func(cer *types.Ceremony) bool {
				assert.Equal(t, tt.args.Msg.GetCredentialId(), cer.CredentialID.Ref().Bytes())
				assert.Equal(t, tt.args.Msg.GetSessionId(), cer.SessionID.Bytes())
				if tt.args.Msg.GetCeremonyType() == "" {
					assert.Equal(t, string(types.AssertCeremony), string(cer.CeremonyType))
				} else {
					assert.Equal(t, tt.args.Msg.GetCeremonyType(), string(cer.CeremonyType))
				}
				return true
			})).Return(nil).Maybe()

			got, err := Handler.CreateChallenge(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.Invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)

			dynamoClient.AssertExpectations(t)

		})
	}
}
