package main

import (
	"context"
	"encoding/base64"
	"nugg-auth/core/pkg/cognito"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/webauthn/protocol"

	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

func DummyHandler(t *testing.T) *Handler {
	dynamoClient := dynamo.NewMockClient(t)

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
		Id:      "test",
		Ctx:     context.Background(),
		Dynamo:  dynamoClient,
		Config:  nil,
		Cognito: cognito.NewMockClient(),
		Logger:  zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Timestamp().Logger(),
		counter: 0,
	}
}

func TestHandler_Invoke(t *testing.T) {

	Handler := DummyHandler(t)

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
					"Content-Type":             "application/json",
					"x-nugg-webauthn-creation": "eyJyYXdDbGllbnREYXRhSlNPTiI6ImV5SmphR0ZzYkdWdVoyVWlPaUp3Vm5JeVVGVkhYMnhsTm14a1pUbDNlR1ZKYlVoQklpd2liM0pwWjJsdUlqb2lhSFIwY0hNNkx5OXVkV2RuTG5oNWVpSXNJblI1Y0dVaU9pSjNaV0poZFhSb2JpNWpjbVZoZEdVaWZRPT0iLCJyYXdBdHRlc3RhdGlvbk9iamVjdCI6Im8yTm1iWFJrYm05dVpXZGhkSFJUZEcxMG9HaGhkWFJvUkdGMFlWaVlxYm1yOS94R3NUVmt0SjFjK0Z2TDgzSDV5Mk1PRFdzMVM4WUxVZUJsMmtoZEFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFGSEJUN1FrQURQcjkxdUhaalpLWGx2bkFmRVpycFFFQ0F5WWdBU0ZZSUREZnVESHJzNEs4dlVXc2JMRjBVaUszMkJyWTFFcXpQaURTdmFZeXRXa3FJbGdnOWtsdEE5TlhjWDEyYWFldlNReUhCdjd3VXNDQm1nSzl5a3VTdlVKRm1nQT0iLCJjcmVkZW50aWFsSWQiOiJjRlB0Q1FBTSt2M1c0ZG1Oa3BlVytjQjhSbXM9In0=",
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

			r, err := Handler.Dynamo.GetItem(Handler.Ctx, &dynamodb.GetItemInput{
				TableName: &Handler.Dynamo.CeremonyTableName,
				Key: map[string]types.AttributeValue{
					"ceremony_id": &types.AttributeValueMemberS{Value: "pVr2PUG_le6lde9wxeImHA"},
				},
			})
			if err != nil {
				t.Error(err)
			}

			if r.Item == nil {
				t.Error("item is nil")
			}

			c, err := Handler.Dynamo.GetItem(Handler.Ctx, &dynamodb.GetItemInput{
				TableName: &Handler.Dynamo.CredentialTableName,
				Key: map[string]types.AttributeValue{
					"credential_id": &types.AttributeValueMemberS{Value: "cFPtCQAM-v3W4dmNkpeW-cB8Rms"},
				},
			})

			if err != nil {
				t.Error(err)
			}

			if c.Item == nil {
				t.Error("item is nil")
			}

		})
	}

	t.Run("count check", func(t *testing.T) {
		if Handler.counter != len(tests) {
			t.Errorf("Handler.Invoke() counter = %v, want 1", Handler.counter)
		}

	})
}
