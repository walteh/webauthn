package secretsmanager

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/smithy-go/middleware"
)

type ClientWithAWSMocked struct {
	*Client
}

func NewMockClient() (client *ClientWithAWSMocked) {
	return &ClientWithAWSMocked{&Client{}}
}

func (c *ClientWithAWSMocked) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return &secretsmanager.GetSecretValueOutput{
		SecretString:   params.SecretId,
		VersionId:      params.VersionId,
		VersionStages:  []string{*params.VersionStage},
		ARN:            params.SecretId,
		Name:           params.SecretId,
		CreatedDate:    aws.Time(time.Now()),
		SecretBinary:   []byte(*params.SecretId),
		ResultMetadata: middleware.Metadata{},
	}, nil
}
