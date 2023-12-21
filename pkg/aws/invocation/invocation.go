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
	Cancel()
	Error(err error, extra ...interface{}) (O, error)
	Success(output O, extra ...interface{}) (O, error)
	HandledError(output O, err error, extra ...interface{}) (O, error)
}

type BaseInvocation[I any, O any] struct {
	start     time.Time
	cancel    context.CancelFunc
	handlerId string
	counter   int
	span      trace.Span
	// labels    []string
}

func (h *BaseInvocation[I, O]) Cancel() {
	h.cancel()
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

func NewInvocation[I any, O any](ctx context.Context, me Handler[I, O], input I, labels ...string) (Invocation[I, O], context.Context) {
	ctx, cnl := context.WithCancel(ctx)

	counter := me.IncrementCounter()
	handlerID := me.ID()
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(attribute.StringSlice("tags", labels))

	log := zerolog.Ctx(ctx).With().Int("counter", counter).Array("tags", zerolabels(labels)).Str("handler", handlerID).Logger()

	ctx = log.WithContext(ctx)

	log.Trace().CallerSkipFrame(1).Interface("input", input).Msg("input parsed")

	log.Info().Msg("invocation started")

	return &BaseInvocation[I, O]{
		cancel:    cnl,
		start:     time.Now(),
		handlerId: handlerID,
		counter:   counter,
		span:      span,
	}, ctx
}
