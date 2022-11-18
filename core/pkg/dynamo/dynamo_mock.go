package dynamo

import (
	"context"
	"encoding/base64"
	"log"
	"nugg-auth/core/pkg/webauthn/webauthn"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Client is a mock of Client interface
type MockClient struct {
	*Client
}

func (cli *MockClient) MockCreateTable(t *testing.T, name string, pk string) string {
	_, err := cli.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String(name),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(pk),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(pk),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	return name
}

func (cli *MockClient) MockDeleteTable(t *testing.T, name string) {
	_, err := cli.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
		TableName: aws.String(name),
	})

	if err != nil {
		log.Println(err)
		t.Fail()
	}
}

func NewMockClient(t *testing.T) *Client {
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
	os.Setenv("AWS_REGION", "local")

	conf, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	cli := dynamodb.NewFromConfig(conf, func(o *dynamodb.Options) {
		o.EndpointResolver = dynamodb.EndpointResolverFromURL("http://localhost:8000")
	})

	mocked := MockClient{
		Client: &Client{Client: cli},
	}

	t.Cleanup(func() {
		log.Println("teardown")
		mocked.MockDeleteTable(t, "credential-table")
		mocked.MockDeleteTable(t, "user-table")
		mocked.MockDeleteTable(t, "ceremony-table")
	})

	return &Client{
		Client:              cli,
		UserTableName:       mocked.MockCreateTable(t, "user-table", "user_id"),
		CeremonyTableName:   mocked.MockCreateTable(t, "ceremony-table", "ceremony_id"),
		CredentialTableName: mocked.MockCreateTable(t, "credential-table", "credential_id"),
	}
}

func (dynamoClient *Client) MockSetCeremony(t *testing.T, wan *webauthn.WebAuthn, challenge string) {
	t.Helper()

	_, cer, err := wan.BeginRegistration("tester1")
	if err != nil {
		t.Fatal(err)
	}

	cer.Challenge, err = base64.URLEncoding.DecodeString(challenge)
	if err != nil {
		t.Fatal(err)
	}
	put, err := dynamoClient.NewCeremonyPut(cer)
	if err != nil {
		t.Fatal(err)
	}

	err = dynamoClient.TransactWrite(context.Background(),
		types.TransactWriteItem{Put: put},
	)
	if err != nil {
		t.Fatal(err)
	}
}

// func (dynamoClient *Client) MockPretendRegisterHappened(t *testing.T, wan *webauthn.WebAuthn, challenge string, attestation string) {
// 	t.Helper()

// 	_, cer, err := wan.BeginRegistration()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	cer.Challenge = challenge

// 	ceremonyPut, err := dynamoClient.NewCeremonyPut(cer)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	dynamoClient.TransactWrite(context.Background(),
// 		types.TransactWriteItem{Put: userput},
// 		types.TransactWriteItem{Put: credput},
// 		types.TransactWriteItem{Put: ceremonyPut},
// 	)

// 	log.
// }
