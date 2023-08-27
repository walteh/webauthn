package invocation

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Invocation[I any, O any] interface {
	Start() time.Time
	Logger() *zerolog.Logger
	Cancel()
	Error(err error, extra ...interface{}) (O, error)
	Success(output O, extra ...interface{}) (O, error)
	HandledError(output O, err error, extra ...interface{}) (O, error)
	ID() string
}

type BaseInvocation[I any, O any] struct {
	Invocation[I, O]
	logger    zerolog.Logger
	start     time.Time
	cancel    context.CancelFunc
	handlerId string
	counter   int
	span      trace.Span
	// labels    []string
}

func (h *BaseInvocation[I, O]) SetSpan(span trace.Span) {
	h.span = span
}

func (h *BaseInvocation[I, O]) Cancel() {
	h.cancel()
}

func (h *BaseInvocation[I, O]) Logger() *zerolog.Logger {
	return &h.logger
}

func (h *BaseInvocation[I, O]) Start() time.Time {
	return h.start
}

func (h *BaseInvocation[I, O]) ID() string {
	return h.handlerId + "-" + strconv.Itoa(h.counter)
}

type zerolabels []string

func (l zerolabels) MarshalZerologArray(e *zerolog.Array) {
	for _, s := range l {
		e.Str(s)
	}
}

func (me *Handler[C, G]) NewInvocation(ctx context.Context, input C, labels ...string) (Invocation[C, G], context.Context) {
	ctx, cnl := context.WithCancel(ctx)

	counter := me.IncrementCounter()
	handlerID := me.ID()
	// ctx, span := me.options.TracerProvider().Tracer("nugg.xyz/go/invocation").Start(ctx, "invocation")
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(attribute.StringSlice("tags", labels))

	log := me.Opts().Logger(ctx).With().Int("counter", counter).Array("tags", zerolabels(labels)).Str("handler", handlerID).Logger()

	ctx = log.WithContext(ctx)

	inputlog := log.Trace().CallerSkipFrame(1)

	if g, ok := (interface{})(input).(zerolog.LogArrayMarshaler); ok {
		inputlog.Array("input", g)
	} else if g, ok := (interface{})(input).(zerolog.LogObjectMarshaler); ok {
		inputlog.Object("input", g)
	} else {
		inputlog.Interface("input", input)
	}

	inputlog.Msg("parsed input")

	log.Info().Msg("invocation started")

	return &BaseInvocation[C, G]{
		cancel:    cnl,
		logger:    log,
		start:     time.Now(),
		handlerId: handlerID,
		counter:   counter,
		span:      span,
	}, ctx
}
