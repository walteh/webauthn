package invocation

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type Invokable[I any, O any] interface {
	Invoke(ctx context.Context, input I) (O, error)
}

type Tracable interface {
	Opts() *Options
}

type Handler[I any, O any] interface {
	ID() string
	IncrementCounter() int
	Cleanup()
	Recover()
	Cancel()
}

type baseHandler[I any, O any] struct {
	id         string
	start      time.Time
	counter    int
	context    context.Context
	cancelFunc context.CancelFunc
	options    *Options
	invoker    Invokable[I, O]
	cleanups   []func(*baseHandler[I, O])
}

func (h *baseHandler[I, O]) AddCleanup(cleanup func(*baseHandler[I, O])) {
	h.cleanups = append(h.cleanups, cleanup)
}

func (h *baseHandler[I, O]) ID() string {
	return h.id
}

func (h *baseHandler[I, O]) IncrementCounter() int {
	h.counter += 1
	return h.counter
}

func (h *baseHandler[I, O]) Cleanup() {
	zerolog.Ctx(h.Ctx()).Info().Int("invocations", h.counter).Dur("duration", time.Since(h.start)).Msgf("graceful shutdown of handler %s", h.id)
	for _, cleanup := range h.cleanups {
		cleanup(h)
	}
	h.Cancel()
}

func (h *baseHandler[I, O]) Opts() *Options {
	return h.options
}

func (h *baseHandler[I, O]) Ctx() context.Context {
	return h.context
}

func (h *baseHandler[I, O]) Cancel() {
	h.cancelFunc()
}

type OptFunc func(context.Context, *Options)

func (h *baseHandler[I, O]) Invoke(ctx context.Context, input I) (O, error) {
	ctx, cnl := context.WithCancel(ctx)
	defer cnl()

	counter := h.IncrementCounter()
	handlerID := h.ID()
	span := trace.SpanFromContext(ctx)

	log := zerolog.Ctx(h.context).With().Int("counter", counter).Str("handler", handlerID).Logger()

	ctx = log.WithContext(ctx)

	log.Trace().CallerSkipFrame(1).Interface("input", input).Msg("input parsed")

	log.Info().Msg("invocation started")

	defer func() {
		log.Info().Msg("invocation finished")
		h.Cleanup()
		if r := recover(); r != nil {
			log.Warn().Interface("panic", r).Msg("invocation panicked")
		}
	}()

	return h.invoker.Invoke(ctx, input)
}

func NewBaseHandler[I any, O any](ctx context.Context, opt ...OptFunc) Handler[I, O] {

	handler := &baseHandler[I, O]{}

	handler.start = time.Now()

	handler.id = xid.New().String()

	handler.counter = 0

	ctx, handler.cancelFunc = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGSEGV, syscall.SIGABRT)

	handler.options = newDeadOptions()

	for _, pre := range opt {
		pre(handler.context, handler.options)
	}

	for _, post := range handler.options.postapply {
		post(handler.context, handler.options)
	}

	return handler
}

func NewDefaultTestbaseHandler[I any, O any]() Handler[I, O] {
	return NewBaseHandler[I, O](context.Background(), func(ctx context.Context, o *Options) {
		o.WithLocalAwsConfig()
		// o.WithVerboseLoggerWithLevel(zerolog.TraceLevel)
		o.WithHttpClientAppliedToAwsConfig(ctx)
		o.WithLogLevel(zerolog.DebugLevel)
	})
}
