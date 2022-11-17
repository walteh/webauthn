package main

import (
	"context"
	"fmt"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/webauthn/protocol"
	"nugg-auth/core/pkg/webauthn/webauthn"

	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"
)

func TestHandler_Invoke(t *testing.T) {

	dynamoClient, cleanup := dynamo.NewMockClient(t)
	defer cleanup()

	wan, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "nugg.xyz",
		RPID:          "nugg.xyz",
		RPOrigin:      "https://nugg.xyz",
	})
	if err != nil {
		t.Fatal(err)
	}

	chal := "20fxslV7kyhXFe3_EvJFgQ"

	dynamo.MockSetCeremony(t, dynamoClient, wan, chal)

	expected := protocol.MockSetRander("ABCD")

	type fields struct {
		Id       string
		Ctx      context.Context
		Dynamo   *dynamo.Client
		Config   config.Config
		Logger   zerolog.Logger
		WebAuthn *webauthn.WebAuthn
		counter  int
	}
	type args struct {
		ctx     context.Context
		payload Input
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Output
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				Id:       "test",
				Ctx:      context.Background(),
				Dynamo:   dynamoClient,
				Config:   nil,
				Logger:   zerolog.New(zerolog.NewTestWriter(t)).With().Caller().Timestamp().Logger(),
				WebAuthn: wan,
				counter:  0,
			},
			args: args{
				ctx: context.Background(),
				payload: Input{
					Headers: map[string]string{
						"Content-Type":                  "application/json",
						"x-nugg-webauthn-attestation":   "o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YViYqbmr9/xGsTVktJ1c+FvL83H5y2MODWs1S8YLUeBl2khdAAAAAAAAAAAAAAAAAAAAAAAAAAAAFCY5NdSMconPGlW0iJMfKJqyKKNfpQECAyYgASFYIFmeyTG8GDlVd28qvhRoKIzb/w7/2A+/+5ATgnztvYmQIlggLk7V2sw9GlIuHdbyGsELL+I36j1dlEdR7IByP68ynHA=",
						"x-nugg-webauthn-clientdata":    fmt.Sprintf("{\"type\":\"webauthn.create\",\"challenge\":\"%s\",\"origin\":\"https://nugg.xyz\"}", chal),
						"x-nugg-webauthn-credential-id": "Jjk11Ixyic8aVbSIkx8omrIoo18=",
					},
				},
			},
			want: Output{
				StatusCode: 204,
				Headers: map[string]string{
					"Content-Length":   "0",
					"x-nugg-challenge": expected,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Id:       tt.fields.Id,
				Ctx:      tt.fields.Ctx,
				Dynamo:   tt.fields.Dynamo,
				Config:   tt.fields.Config,
				Logger:   tt.fields.Logger,
				WebAuthn: tt.fields.WebAuthn,
				counter:  tt.fields.counter,
			}
			got, err := h.Invoke(tt.args.ctx, tt.args.payload)
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
