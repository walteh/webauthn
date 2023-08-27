package clientdata

import (
	"errors"
	"reflect"
)

var errref struct{}

var (
	ErrInvalidCeremonyType    = errors.New(reflect.TypeOf(errref).PkgPath() + ":ErrInvalidCeremonyType")
	ErrChallengeMismatch      = errors.New(reflect.TypeOf(errref).PkgPath() + ":ErrChallengeMismatch")
	ErrOriginMismatch         = errors.New(reflect.TypeOf(errref).PkgPath() + ":ErrOriginMismatch")
	ErrTokenMissingStatus     = errors.New(reflect.TypeOf(errref).PkgPath() + ":ErrTokenMissingStatus")
	ErrTokenInvalidStatus     = errors.New(reflect.TypeOf(errref).PkgPath() + ":ErrTokenInvalidStatus")
	ErrOriginNotParsableAsURL = errors.New(reflect.TypeOf(errref).PkgPath() + ":ErrOriginNotParsableAsURL")
)
