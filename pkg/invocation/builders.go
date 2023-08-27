package invocation

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func (me *Options) NewS3Client() *s3.Client {
	return s3.NewFromConfig(*me.awsConfig)
}

func (me *Options) NewSNSClient() *sns.Client {
	return sns.NewFromConfig(*me.awsConfig)
}

func (me *Options) NewSQSClient() *sqs.Client {
	return sqs.NewFromConfig(*me.awsConfig)
}

func (me *Options) NewDynamoDBClient() *dynamodb.Client {
	return dynamodb.NewFromConfig(*me.awsConfig)
}

func (me *Options) NewIOTCoreClient() *iotdataplane.Client {
	return iotdataplane.NewFromConfig(*me.awsConfig)
}
