package safeid

import "errors"

var (
	// ErrDataSize is returned when parsing or unmarshaling ULIDs with the wrong
	// data size.
	ErrDataSize = errors.New("ulid: bad data size when unmarshaling")

	// ErrInvalidCharacters is returned when parsing or unmarshaling ULIDs with
	// invalid Base32 encodings.
	ErrInvalidCharacters = errors.New("ulid: bad data characters when unmarshaling")

	// ErrBufferSize is returned when marshalling ULIDs to a buffer of insufficient
	// size.
	ErrBufferSize = errors.New("ulid: bad buffer size when marshaling")

	// ErrBigTime is returned when constructing an ULID with a time that is larger
	// than MaxTime.
	ErrBigTime = errors.New("ulid: time too big")

	// ErrOverflow is returned when unmarshaling a ULID whose first character is
	// larger than 7, thereby exceeding the valid bit depth of 128.
	ErrOverflow = errors.New("ulid: overflow when unmarshaling")

	// ErrMonotonicOverflow is returned by a Monotonic entropy source when
	// incrementing the previous ULID's entropy bytes would result in overflow.
	ErrMonotonicOverflow = errors.New("ulid: monotonic entropy overflow")

	// ErrScanValue is returned when the value passed to scan cannot be unmarshaled
	// into the ULID.
	ErrScanValue = errors.New("ulid: source value must be a string or byte slice")
)
