package invocation

import "github.com/rs/zerolog"

type Handler interface {
	IncrementCounter() int
	ID() string
	Logger() zerolog.Logger
}
