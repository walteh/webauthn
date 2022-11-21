package invocation

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
)

type Input = events.APIGatewayV2HTTPRequest
type Output = events.APIGatewayV2HTTPResponse

type Invocation struct {
	zerolog.Logger
	Start  time.Time
	Ctx    context.Context
	cancel context.CancelFunc
}

func NewInvocation(ctx context.Context, handler Handler, input Input) *Invocation {
	ctx, cnl := context.WithCancel(ctx)

	logg := handler.Logger().With().
		Int("counter", handler.IncrementCounter()).
		Str("handler", handler.ID()).
		Logger()

	logg.Info().Interface("input", input).Msg("starting invocation")

	ctx = logg.WithContext(ctx)

	return &Invocation{
		cancel: cnl,
		Ctx:    ctx,
		Logger: logg,
		Start:  time.Now(),
	}
}

func (h *Invocation) Error(err error, code int, message string) (Output, error) {
	h.Logger.Error().Err(err).
		Int("status_code", code).
		Str("body", "").
		CallerSkipFrame(2).
		TimeDiff("duration", time.Now(), h.Start).
		Msg(message)

	go h.cancel()

	return Output{
		StatusCode: code,
	}, nil
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
