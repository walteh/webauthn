package dynamo

import (
	"context"
	"log"
	"nugg-auth/core/pkg/random"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Challenge struct {
	Id    string `dynamodbav:"id"`
	State string `dynamodbav:"state"`
	Ttl   int64  `dynamodbav:"ttl"`
}

// AddMovie adds a movie the DynamoDB table.
func (client *Client) GenerateChallenge(ctx context.Context, state string, ttl time.Time) (string, error) {

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

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(client.TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
		return "", err
	}
	return challenge.Id, nil
}
