package safeid

import "time"

// Time returns the Unix time in milliseconds encoded in the SafeID.
// Use the top level Time function to convert the returned value to
// a time.Time.
func (id SafeID) Time() uint64 {
	return uint64(id[5]) | uint64(id[4])<<8 |
		uint64(id[3])<<16 | uint64(id[2])<<24 |
		uint64(id[1])<<32 | uint64(id[0])<<40
}

// maxTime is the maximum Unix time in milliseconds that can be
// represented in an SafeID.
var maxTime = SafeID{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}.Time()

// MaxTime returns the maximum Unix time in milliseconds that
// can be encoded in an SafeID.
func MaxTime() uint64 { return maxTime }

// Now is a convenience function that returns the current
// UTC time in Unix milliseconds. Equivalent to:
//
//	Timestamp(time.Now().UTC())
func Now() uint64 { return Timestamp(time.Now().UTC()) }

// Timestamp converts a time.Time to Unix milliseconds.
//
// Because of the way SafeID stores time, times from the year
// 10889 produces undefined results.
func Timestamp(t time.Time) uint64 {
	return uint64(t.Unix())*1000 +
		uint64(t.Nanosecond()/int(time.Millisecond))
}

// Time converts Unix milliseconds in the format
// returned by the Timestamp function to a time.Time.
func Time(ms uint64) time.Time {
	s := int64(ms / 1e3)
	ns := int64((ms % 1e3) * 1e6)
	return time.Unix(s, ns)
}

// SetTime sets the time component of the SafeID to the given Unix time
// in milliseconds.
func (id *SafeID) SetTime(ms uint64) error {
	if ms > maxTime {
		return ErrBigTime
	}

	(*id)[0] = byte(ms >> 40)
	(*id)[1] = byte(ms >> 32)
	(*id)[2] = byte(ms >> 24)
	(*id)[3] = byte(ms >> 16)
	(*id)[4] = byte(ms >> 8)
	(*id)[5] = byte(ms)

	return nil
}
