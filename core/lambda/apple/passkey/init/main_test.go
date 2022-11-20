package main

import (
	"context"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/webauthn/protocol"
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

	Handler := DummyHandler(t)

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
					"Content-Type":              "application/json",
					"x-nugg-webauthn-sessionId": "0xff33ff",
				},
			},
			want: Output{
				StatusCode: 204,
				Headers: map[string]string{
					"Content-Length":            "0",
					"x-nugg-webauthn-challenge": expected.CalculateDeterministicHash(1).Hex(),
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
