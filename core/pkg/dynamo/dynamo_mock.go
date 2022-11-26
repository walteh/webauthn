package dynamo

import (
	"context"
	"log"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/webauthn/types"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Client is a mock of Client interface
type MockClient struct {
	*Client
}

func MockCeremonyTableName() *string {
	return aws.String("ceremony-table")
}

func MockCredentialTableName() *string {
	return aws.String("credential-table")
}

func (cli *MockClient) MockCreateTable(t *testing.T, name *string, pk string) string {
	_, err := cli.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: name,
		AttributeDefinitions: []dtypes.AttributeDefinition{
			{
				AttributeName: aws.String(pk),
				AttributeType: dtypes.ScalarAttributeTypeS,
			},
		},
		KeySchema: []dtypes.KeySchemaElement{
			{
				AttributeName: aws.String(pk),
				KeyType:       dtypes.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &dtypes.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})

	if err != nil {
		t.Error(err)
	}

	return *name
}

func (cli *MockClient) MockDeleteTable(t *testing.T, name *string) {
	_, err := cli.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
		TableName: name,
	})

	if err != nil {
		log.Println(err)
		t.Fail()
	}
}

func MockBatchPut(t *testing.T, cli *Client, items ...dtypes.Put) {
	for _, put := range items {
		if _, err := cli.PutItem(context.Background(), &dynamodb.PutItemInput{
			Item:      put.Item,
			TableName: put.TableName,
		}); err != nil {
			t.Errorf("handler.Dynamo.PutItem() error = %v", err)
			return
		}
	}
}

func MockBatchCheck(t *testing.T, cli *Client, items ...dtypes.Update) {
	for _, put := range items {
		if val, err := cli.GetItem(context.Background(), &dynamodb.GetItemInput{
			Key:       put.Key,
			TableName: put.TableName,
		}); err != nil {
			t.Errorf("handler.Dynamo.PutItem() error = %v", err)
			return
		} else {
			for k, want := range put.ExpressionAttributeValues {
				if !reflect.DeepEqual(val.Item[k], want) {
					t.Errorf("dynamo.MockBatchCheck() [\"%s\"] = %v, want %v", k, val.Item[k], want)
					return
				}
			}
		}

	}
}

func AttachLocalDynamoServer(t *testing.T) {
	cmd := exec.Command("docker", "run", "-d", "-p", "8777:8000", "amazon/dynamodb-local")
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		log.Println("teardown - server")

		cmd.Process.Kill()
	})
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
		o.EndpointResolver = dynamodb.EndpointResolverFromURL("http://localhost:8777")
	})

	mocked := MockClient{
		Client: &Client{Client: cli},
	}

	t.Cleanup(func() {
		log.Println("teardown - tables")
		mocked.MockDeleteTable(t, MockCredentialTableName())
		mocked.MockDeleteTable(t, MockCeremonyTableName())
	})

	return &Client{
		Client:              cli,
		CeremonyTableName:   mocked.MockCreateTable(t, MockCeremonyTableName(), "challenge_id"),
		CredentialTableName: mocked.MockCreateTable(t, MockCredentialTableName(), "credential_id"),
	}
}

func (dynamoClient *MockClient) MockSetCeremony(t *testing.T, credential hex.Hash, challenge hex.Hash, type_ types.CeremonyType) {
	t.Helper()

	cer := types.NewCeremony(credential, challenge, type_)

	maper, err := cer.Put()
	if err != nil {
		t.Fatal(err)
	}

	maper.TableName = aws.String(dynamoClient.CeremonyTableName)

	err = dynamoClient.TransactWrite(context.Background(),
		dtypes.TransactWriteItem{Put: maper},
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
