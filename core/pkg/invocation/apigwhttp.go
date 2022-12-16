package invocation

import (
	"context"
	"time"

	"github.com/nuggxyz/webauthn/pkg/errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Invocation struct {
	zerolog.Logger
	Start  time.Time
	cancel context.CancelFunc
}

func NewInvocation(ctx context.Context, handler Handler, input Input) (*Invocation, context.Context) {
	ctx, cnl := context.WithCancel(ctx)

	logg := handler.Logger().With().
		Int("counter", handler.IncrementCounter()).
		Str("handler", handler.ID()).
		Logger()

	logg.Info().Interface("input", input).Msg("starting invocation")

	ctx = logg.WithContext(ctx)

	return &Invocation{
		cancel: cnl,
		Logger: logg,
		Start:  time.Now(),
	}, ctx
}

func (h *Invocation) Error(err error, code int, message string) (Output, error) {

	event := h.Logger.Error()

	// if error is an error type, then we can use the With* methods
	if e, ok := err.(*errors.Error); ok {
		err = e.Roots()[0]
		roots := e.Roots()
		event = event.Errs("roots", roots[1:])

	}

	event.
		Err(err).
		Str("body", "").
		CallerSkipFrame(1).
		TimeDiff("duration", time.Now(), h.Start).
		Msg("returning error: " + message)

	go h.cancel()

	return Output{StatusCode: code}, nil
}

func (h *Invocation) Success(code int, headers map[string]string, message string) (Output, error) {

	output := Output{
		StatusCode: code,
		Headers:    headers,
	}

	if message != "" && code != 204 {
		output.Body = message
	}

	r := zerolog.Dict()
	for k, v := range output.Headers {
		r = r.Str(k, v)
	}

	if message == "" {
		message = "empty"
	}

	if code == 204 && headers["Content-Length"] == "" {
		output.Headers["Content-Length"] = "0"
	}

	h.Logger.Info().
		Int("status_code", code).
		Str("body", output.Body).
		Dict("headers", r).
		CallerSkipFrame(1).
		TimeDiff("duration", time.Now(), h.Start).
		Msg(message)

	go h.cancel()

	return output, nil
}
