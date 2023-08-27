package otelawsc

import (
	"context"
	"reflect"

	v2Middleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/smithy-go/middleware"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var servicemap = map[string]otelaws.AttributeSetter{
	dynamodb.ServiceID: DynamoDBAttributeSetter,
	sqs.ServiceID:      otelaws.SQSAttributeSetter,
}

// OperationAttr returns the AWS operation attribute.
func OperationAttr(operation string) attribute.KeyValue {
	return otelaws.OperationAttr(operation)
}

// RegionAttr returns the AWS region attribute.
func RegionAttr(region string) attribute.KeyValue {
	return otelaws.RegionKey.String(region)
}

// ServiceAttr returns the AWS service attribute.
func ServiceAttr(service string) attribute.KeyValue {
	return otelaws.OperationAttr(service)
}

// RequestIDAttr returns the AWS request ID attribute.
func RequestIDAttr(requestID string) attribute.KeyValue {
	return otelaws.RequestIDKey.String(requestID)
}

// DefaultAttributeSetter checks to see if there are service specific attributes available to set for the AWS service.
// If there are service specific attributes available then they will be included.
func DefaultAttributeSetter(ctx context.Context, in middleware.InitializeInput) []attribute.KeyValue {
	serviceID := v2Middleware.GetServiceID(ctx)

	if fn, ok := servicemap[serviceID]; ok {
		return fn(ctx, in)
	}

	return []attribute.KeyValue{}
}

// DynamoDBAttributeSetter sets DynamoDB specific attributes depending on the DynamoDB operation being performed.
func DynamoDBAttributeSetter(ctx context.Context, in middleware.InitializeInput) []attribute.KeyValue {
	dynamodbAttributes := []attribute.KeyValue{semconv.DBSystemKey.String("dynamodb")}

	zerolog.Ctx(ctx).Debug().Any("paramaters", in.Parameters).Str("argtype", reflect.TypeOf(in.Parameters).String()).Msg("DynamoDBAttributeSetter")

	switch v := in.Parameters.(type) {
	case *dynamodb.TransactWriteItemsInput:
		tableNames := make(map[string]bool)

		for _, z := range v.TransactItems {
			if z.Update != nil {
				tableNames[*z.Update.TableName] = true
			} else if z.Put != nil {
				tableNames[*z.Put.TableName] = true
			} else if z.Delete != nil {
				tableNames[*z.Delete.TableName] = true
			} else if z.ConditionCheck != nil {
				tableNames[*z.ConditionCheck.TableName] = true
			}
		}

		tableNamesSlice := make([]string, 0, len(tableNames))
		for k := range tableNames {
			tableNamesSlice = append(tableNamesSlice, k)
		}

		if len(tableNamesSlice) == 1 {
			dynamodbAttributes = append(dynamodbAttributes, semconv.AWSDynamoDBTableNamesKey.String(tableNamesSlice[0]))
		} else {
			dynamodbAttributes = append(dynamodbAttributes, semconv.AWSDynamoDBTableNamesKey.StringSlice(tableNamesSlice))
		}

		trace.SpanFromContext(ctx).SetAttributes(attribute.String("origin", "AWS::DynamoDB::Table"))

	case *dynamodb.TransactGetItemsInput:

		tableNames := make(map[string]bool)

		for _, z := range v.TransactItems {
			tableNames[*z.Get.TableName] = true
		}

		tableNamesSlice := make([]string, 0, len(tableNames))
		for k := range tableNames {
			tableNamesSlice = append(tableNamesSlice, k)
		}

		if len(tableNamesSlice) == 1 {
			dynamodbAttributes = append(dynamodbAttributes, semconv.AWSDynamoDBTableNamesKey.String(tableNamesSlice[0]))
		} else {
			dynamodbAttributes = append(dynamodbAttributes, semconv.AWSDynamoDBTableNamesKey.StringSlice(tableNamesSlice))
		}
	default:
		dynamodbAttributes = append(dynamodbAttributes, otelaws.DynamoDBAttributeSetter(ctx, in)...)
	}

	return dynamodbAttributes
}
