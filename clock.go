package handwound

import "time"

/**
 * The implemenation of SystemClock is a thin wrapper over the
 * go SDK's time package. A clock must hold to the same contract.
 * however, as other implemtations are primarily for test purposes
 * they have no obligation to respect "real" time.
 *
 * See https://pkg.go.dev/time
 */
type Clock interface {
	Now() time.Time

	Sleep(d time.Duration)

	After(d time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) Timer
	NewTimer(d time.Duration) Timer
}

type Timer interface {
	C() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}
