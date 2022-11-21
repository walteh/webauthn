package main

import (
	"context"
	"nugg-auth/core/pkg/cognito"
	"nugg-auth/core/pkg/dynamo"

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
				"challenge_id": &types.AttributeValueMemberS{
					Value: "0xa55af63d41bf95eea575ef70c5e2261c",
				},
				"session_id": &types.AttributeValueMemberS{
					Value: "0xa55af63d41bf95eea575ef70c5e2261c",
				},
				"credential_id": &types.AttributeValueMemberS{
					Value: "0x",
				},
				"ceremony_type": &types.AttributeValueMemberS{
					Value: "webauthn.create",
				},
				"created_at": &types.AttributeValueMemberN{
					Value: "1668984054",
				},
				"ttl": &types.AttributeValueMemberN{
					Value: "1668984354",
				},
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
		logger:  zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Timestamp().Logger(),
		counter: 0,
	}
}

func TestHandler_Invoke(t *testing.T) {

	Handler := DummyHandler(t)

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
					"Content-Type":                  "application/json",
					"x-nugg-hex-attestation-object": "0xa363666d74646e6f6e656761747453746d74a06861757468446174615898a9b9abf7fc46b13564b49d5cf85bcbf371f9cb630e0d6b354bc60b51e065da485d000000000000000000000000000000000000000000147053ed09000cfafdd6e1d98d929796f9c07c466ba501020326200121582030dfb831ebb382bcbd45ac6cb1745222b7d81ad8d44ab33e20d2bda632b5692a225820f6496d03d357717d7669a7af490c8706fef052c0819a02bdca4b92bd42459a00",
					"x-nugg-hex-credential-id":      "0x7053ed09000cfafdd6e1d98d929796f9c07c466b",
					"x-nugg-utf-client-data-json":   `{"challenge":"pVr2PUG_le6lde9wxeImHA","origin":"https://nugg.xyz","type":"webauthn.create"}`,
					// "x-nugg-webauthn-creation":      "eyJyYXdDbGllbnREYXRhSlNPTiI6ImV5SmphR0ZzYkdWdVoyVWlPaUp3Vm5JeVVGVkhYMnhsTm14a1pUbDNlR1ZKYlVoQklpd2liM0pwWjJsdUlqb2lhSFIwY0hNNkx5OXVkV2RuTG5oNWVpSXNJblI1Y0dVaU9pSjNaV0poZFhSb2JpNWpjbVZoZEdVaWZRPT0iLCJyYXdBdHRlc3RhdGlvbk9iamVjdCI6Im8yTm1iWFJrYm05dVpXZGhkSFJUZEcxMG9HaGhkWFJvUkdGMFlWaVlxYm1yOS94R3NUVmt0SjFjK0Z2TDgzSDV5Mk1PRFdzMVM4WUxVZUJsMmtoZEFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFGSEJUN1FrQURQcjkxdUhaalpLWGx2bkFmRVpycFFFQ0F5WWdBU0ZZSUREZnVESHJzNEs4dlVXc2JMRjBVaUszMkJyWTFFcXpQaURTdmFZeXRXa3FJbGdnOWtsdEE5TlhjWDEyYWFldlNReUhCdjd3VXNDQm1nSzl5a3VTdlVKRm1nQT0iLCJjcmVkZW50aWFsSWQiOiJjRlB0Q1FBTSt2M1c0ZG1Oa3BlVytjQjhSbXM9In0=",
				},
			},
			want: Output{
				StatusCode: 204,
				Headers: map[string]string{
					"Content-Length":      "0",
					"x-nugg-access-token": "OpenIdToken",
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
					"challenge_id": &types.AttributeValueMemberS{Value: "0xa55af63d41bf95eea575ef70c5e2261c"},
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
					"credential_id": &types.AttributeValueMemberS{Value: "0x7053ed09000cfafdd6e1d98d929796f9c07c466b"},
				},
			})

			if err != nil {
				t.Error(err)
			}

			if c.Item == nil {
				t.Error("item is nil")
			}
			if Handler.counter != len(tests) {
				t.Errorf("Handler.Invoke() counter = %v, want 1", Handler.counter)
			}

		})

	}

}
