package dynamo

import "errors"

var ErrConditionalCheckFailed = errors.New("conditional check failed")
var ErrNotFound = errors.New("not found")

var ErrCeremonyExpired = errors.New("ceremony_expired")
var ErrDynamoConnectionError = errors.New("dynamo_connection_error")
