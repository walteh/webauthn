package dynamo

import (
	"context"
	"log"
	"nugg-auth/challenge/pkg/random"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Client struct {
	*dynamodb.Client
	TableName string
	// contains filtered or unexported fields
}

func NewClient(config aws.Config, tableName string) *Client {
	return &Client{
		dynamodb.NewFromConfig(config),
		tableName,
	}
}

type Challenge struct {
	Id    string `dynamodbav:"id"`
	State string `dynamodbav:"state"`
	Ttl   int64  `dynamodbav:"ttl"`
}

// AddMovie adds a movie the DynamoDB table.
func (basics *Client) GenerateChallenge(ctx context.Context, state string, ttl time.Time) (string, error) {

	challenge := Challenge{
		Id:    random.KSUID(),
		State: state,
		Ttl:   ttl.Unix(),
	}

	log.Println("challenge", challenge)

	item, err := attributevalue.MarshalMap(challenge)
	if err != nil {
		log.Printf("failed to marshal challenge, %v", err)
		return "", err
	}

	log.Println("putting item in dynamo", item)

	_, err = basics.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
		return "", err
	}
	return challenge.Id, nil
}
