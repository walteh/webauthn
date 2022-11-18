package main

import (
	"context"
	"encoding/base64"
	"log"
	"nugg-auth/core/pkg/dynamo"
	"nugg-auth/core/pkg/webauthn/webauthn"

	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

func DummyHandler(t *testing.T) *Handler {
	dynamoClient := dynamo.NewMockClient(t)

	wan, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "nugg.xyz",
		RPID:          "nugg.xyz",
		RPOrigin:      "https://nugg.xyz",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = dynamoClient.TransactWrite(context.Background(),
		types.TransactWriteItem{
			Put: &types.Put{
				Item: map[string]types.AttributeValue{
					"user_id": &types.AttributeValueMemberS{
						Value: "01GJ64VW4WB50CYYN2KDN94WKH",
					},
					"created_at": &types.AttributeValueMemberS{
						Value: "2022-11-18 07:22:54.668188 -0600 CST m=+0.032937709",
					},
					"updated_at": &types.AttributeValueMemberS{
						Value: "2022-11-18 07:22:54.668281 -0600 CST m=+0.033030751",
					},
				},
				TableName: &dynamoClient.UserTableName,
			},
		},
		types.TransactWriteItem{
			Put: &types.Put{
				Item: map[string]types.AttributeValue{
					"dat": &types.AttributeValueMemberB{
						Value: []uint8{
							0x7b, 0x22, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x22,
							0x3a, 0x22, 0x63, 0x46, 0x50, 0x74, 0x43, 0x51, 0x41, 0x4d, 0x2b, 0x76, 0x33, 0x57, 0x34, 0x64,
							0x6d, 0x4e, 0x6b, 0x70, 0x65, 0x57, 0x2b, 0x63, 0x42, 0x38, 0x52, 0x6d, 0x73, 0x3d, 0x22, 0x2c,
							0x22, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x22, 0x3a, 0x22, 0x70, 0x51,
							0x45, 0x43, 0x41, 0x79, 0x59, 0x67, 0x41, 0x53, 0x46, 0x59, 0x49, 0x44, 0x44, 0x66, 0x75, 0x44,
							0x48, 0x72, 0x73, 0x34, 0x4b, 0x38, 0x76, 0x55, 0x57, 0x73, 0x62, 0x4c, 0x46, 0x30, 0x55, 0x69,
							0x4b, 0x33, 0x32, 0x42, 0x72, 0x59, 0x31, 0x45, 0x71, 0x7a, 0x50, 0x69, 0x44, 0x53, 0x76, 0x61,
							0x59, 0x79, 0x74, 0x57, 0x6b, 0x71, 0x49, 0x6c, 0x67, 0x67, 0x39, 0x6b, 0x6c, 0x74, 0x41, 0x39,
							0x4e, 0x58, 0x63, 0x58, 0x31, 0x32, 0x61, 0x61, 0x65, 0x76, 0x53, 0x51, 0x79, 0x48, 0x42, 0x76,
							0x37, 0x77, 0x55, 0x73, 0x43, 0x42, 0x6d, 0x67, 0x4b, 0x39, 0x79, 0x6b, 0x75, 0x53, 0x76, 0x55,
							0x4a, 0x46, 0x6d, 0x67, 0x41, 0x3d, 0x22, 0x2c, 0x22, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x61,
							0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x22, 0x6e, 0x6f, 0x6e, 0x65,
							0x22, 0x2c, 0x22, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72,
							0x22, 0x3a, 0x7b, 0x22, 0x61, 0x61, 0x67, 0x75, 0x69, 0x64, 0x22, 0x3a, 0x22, 0x41, 0x41, 0x41,
							0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
							0x41, 0x41, 0x41, 0x3d, 0x3d, 0x22, 0x2c, 0x22, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x63, 0x6f, 0x75,
							0x6e, 0x74, 0x22, 0x3a, 0x30, 0x2c, 0x22, 0x63, 0x6c, 0x6f, 0x6e, 0x65, 0x5f, 0x77, 0x61, 0x72,
							0x6e, 0x69, 0x6e, 0x67, 0x22, 0x3a, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x7d, 0x7d,
						},
					},
					"credential_id": &types.AttributeValueMemberS{
						Value: "cFPtCQAM-v3W4dmNkpeW-cB8Rms",
					},
					"nugg_id": &types.AttributeValueMemberS{
						Value: "01GJ64VW4WB50CYYN2KDN94WKH",
					},
					"credential_user_id": &types.AttributeValueMemberB{
						Value: []uint8{
							0x68, 0x4f, 0x43, 0x65, 0x35, 0x38, 0x38, 0x64, 0x61, 0x4a, 0x5a, 0x75, 0x47, 0x41, 0x36, 0x50,
							0x65, 0x4a, 0x54, 0x63, 0x7a, 0x51,
						},
					},
					"type": &types.AttributeValueMemberS{
						Value: "apple-passkey",
					},
					"created_at": &types.AttributeValueMemberN{
						Value: "1668801228",
					},
					"updated_at": &types.AttributeValueMemberN{
						Value: "1668801228",
					},
				},
				TableName: &dynamoClient.CredentialTableName,
			},
		}, types.TransactWriteItem{
			Put: &types.Put{
				Item: map[string]types.AttributeValue{
					"session_data": &types.AttributeValueMemberM{
						Value: map[string]types.AttributeValue{
							"challenge": &types.AttributeValueMemberB{
								Value: []uint8{
									0xe1, 0x2e, 0x11, 0x5a, 0xcf, 0x45, 0x52, 0xb2, 0x56, 0x8b, 0x55, 0xe9, 0x3c, 0xbd, 0x39, 0x39,
								},
							},
							"user_id": &types.AttributeValueMemberB{
								Value: []uint8{
									0x84, 0xe0, 0x9e, 0xe7, 0xcf, 0x1d, 0x68, 0x96, 0x6e, 0x18, 0x0e, 0x8f, 0x78, 0x94, 0xdc, 0xcd,
								},
							},
							"allowed_credentials": &types.AttributeValueMemberL{
								Value: []types.AttributeValue{
									&types.AttributeValueMemberB{
										Value: []uint8{
											0x70, 0x53, 0xed, 0x09, 0x00, 0x0c, 0xfa, 0xfd, 0xd6, 0xe1, 0xd9, 0x8d, 0x92, 0x97, 0x96, 0xf9,
											0xc0, 0x7c, 0x46, 0x6b,
										},
									},
								},
							},
							"user_verification": &types.AttributeValueMemberS{
								Value: "",
							},
						},
					},
					"ttl": &types.AttributeValueMemberN{
						Value: "1669101228",
					},
					"ceremony_id": &types.AttributeValueMemberS{
						Value: "4S4RWs9FUrJWi1XpPL05OQ",
					},
				},
				TableName: &dynamoClient.CeremonyTableName,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

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

	Handler := DummyHandler(t)

	val := []byte{
		0x4d, 0x44, 0x41, 0x78, 0x4e, 0x44, 0x4d, 0x33, 0x4c, 0x6d, 0x52, 0x6c, 0x5a, 0x6a, 0x55, 0x7a,
		0x4e, 0x57, 0x52, 0x6b, 0x5a, 0x44, 0x6c, 0x6c, 0x4d, 0x6a, 0x52, 0x6a, 0x4e, 0x47, 0x5a, 0x68,
		0x4e, 0x44, 0x4d, 0x32, 0x4e, 0x32, 0x52, 0x6a, 0x59, 0x54, 0x55, 0x77, 0x5a, 0x6d, 0x52, 0x6d,
		0x5a, 0x57, 0x52, 0x69, 0x4c, 0x6a, 0x45, 0x35, 0x4e, 0x54, 0x45, 0x3d,
	}

	res, _ := base64.StdEncoding.DecodeString(string(val))

	log.Println("val", string(res))

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
					// "x-nugg-webauthn-assertion": "eyJyYXdBdXRoZW50aWNhdG9yRGF0YSI6InFibXI5XC94R3NUVmt0SjFjK0Z2TDgzSDV5Mk1PRFdzMVM4WUxVZUJsMmtnZEFBQUFBQT09Iiwic2lnbmF0dXJlIjoiTUVZQ0lRRG1oUTRnMjRPb3pEV0Y4UlZZNEplZHBvVHlUT1JhbnpjRkIxWnc1S2lvOGdJaEFMUjJZUjRnTHI0ZGZ5ajJWQlhxeTc5SzJCT2xDOWY4dWYzeDcyQldmb3AzIiwidXNlcklEIjoiTURVd1FVUXdRVUV0UWpneU9TMDBRVUl6TFVGR1F6UXROVE0wTURBNVJqWkdOalkwIiwicmF3Q2xpZW50RGF0YUpTT04iOiJleUowZVhCbElqb2lkMlZpWVhWMGFHNHVaMlYwSWl3aVkyaGhiR3hsYm1kbElqb2lVVlZLUkZKQklpd2liM0pwWjJsdUlqb2lhSFIwY0hNNkx5OXVkV2RuTG5oNWVpSjkiLCJjcmVkZW50aWFsSUQiOiJPWFdGTk5qSDFNK3pTNDB2MHRISWRudlo4T009IiwiY3JlZGVudGlhbFR5cGUiOiJwdWJsaWMta2V5In0=",
					"x-nugg-webauthn-assertion": "eyJyYXdBdXRoZW50aWNhdG9yRGF0YSI6InFibXI5XC94R3NUVmt0SjFjK0Z2TDgzSDV5Mk1PRFdzMVM4WUxVZUJsMmtnZEFBQUFBQT09Iiwic2lnbmF0dXJlIjoiTUVVQ0lBOERoRHNkRjhYY3dENlQ5WDBSMUM2OG9lRncrZ05neTFsSE1ZR3hpXC9XSEFpRUE2K0pYcGJMZFkzOWQ2Zks5b0RSRHBMdERBdjdEcGxTbDdwK05tXC9OaUZKYz0iLCJ1c2VySUQiOiJoT0NlNTg4ZGFKWnVHQTZQZUpUY3pRPT0iLCJyYXdDbGllbnREYXRhSlNPTiI6ImV5SjBlWEJsSWpvaWQyVmlZWFYwYUc0dVoyVjBJaXdpWTJoaGJHeGxibWRsSWpvaU5GTTBVbGR6T1VaVmNrcFhhVEZZY0ZCTU1EVlBVU0lzSW05eWFXZHBiaUk2SW1oMGRIQnpPaTh2Ym5Wblp5NTRlWG9pZlE9PSIsImNyZWRlbnRpYWxJRCI6ImNGUHRDUUFNK3YzVzRkbU5rcGVXK2NCOFJtcz0iLCJjcmVkZW50aWFsVHlwZSI6InB1YmxpYy1rZXkifQ==",
					// "x-nugg-webauthn-signature":          "MEUCIQDjOjqf0rsEHbD1tTBS7dak7cDGXTcpQT0URWwDEWckfQIgSX9x74X6Bx4cL7Du6qUZ+pcUD74pPUnjLvq9/HlQoMQ=",
					// "x-nugg-webauthn-clientdata":         "{\"type\":\"webauthn.get\",\"challenge\":\"4S4RWs9FUrJWi1XpPL05OQ\",\"origin\":\"https://nugg.xyz\"}",
					// "x-nugg-webauthn-credential-id":      "cFPtCQAM-v3W4dmNkpeW-cB8Rms",
					// "x-nugg-webauthn-user-id":            "hOCe588daJZuGA6PeJTczQ",
					// "x-nugg-webauthn-authenticator-data": "qbmr9_xGsTVktJ1c-FvL83H5y2MODWs1S8YLUeBl2kgdAAAAAA",
				},
			},
			want: Output{
				StatusCode: 204,
				Headers: map[string]string{
					"Content-Length": "0",
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
