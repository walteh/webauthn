package types

import (
	"errors"
	"reflect"
)

var errref struct{}

var (
	ErrUnmarshaling          = errors.New(reflect.TypeOf(errref).PkgPath() + ":ErrUnmarshaling")
	ErrInvalidAssertionInput = errors.New(reflect.TypeOf(errref).PkgPath() + ":ErrInvalidAssertionInput")
)
