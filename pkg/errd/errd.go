package errd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

func Wrap(ctx context.Context, e error, s ...string) error {
	event := zerolog.Ctx(ctx).Error().Err(e).CallerSkipFrame(1)
	for i, msg := range s {
		event = event.Str(fmt.Sprintf("extra[%d]", i), msg)
	}
	zerolog.Ctx(ctx).Error().Err(e).CallerSkipFrame(1).Msg("error")
	return e
}

func Mismatch(ctx context.Context, e error, expected, actual string) error {
	zerolog.Ctx(ctx).Error().Err(e).Str("expected", expected).Str("actual", actual).CallerSkipFrame(1).Msg("mismatch")
	return e
}
