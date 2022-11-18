package main

import (
	"context"
	"fmt"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/webauthn/protocol"
	"nugg-auth/core/pkg/webauthn/webauthn"

	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

func DummyHandler(t *testing.T, chal string) *Handler {
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
				"ceremony_id": &types.AttributeValueMemberS{Value: chal},
				"session_data": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"challenge": &types.AttributeValueMemberS{Value: chal},
						"user_id": &types.AttributeValueMemberB{
							Value: []uint8{
								0x45, 0x45, 0x45}},
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
		Logger:   zerolog.New(zerolog.NewTestWriter(t)).With().Caller().Timestamp().Logger(),
		WebAuthn: wan,
		counter:  0,
	}
}

func TestHandler_Invoke(t *testing.T) {

	chal := "xsTWpSak5HWm"

	Handler := DummyHandler(t, chal)

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
					"Content-Type":                  "application/json",
					"x-nugg-webauthn-attestation":   "o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YViYqbmr9/xGsTVktJ1c+FvL83H5y2MODWs1S8YLUeBl2khdAAAAAAAAAAAAAAAAAAAAAAAAAAAAFDl1hTTYx9TPs0uNL9LRyHZ72fDjpQECAyYgASFYIOvGCF1LRLSbI+58Wx7AQIGH2MKBPJvrA5lTDG/yqKbEIlggrrAu3x94Y7zBa8DJjwXIIUZ1/0bDWqpGh7BkF1ZrACU=",
					"x-nugg-webauthn-clientdata":    fmt.Sprintf("{\"type\":\"webauthn.create\",\"challenge\":\"%s\",\"origin\":\"https://nugg.xyz\"}", chal),
					"x-nugg-webauthn-credential-id": "OXWFNNjH1M+zS40v0tHIdnvZ8OM=",
				},
			},
			want: Output{
				StatusCode: 204,
				Headers: map[string]string{
					"Content-Length":   "0",
					"x-nugg-challenge": rander.CalculateDeterministicHash(1),
					"x-nugg-user-id":   rander.CalculateDeterministicHash(2),
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
