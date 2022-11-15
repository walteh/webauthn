package safeid

import (
	"io"
)

type SafeID [16]byte

// Make returns an SafeID with the current time in Unix milliseconds and
// monotonically increasing entropy for the same millisecond.
// It is safe for concurrent use, leveraging a sync.Pool underneath for minimal
// contention.
func Make() (id SafeID) {
	// NOTE: MustNew can't panic since DefaultEntropy never returns an error.
	return MustNew(Now(), DefaultEntropy())
}

// New returns an SafeID with the given Unix milliseconds timestamp and an
// optional entropy source. Use the Timestamp function to convert
// a time.Time to Unix milliseconds.
//
// ErrBigTime is returned when passing a timestamp bigger than MaxTime.
// Reading from the entropy source may also return an error.
//
// Safety for concurrent use is only dependent on the safety of the
// entropy source.
func New(ms uint64, entropy io.Reader) (id SafeID, err error) {
	if err = id.SetTime(ms); err != nil {
		return id, err
	}

	switch e := entropy.(type) {
	case nil:
		return id, err
	case MonotonicReader:
		err = e.MonotonicRead(ms, id[6:])
	default:
		_, err = io.ReadFull(e, id[6:])
	}

	return id, err
}

// MustNew is a convenience function equivalent to New that panics on failure
// instead of returning an error.
func MustNew(ms uint64, entropy io.Reader) SafeID {
	id, err := New(ms, entropy)
	if err != nil {
		panic(err)
	}
	return id
}
