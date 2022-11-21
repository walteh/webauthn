package main

import (
	"context"
	"nugg-webauthn/core/pkg/dynamo"
	"nugg-webauthn/core/pkg/webauthn/protocol"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

func DummyHandler(t *testing.T) *Handler {
	dynamoClient := dynamo.NewMockClient(t)

	return &Handler{
		Id:      "test",
		Ctx:     context.Background(),
		Dynamo:  dynamoClient,
		Config:  nil,
		Logger:  zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Timestamp().Logger(),
		counter: 0,
	}
}

func TestHandler_Invoke_UnitTest1234(t *testing.T) {

	expected := protocol.MockSetRander(t, "xsTWpSak5HWm")

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
					"Content-Type":          "application/json",
					"x-nugg-hex-session-id": "0xff33ff",
				},
			},
			want: Output{
				StatusCode: 204,
				Headers: map[string]string{
					"Content-Length":       "0",
					"x-nugg-hex-challenge": expected.CalculateDeterministicHash(1).Hex(),
				},
			},
			wantErr: false,
		},
		{
			name: "B",
			args: Input{
				Headers: map[string]string{
					"Content-Type":             "application/json",
					"x-nugg-hex-session-id":    "0xff33ff",
					"x-nugg-hex-credential-id": "0x7053ed09000cfafdd6e1d98d929796f9c07c466b",
				},
			},
			want: Output{
				StatusCode: 204,
				Headers: map[string]string{
					"Content-Length":       "0",
					"x-nugg-hex-challenge": expected.CalculateDeterministicHash(1).Hex(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			expected = protocol.MockSetRander(t, "xsTWpSak5HWm")

			Handler := DummyHandler(t)

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
