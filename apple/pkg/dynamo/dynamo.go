package dynamo

import (
	"context"
	"log"
	"nugg-auth/apple/pkg/random"
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
		Client: dynamodb.NewFromConfig(config),
	}
}

type Challenge struct {
	Challenge string `json:"challenge"`
	UserId    string `json:"userId"`
	Ttl       int64  `json:"ttl"`
}

// AddMovie adds a movie the DynamoDB table.
func (basics *Client) GenerateChallenge(ctx context.Context, userId string, life int64) (string, error) {

	challenge := Challenge{
		Challenge: random.Sequence(userId),
		UserId:    userId,
		Ttl:       time.Now().Unix() + life,
	}

	item, err := attributevalue.MarshalMap(challenge)
	if err != nil {
		panic(err)
	}
	_, err = basics.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return challenge.Challenge, nil
}
