package otelawsc

// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"reflect"

	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/rs/zerolog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv/v1.13.0/httpconv"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "git.nugg.xyz/go-sdk/otel/otelawsc"
)

// type spanTimestampKey struct{}

// // AttributeSetter returns an array of KeyValue pairs, it can be used to set custom attributes.
// type AttributeSetter func(context.Context, middleware.InitializeInput) []attribute.KeyValue

// type otelMiddlewares struct {
// 	tracer          trace.Tracer
// 	propagator      propagation.TextMapPropagator
// 	attributeSetter []AttributeSetter
// }

// func (m otelMiddlewares) initializeMiddlewareBefore(stack *middleware.Stack) error {
// 	return stack.Initialize.Add(middleware.InitializeMiddlewareFunc("OTelInitializeMiddlewareBefore", func(
// 		ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
// 		out middleware.InitializeOutput, metadata middleware.Metadata, err error) {

// 		serviceID := v2Middleware.GetServiceID(ctx)

// 		attributes := []attribute.KeyValue{
// 			ServiceAttr(serviceID),
// 			RegionAttr(v2Middleware.GetRegion(ctx)),
// 			OperationAttr(v2Middleware.GetOperationName(ctx)),
// 		}
// 		for _, setter := range m.attributeSetter {
// 			attributes = append(attributes, setter(ctx, in)...)
// 		}

// 		ctx, _ = m.tracer.Start(ctx, serviceID,
// 			trace.WithTimestamp(time.Now()),
// 			trace.WithSpanKind(trace.SpanKindClient),
// 			trace.WithAttributes(attributes...),
// 		)

// 		return next.HandleInitialize(ctx, in)
// 	}),
// 		middleware.Before)
// }

// func (m otelMiddlewares) initializeMiddlewareAfter(stack *middleware.Stack) error {
// 	return stack.Initialize.Add(middleware.InitializeMiddlewareFunc("OTelInitializeMiddlewareAfter", func(
// 		ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
// 		out middleware.InitializeOutput, metadata middleware.Metadata, err error) {

// 		span := trace.SpanFromContext(ctx)

// 		defer span.End()

// 		out, metadata, err = next.HandleInitialize(ctx, in)
// 		if err != nil {
// 			span.RecordError(err)
// 			span.SetStatus(codes.Error, err.Error())
// 		}

// 		return out, metadata, err
// 	}),
// 		middleware.After)
// }

// func (m otelMiddlewares) deserializeMiddleware(stack *middleware.Stack) error {
// 	return stack.Deserialize.Add(middleware.DeserializeMiddlewareFunc("OTelDeserializeMiddleware", func(
// 		ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
// 		out middleware.DeserializeOutput, metadata middleware.Metadata, err error) {
// 		out, metadata, err = next.HandleDeserialize(ctx, in)
// 		resp, ok := out.RawResponse.(*smithyhttp.Response)
// 		if !ok {
// 			// No raw response to wrap with.
// 			return out, metadata, err
// 		}
// 		span := trace.SpanFromContext(ctx)
// 		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(resp.StatusCode))
// 		span.SetAttributes(httpconv.ClientResponse(resp.Response)...)
// 		span.SetStatus(httpconv.ClientStatus(resp.StatusCode))

// 		requestID, ok := v2Middleware.GetRequestIDMetadata(metadata)
// 		if ok {
// 			span.SetAttributes(RequestIDAttr(requestID))
// 		}

// 		return out, metadata, err
// 	}),
// 		middleware.Before)
// }

// func (m otelMiddlewares) finalizeMiddleware(stack *middleware.Stack) error {
// 	return stack.Finalize.Add(middleware.FinalizeMiddlewareFunc("OTelFinalizeMiddleware", func(
// 		ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (
// 		out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
// 		// Propagate the Trace information by injecting it into the HTTP request.
// 		switch req := in.Request.(type) {
// 		case *smithyhttp.Request:
// 			m.propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
// 		default:
// 		}

// 		return next.HandleFinalize(ctx, in)
// 	}),
// 		middleware.Before)
// }

// // AppendMiddlewares attaches OTel middlewares to the AWS Go SDK V2 for instrumentation.
// // OTel middlewares can be appended to either all aws clients or a specific operation.
// // Please see more details in https://aws.github.io/aws-sdk-go-v2/docs/middleware/
// func AppendMiddlewares(apiOptions *[]func(*middleware.Stack) error, tracerProvider trace.TracerProvider, abc propagation.TextMapPropagator) {

// 	m := otelMiddlewares{tracer: tracerProvider.Tracer(tracerName,
// 		trace.WithInstrumentationVersion(otelaws.SemVersion())),
// 		propagator:      abc,
// 		attributeSetter: []AttributeSetter{DefaultAttributeSetter},
// 	}
// 	*apiOptions = append(*apiOptions, m.initializeMiddlewareBefore, m.initializeMiddlewareAfter, m.finalizeMiddleware, m.deserializeMiddleware)
// }

func DeserializeHttpMiddleware(stack *middleware.Stack) error {
	return stack.Deserialize.Add(middleware.DeserializeMiddlewareFunc("OTelDeserializeHttpMiddleware", func(
		ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
		out middleware.DeserializeOutput, metadata middleware.Metadata, err error) {
		out, metadata, err = next.HandleDeserialize(ctx, in)
		resp, ok := out.RawResponse.(*smithyhttp.Response)
		if !ok {
			// No raw response to wrap with.
			return out, metadata, err
		}

		span := trace.SpanFromContext(ctx)

		span.SetAttributes(httpconv.ClientResponse(resp.Response)...)
		span.SetStatus(httpconv.ClientStatus(resp.StatusCode))

		return out, metadata, err
	}),
		middleware.Before)
}

func InitializeHttpMiddleware(stack *middleware.Stack) error {
	return stack.Initialize.Add(middleware.InitializeMiddlewareFunc("OTelInitializeHttpMiddleware", func(
		ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
		out middleware.InitializeOutput, metadata middleware.Metadata, err error) {

		span := trace.SpanFromContext(ctx)

		children := make([]trace.Span, 0)

		zerolog.Ctx(ctx).Debug().Any("paramaters", in.Parameters).Str("argtype", reflect.TypeOf(in.Parameters).String()).Msg("DynamoDBAttributeSetter")

		// switch v := in.Parameters.(type) {
		// case *dynamodb.TransactWriteItemsInput:

		// 	for _, item := range v.TransactItems {
		// 		_, spn := span.TracerProvider().Tracer(tracerName).Start(ctx, "DynamoDB", trace.WithSpanKind(trace.SpanKindServer))
		// 		spn.SetAttributes(semconv.DBSystemKey.String("dynamodb"))
		// 		spn.SetAttributes(otelaws.ServiceAttr("dynamodb"))
		// 		spn.SetAttributes(attribute.String("namespace", "aws"))

		// 		if item.Put != nil {
		// 			spn.SetAttributes(attribute.String(AWSOperationAttribute, "PutItem"))
		// 			spn.SetAttributes(semconv.AWSDynamoDBTableNamesKey.String(*item.Put.TableName))
		// 		} else if item.Update != nil {
		// 			spn.SetAttributes(attribute.String(AWSOperationAttribute, "UpdateItem"))
		// 			spn.SetAttributes(semconv.AWSDynamoDBTableNamesKey.String(*item.Update.TableName))
		// 		} else if item.Delete != nil {
		// 			spn.SetAttributes(attribute.String(AWSOperationAttribute, "DeleteItem"))
		// 			spn.SetAttributes(semconv.AWSDynamoDBTableNamesKey.String(*item.Delete.TableName))
		// 		} else if item.ConditionCheck != nil {
		// 			spn.SetAttributes(attribute.String(AWSOperationAttribute, "ConditionCheck"))
		// 			spn.SetAttributes(semconv.AWSDynamoDBTableNamesKey.String(*item.ConditionCheck.TableName))
		// 		}

		// 		children = append(children, spn)

		// 	}

		// }

		defer span.End()

		span.SetAttributes(attribute.String(AWSOperationAttribute, "UpdateItem"))

		defer func() {
			for _, child := range children {
				child.End()
			}
		}()

		out, metadata, err = next.HandleInitialize(ctx, in)
		if err != nil {
			for _, child := range children {
				child.RecordError(err)
				child.SetStatus(codes.Error, err.Error())
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		return out, metadata, err
	}),
		middleware.After)
}

// func MakeTraceOkOnConditionalCheckFailMiddleware(stack *middleware.Stack) error {
// 	return stack.Deserialize.Add(middleware.DeserializeMiddlewareFunc("OTelDeserializeMakeTraceOkOnConditionalCheckFailMiddleware", func(
// 		ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
// 		out middleware.DeserializeOutput, metadata middleware.Metadata, err error) {
// 		out, metadata, err = next.HandleDeserialize(ctx, in)
// 		zerolog.Ctx(ctx).Debug().Any("out", out).Msg("DynamoDBAttributeSetter")
// 		zerolog.Ctx(ctx).Debug().Any("metadata", metadata).Msg("DynamoDBAttributeSetter")
// 		if err != nil {

// 			if err == nil {
// 				zerolog.Ctx(ctx).Debug().Any("err", err).Msg("DynamoDBAttributeSetter")
// 				if _, ok := errd.As[*types.ConditionalCheckFailedException](err); ok {
// 					err = nil
// 				} else if er, ok := errd.As[*types.TransactionCanceledException](err); ok {
// 					ok := true
// 					for _, cause := range er.CancellationReasons {
// 						if *cause.Code == "ConditionalCheckFailed" {
// 							continue
// 						}
// 						ok = false
// 						break
// 					}
// 					if ok {
// 						err = nil
// 					}
// 				}
// 			}
// 		}

// 		return out, metadata, err
// 	}),
// 		middleware.Before)
// }
