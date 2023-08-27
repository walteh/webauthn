package invocation

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/walteh/webauthn/pkg/rhttp"
	"golang.org/x/sync/errgroup"
)

type Invokable[I any, O any] interface {
	Invoke(ctx context.Context, input I) (O, error)
}

type Tracable interface {
	Opts() *Options
}

type Handler[I any, O any] struct {
	id         string
	start      time.Time
	counter    int
	context    context.Context
	cancelFunc context.CancelFunc
	errorGroup *errgroup.Group
	options    *Options
	cleanups   []func(*Handler[I, O])
}

func (h *Handler[I, O]) AddCleanup(cleanup func(*Handler[I, O])) {
	h.cleanups = append(h.cleanups, cleanup)
}

func (h *Handler[I, O]) ID() string {
	return h.id
}

func (h *Handler[I, O]) IncrementCounter() int {
	h.counter += 1
	return h.counter
}

func (h *Handler[I, O]) Cleanup() {
	h.Opts().Logger(h.Ctx()).Info().Int("invocations", h.counter).Dur("duration", time.Since(h.start)).Msgf("graceful shutdown of handler %s", h.id)
	for _, cleanup := range h.cleanups {
		cleanup(h)
	}
	h.Cancel()
}

func (h *Handler[I, O]) Recover() {
	if r := recover(); r != nil {
		h.Opts().Logger(h.Ctx()).Warn().Msgf("Recovered from panic: %v", r)
	}
}

func (h *Handler[I, O]) Opts() *Options {
	return h.options
}

func (h *Handler[I, O]) Ctx() context.Context {
	return h.context
}

func (h *Handler[I, O]) Cancel() {
	h.cancelFunc()
}

type OptFunc func(context.Context, *Options)

func NewHandler[I any, O any](ctx context.Context, opt ...OptFunc) *Handler[I, O] {

	handler := &Handler[I, O]{}

	handler.start = time.Now()

	handler.id = xid.New().String()

	handler.counter = 0

	ctx, handler.cancelFunc = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGSEGV, syscall.SIGABRT)

	handler.errorGroup, handler.context = errgroup.WithContext(ctx)

	handler.options = newDeadOptions()

	for _, pre := range opt {
		pre(handler.context, handler.options)
	}

	for _, post := range handler.options.postapply {
		post(handler.context, handler.options)
	}

	return handler
}

func NewDefaultTestHandler[I any, O any]() *Handler[I, O] {
	return NewHandler[I, O](context.Background(), func(ctx context.Context, o *Options) {
		o.WithLocalAwsConfig()
		o.WithVerboseLoggerWithLevel(zerolog.TraceLevel)
		o.WithHttpClientAppliedToAwsConfig(ctx)

		cli := rhttp.NewClientContext(ctx)
		cli.RetryMax = 0
		o.WithRetryRoundTripper(cli.Transport())
		o.WithLogLevel(zerolog.DebugLevel)
	})
}
