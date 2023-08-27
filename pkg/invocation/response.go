package invocation

import (
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/codes"
)

// func NewApigwV2Invocation[I ApigwV2Input, O ApigwV2Output](ctx context.Context, handler Handler[I, O], input I) (Invocation[I, O], context.Context) {
// 	bi, ctx := Build[I, O](ctx, handler, input)
// 	return &bi, ctx
// }

func (h *BaseInvocation[I, O]) Error(err error, extra ...interface{}) (O, error) {

	if err == nil {
		err = errors.New("unknown error")
	}

	event := h.Logger().Error()

	if err == nil {
		err = errors.New("unspecified error")
	}

	for i, e := range extra {
		event = event.Interface(fmt.Sprintf("extra[%d]", i), e)
	}

	event.
		Err(err).
		CallerSkipFrame(1).
		TimeDiff("duration", time.Now(), h.Start()).
		Msg("returning error")

	h.span.SetStatus(codes.Error, err.Error())

	go h.Cancel()

	var z O

	return z, err

}

func (h *BaseInvocation[I, O]) HandledError(output O, err error, extra ...interface{}) (O, error) {

	event := h.Logger().Error()

	event.
		Err(err).
		Interface("output", output).
		CallerSkipFrame(1).
		TimeDiff("duration", time.Now(), h.Start()).
		Msg("handling error")

	h.span.SetStatus(codes.Error, err.Error())

	go h.Cancel()

	return output, nil

}

func (h *BaseInvocation[I, O]) Success(output O, extra ...interface{}) (O, error) {

	event := h.Logger().Info()

	for i, e := range extra {
		event = event.Interface(fmt.Sprintf("extra[%d]", i), e)
	}

	event.
		Interface("output", output).
		CallerSkipFrame(1).
		TimeDiff("duration", time.Now(), h.Start()).
		Msg("returning success")

	h.span.SetStatus(codes.Ok, "success")

	go h.Cancel()

	return output, nil
}
