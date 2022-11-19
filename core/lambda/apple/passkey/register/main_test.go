package main

import (
	"context"
	"encoding/base64"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/webauthn/protocol"
	"nugg-auth/core/pkg/webauthn/webauthn"

	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

func DummyHandler(t *testing.T, chal string, userid string) *Handler {
	dynamoClient := dynamo.NewMockClient(t)

	wan, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "nugg.xyz",
		RPID:          "nugg.xyz",
		RPOrigin:      "https://nugg.xyz",
	})
	if err != nil {
		t.Fatal(err)
	}

	dynamoClient.TransactWrite(context.TODO(), types.TransactWriteItem{
		ConditionCheck: (*types.ConditionCheck)(nil),
		Delete:         (*types.Delete)(nil),
		Put: &types.Put{
			Item: map[string]types.AttributeValue{
				"ceremony_id": &types.AttributeValueMemberS{Value: "pVr2PUG_le6lde9wxeImHA"},
				"session_data": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"challenge": &types.AttributeValueMemberB{
							Value: []uint8{
								0xa5, 0x5a, 0xf6, 0x3d, 0x41, 0xbf, 0x95, 0xee, 0xa5, 0x75, 0xef, 0x70, 0xc5, 0xe2, 0x26, 0x1c,
							},
						},
						"user_id": &types.AttributeValueMemberB{
							Value: []uint8{
								0x84, 0xe0, 0x9e, 0xe7, 0xcf, 0x1d, 0x68, 0x96, 0x6e, 0x18, 0x0e, 0x8f, 0x78, 0x94, 0xdc, 0xcd,
							},
						},
						"user_verification": &types.AttributeValueMemberS{Value: ""}}},
				"ttl": &types.AttributeValueMemberN{Value: "1669028360"},
			},
			TableName: &dynamoClient.CeremonyTableName,
		},
	})

	return &Handler{
		Id:       "test",
		Ctx:      context.Background(),
		Dynamo:   dynamoClient,
		Config:   nil,
		Logger:   zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Timestamp().Logger(),
		WebAuthn: wan,
		counter:  0,
	}
}

func TestHandler_Invoke(t *testing.T) {

	chal := "pVr2PUG_le6lde9wxeImHA"

	Handler := DummyHandler(t, chal, "hOCe588daJZuGA6PeJTczQ")

	rander := protocol.MockSetRander(t, "ABCD")

	tests := []struct {
		name    string
		args    Input
		want    Output
		wantErr bool
	}{
		{
			name: "test",
			args: Input{
				Headers: map[string]string{
					"Content-Type":                      "application/json",
					"x-nugg-apple-passkey-clientdata":   "{\"type\":\"webauthn.create\",\"challenge\":\"pVr2PUG_le6lde9wxeImHA\",\"origin\":\"https://nugg.xyz\"}",
					"x-nugg-apple-passkey-attestation":  "o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YViYqbmr9_xGsTVktJ1c-FvL83H5y2MODWs1S8YLUeBl2khdAAAAAAAAAAAAAAAAAAAAAAAAAAAAFHBT7QkADPr91uHZjZKXlvnAfEZrpQECAyYgASFYIDDfuDHrs4K8vUWsbLF0UiK32BrY1EqzPiDSvaYytWkqIlgg9kltA9NXcX12aaevSQyHBv7wUsCBmgK9ykuSvUJFmgA",
					"x-nugg-apple-passkey-credentialid": "cFPtCQAM-v3W4dmNkpeW-cB8Rms",
					// {"rawClientDataJSON":"eyJ0eXBlIjoid2ViYXV0aG4uY3JlYXRlIiwiY2hhbGxlbmdlIjoicFZyMlBVR19sZTZsZGU5d3hlSW1IQSIsIm9yaWdpbiI6Imh0dHBzOi8vbnVnZy54eXoifQ","rawAttestationObject":"o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YViYqbmr9_xGsTVktJ1c-FvL83H5y2MODWs1S8YLUeBl2khdAAAAAAAAAAAAAAAAAAAAAAAAAAAAFHBT7QkADPr91uHZjZKXlvnAfEZrpQECAyYgASFYIDDfuDHrs4K8vUWsbLF0UiK32BrY1EqzPiDSvaYytWkqIlgg9kltA9NXcX12aaevSQyHBv7wUsCBmgK9ykuSvUJFmgA","credentialID":"cFPtCQAM-v3W4dmNkpeW-cB8Rms"
					// "x-nugg-payload": "eyJyYXdDbGllbnREYXRhSlNPTiI6ImV5SjBlWEJsSWpvaWQyVmlZWFYwYUc0dVkzSmxZWFJsSWl3aVkyaGhiR3hsYm1kbElqb2ljRlp5TWxCVlIxOXNaVFpzWkdVNWQzaGxTVzFJUVNJc0ltOXlhV2RwYmlJNkltaDBkSEJ6T2k4dmJuVm5aeTU0ZVhvaWZRIiwicmF3QXR0ZXN0YXRpb25PYmplY3QiOiJvMk5tYlhSa2JtOXVaV2RoZEhSVGRHMTBvR2hoZFhSb1JHRjBZVmlZcWJtcjlfeEdzVFZrdEoxYy1Gdkw4M0g1eTJNT0RXczFTOFlMVWVCbDJraGRBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBRkhCVDdRa0FEUHI5MXVIWmpaS1hsdm5BZkVacnBRRUNBeVlnQVNGWUlERGZ1REhyczRLOHZVV3NiTEYwVWlLMzJCclkxRXF6UGlEU3ZhWXl0V2txSWxnZzlrbHRBOU5YY1gxMmFhZXZTUXlIQnY3d1VzQ0JtZ0s5eWt1U3ZVSkZtZ0EiLCJjcmVkZW50aWFsSUQiOiJjRlB0Q1FBTS12M1c0ZG1Oa3BlVy1jQjhSbXMifQ==",
				},
			},
			want: Output{
				StatusCode: 204,
				Headers: map[string]string{
					"Content-Length":   "0",
					"x-nugg-challenge": base64.RawURLEncoding.EncodeToString(rander.CalculateDeterministicHash(1)),
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

	t.Run("count check", func(t *testing.T) {
		if Handler.counter != len(tests) {
			t.Errorf("Handler.Invoke() counter = %v, want 1", Handler.counter)
		}

	})
}
