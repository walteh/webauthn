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

package lambda // import "go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/k0kubun/pp/v3"
	"github.com/rs/zerolog"
	"github.com/walteh/webauthn/pkg/invocation"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var errorLogger = log.New(log.Writer(), "OTel Lambda Error: ", 0)

const (
	tracerName = "git.nugg.xyz/go-sdk/otel/lambda"
)

type HandlerFunc[I any, O any] func(context.Context, I) (O, error)
type RawHandlerFunc HandlerFunc[[]byte, []byte]

type invokerWrapper[I any, O any] struct {
	handler HandlerFunc[I, O]
}

func (iw *invokerWrapper[I, O]) Invoke(ctx context.Context, in I) (O, error) {
	return iw.handler(ctx, in)
}

func InstrumentHandlerFunc[I any, O any](ctx context.Context, handler HandlerFunc[I, O], traceProvider trace.TracerProvider) lambda.Handler {
	return InstrumentHandler[I, O](ctx, &invokerWrapper[I, O]{handler: handler}, traceProvider)
}

type SimpleHandler struct {
	f func(ctx context.Context, eventJSON []byte) ([]byte, error)
}

func (h *SimpleHandler) Invoke(ctx context.Context, eventJSON []byte) ([]byte, error) {
	return h.f(ctx, eventJSON)
}

func InstrumentHandler[I any, O any](ctx context.Context, invoker invocation.Invokable[I, O], traceProvider trace.TracerProvider) lambda.Handler {

	instrumentor := &instrumentor[I, O]{
		tracer:   traceProvider.Tracer(tracerName, trace.WithInstrumentationVersion(Version())),
		resAttrs: []attribute.KeyValue{},
	}

	log := zerolog.Ctx(ctx).With().Str("component", "otel/lambda").Logger()

	log.Info().Str("component", "otel/lambda").Msg("instrumenting lambda handler")

	f := func(ctx context.Context, eventJSON []byte) ([]byte, error) {

		log.Debug().Str("component", "otel/lambda").Msg("invoking lambda handler")

		log.Debug().Str("component", "otel/lambda").Msg("unmarshalling event")

		var in I
		err := json.Unmarshal(eventJSON, &in)
		if err != nil {
			return nil, err
		}

		ctx, span := instrumentor.tracingBegin(ctx, eventJSON, in)
		defer instrumentor.tracingEnd(ctx, span)

		log.Debug().Str("component", "otel/lambda").Msg("invoking lambda handler")

		out, err := invoker.Invoke(ctx, in)
		if err != nil {
			return nil, err
		}

		log.Debug().Str("component", "otel/lambda").Msg("marshalling response")

		outJSON, err := json.Marshal(out)
		if err != nil {
			return nil, err
		}

		log.Debug().Str("component", "otel/lambda").Msg("returning response")

		return outJSON, nil

	}

	return &SimpleHandler{f}
}

type instrumentor[I any, O any] struct {
	tracer   trace.Tracer
	resAttrs []attribute.KeyValue
}

func extractXRayTraceIDs(ev interface{}) []string {
	xrays := []string{}
	if z, ok := ev.(interface {
		XRayTraceID() []string
	}); ok {
		xrays = z.XRayTraceID()
	}
	links := []trace.Link{}

	for _, id := range xrays {
		tmc := propagation.HeaderCarrier{}
		tmc.Set("X-Amzn-Trace-Id", id)
		ctx := xray.Propagator{}.Extract(context.TODO(), tmc)
		links = append(links, trace.LinkFromContext(ctx))
	}

	pp.Println(links)

	return xrays

}

// Logic to start OTel Tracing.
func (i *instrumentor[I, O]) tracingBegin(ctx context.Context, eventJSON []byte, input I) (context.Context, trace.Span) {

	// links := []trace.Link{}

	// xrays := extractXRayTraceIDs(input)

	tmc := propagation.HeaderCarrier{}
	tmc.Set("X-Amzn-Trace-Id", os.Getenv("_X_AMZN_TRACE_ID"))
	ctx = xray.Propagator{}.Extract(ctx, tmc)

	spn := trace.SpanContextFromContext(ctx).WithRemote(false)
	// links = append(links, trace.LinkFromContext(ctx))

	ctx = trace.ContextWithSpanContext(ctx, spn)

	span := trace.SpanFromContext(ctx)
	// spanName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")

	var attributes []attribute.KeyValue
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		errorLogger.Println("failed to load lambda context from context, ensure tracing enabled in Lambda")
	}
	if lc != nil {
		ctxRequestID := lc.AwsRequestID
		attributes = append(attributes, semconv.FaaSExecutionKey.String(ctxRequestID))

		// Some resource attrs added as span attrs because lambda
		// resource detectors are created before a lambda
		// invocation and therefore lack lambdacontext.
		// Create these attrs upon first invocation
		if len(i.resAttrs) == 0 {
			ctxFunctionArn := lc.InvokedFunctionArn
			attributes = append(attributes, semconv.FaaSIDKey.String(ctxFunctionArn))
			arnParts := strings.Split(ctxFunctionArn, ":")
			if len(arnParts) >= 5 {
				attributes = append(attributes, semconv.CloudAccountIDKey.String(arnParts[4]))
			}
		}
		attributes = append(attributes, i.resAttrs...)
	}

	log := zerolog.New(os.Stdout).With().Str("component", "otel/lambda").Logger()

	log.Debug().Any("lambdacontext", lc).Any("arg", span).Msg("starting tracing")

	// ctx, span = i.tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attributes...), trace.WithLinks(links...))

	// xray.Propagator{}.Inject(ctx, &lambdacontext.XRayHeader{Root: span.SpanContext().TraceID().String(), Parent: span.SpanContext().SpanID().String(), Sampled: "1"})

	// trace.SpanFromContext(ctx).AddEvent("invoked", trace.WithAttributes(attribute.String("event", string(eventJSON))))

	return ctx, span
}

// Logic to wrap up OTel Tracing.
func (i *instrumentor[I, O]) tracingEnd(ctx context.Context, span trace.Span) {
	span.End()

	// // force flush any tracing data since lambda may freeze
	// err := i.configuration.Flusher.ForceFlush(ctx)
	// if err != nil {
	// 	errorLogger.Println("failed to force a flush, lambda may freeze before instrumentation exported: ", err)
	// }
}
