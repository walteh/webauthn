package invocation

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/rs/zerolog"
	"github.com/walteh/webauthn/pkg/otel/logging"
	"github.com/walteh/webauthn/pkg/otel/otelawsc"
	"github.com/walteh/webauthn/pkg/rhttp"
	lambdadetector "go.opentelemetry.io/contrib/detectors/aws/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Options struct {

	// TracerProvider is the TracerProvider which will be used
	// to create instrumentation spans
	// The default value of TracerProvider the global otel TracerProvider
	// returned by otel.GetTracerProvider()
	tracerProvider trace.TracerProvider

	propogator propagation.TextMapPropagator

	// Flusher is the mechanism used to flush any unexported spans
	// each Lambda Invocation to avoid spans being unexported for long
	// when periods of time if Lambda freezes the execution environment
	// The default value of Flusher is a noop Flusher, using this
	// default can result in long data delays in asynchronous settings
	logger *zerolog.Logger

	awsConfig *aws.Config

	httpClient *http.Client

	_retryRoundTripper http.RoundTripper
	_otelRoundTripper  bool

	postapply []OptFunc
}

func newDeadOptions() *Options {
	return &Options{
		postapply:      []OptFunc{},
		tracerProvider: otel.GetTracerProvider(),
	}
}

func (options *Options) TracerProvider() trace.TracerProvider {
	if options.tracerProvider == nil {
		options.tracerProvider = otel.GetTracerProvider()
	}
	return options.tracerProvider
}

func (options *Options) Logger(ctx context.Context) *zerolog.Logger {
	if options.logger == nil {
		options.logger = zerolog.DefaultContextLogger
	}

	ctxr := options.logger.With()

	trace.SpanContextFromContext(ctx)

	b := &propagation.MapCarrier{}

	options.Propagator().Inject(ctx, b)

	for k, v := range *b {
		ctxr = ctxr.Str(k, v)
	}

	lgr := ctxr.Logger()

	return &lgr
}

func (options *Options) AwsConfig() *aws.Config {
	if options.awsConfig == nil {
		options.awsConfig = aws.NewConfig()
	}
	return options.awsConfig
}

func (options *Options) HttpClient() *http.Client {
	if options.httpClient == nil {
		options.httpClient = &http.Client{
			Transport: options.RoundTripper(),
		}
	}
	return options.httpClient
}

func (options *Options) WithVerboseLogger() {
	options.logger = logging.NewVerboseLogger()
}

func (options *Options) WithVerboseLoggerWithLevel(level zerolog.Level) {
	lev := logging.NewVerboseLogger().Level(level)
	options.logger = &lev
}

func (options *Options) WithDefaultLogger() {
	options.logger = logging.NewJsonLogger()
}

func (options *Options) WithLogLevel(level zerolog.Level) {
	options.postapply = append(options.postapply, func(_ context.Context, opts *Options) {
		if opts.logger != nil {
			lev := opts.logger.Level(level)
			opts.logger = &lev
		}
	})
}

func (options *Options) WithLogger(logger *zerolog.Logger) {
	options.logger = logger
}

func (options *Options) WithAwsConfig(config *aws.Config) {
	options.awsConfig = config
}

func (options *Options) WithLocalAwsConfig() {
	options.awsConfig = aws.NewConfig()
	options.awsConfig.Credentials = credentials.NewStaticCredentialsProvider("fake", "fake", "")
	options.awsConfig.Region = "us-east-1"
}

func (options *Options) WithDefaultAwsConfig(ctx context.Context) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	options.awsConfig = &cfg
}

func (me *Options) WithOtelAws() {
	me.postapply = append(me.postapply, func(_ context.Context, opts *Options) {
		apiopts := &opts.AwsConfig().APIOptions
		otelaws.AppendMiddlewares(
			apiopts,
			otelaws.WithTracerProvider(opts.TracerProvider()),
			otelaws.WithAttributeSetter(otelawsc.DefaultAttributeSetter),
			otelaws.WithTextMapPropagator(opts.Propagator()),
		)
		*apiopts = append(*apiopts, otelawsc.DeserializeHttpMiddleware)
		// *apiopts = append(*apiopts, otelawsc.InitializeHttpMiddleware)
	})
}

func (me *Options) WithHttpClientAppliedToAwsConfig(ctx context.Context) {
	me.postapply = append(me.postapply, func(_ context.Context, opts *Options) {
		opts.AwsConfig().HTTPClient = &http.Client{
			Transport: opts.RoundTripperNoHttpOtel(),
		}
	})
}

func (options *Options) RoundTripperNoHttpOtel() http.RoundTripper {
	if options._retryRoundTripper == nil && !options._otelRoundTripper {
		return http.DefaultTransport
	} else {
		var trip http.RoundTripper
		if options._retryRoundTripper != nil {
			trip = options._retryRoundTripper
		} else {
			trip = http.DefaultTransport
		}
		return trip
	}
}

func (options *Options) RoundTripper() http.RoundTripper {
	if options._retryRoundTripper == nil && !options._otelRoundTripper {
		return http.DefaultTransport
	} else {
		var trip http.RoundTripper
		if options._retryRoundTripper != nil {
			trip = options._retryRoundTripper
		} else {
			trip = http.DefaultTransport
		}
		if options._otelRoundTripper {
			trip = otelhttp.NewTransport(
				trip,
				otelhttp.WithPropagators(options.Propagator()),
				otelhttp.WithTracerProvider(options.TracerProvider()),
			)
		}
		return trip
	}
}

func (options *Options) WithRetryRoundTripper(client *rhttp.RoundTripper) {
	options._retryRoundTripper = client
}

func (me *Options) WithOtelRoundTripper() {
	me._otelRoundTripper = true
}

func (me *Options) WithPropagator(prop propagation.TextMapPropagator) {
	otel.SetTextMapPropagator(prop)
	me.propogator = prop
	// me.postapply = append(me.postapply, func(_ context.Context, opts *Options) {
	// 	otelaws.AppendMiddlewares(&opts.AwsConfig().APIOptions, otelaws.WithTextMapPropagator(prop))
	// })
}

func (options *Options) WithTracerProvider(tracerProvider trace.TracerProvider) {
	otel.SetTracerProvider(tracerProvider)
	options.tracerProvider = tracerProvider

	// options.postapply = append(options.postapply, func(_ context.Context, opts *Options) {
	// 	otelaws.AppendMiddlewares(&opts.AwsConfig().APIOptions, otelaws.WithTracerProvider(tracerProvider))
	// })

}

func (me *Options) Propagator() propagation.TextMapPropagator {
	if me.propogator == nil {
		me.propogator = otel.GetTextMapPropagator()
	}
	return me.propogator
}

func (me *Options) WithDefaultTracing(ctx context.Context, name string, endpoint string) {

	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		// otlptracegrpc.WithDialOption(grpc.WithTr()),
	)

	if err != nil {
		panic(err)
	}

	idg := xray.NewIDGenerator()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		// the service name used to display traces in backends
		semconv.ServiceNameKey.String(name),
	)
	// handleErr(err, "failed to create resource")

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithIDGenerator(idg),
	)

	me.WithTracerProvider(tp)
	me.WithPropagator(xray.Propagator{})

	me.postapply = append(me.postapply, func(ctx context.Context, opts *Options) {
		*opts.logger = opts.Logger(ctx).With().Str("service", name).Logger()
	})
}

func (me *Options) WithLambdaTracing(ctx context.Context) {
	idg := xray.NewIDGenerator()

	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
		// otlptracegrpc.WithDialOption(grpc.WithTr()),
	)

	if err != nil {
		panic(err)
	}

	name := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	var res *resource.Resource
	res, err = lambdadetector.NewResourceDetector().Detect(ctx)
	if err != nil {
		res = resource.NewWithAttributes(
			semconv.SchemaURL,
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(name),
		)
	}
	// handleErr(err, "failed to create resource")

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithSampler(sdktrace.ParentBased(
			sdktrace.AlwaysSample(),
		)),
		sdktrace.WithResource(res),
		sdktrace.WithIDGenerator(idg),
	)

	me.WithTracerProvider(tp)
	me.WithPropagator(xray.Propagator{})

	me.postapply = append(me.postapply, func(ctx context.Context, opts *Options) {
		*opts.logger = opts.Logger(ctx).With().Str("service", name).Logger()
	})
}
