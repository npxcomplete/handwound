package handwound

import "time"

/**
 * See https://pkg.go.dev/time
 */
type SystemClock struct{}

var _ Clock = SystemClock{}

func (c SystemClock) Now() time.Time {
	return time.Now()
}

func (c SystemClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (c SystemClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (c SystemClock) AfterFunc(d time.Duration, f func()) Timer {
	return SystemTimer{time.AfterFunc(d, f)}
}

func (c SystemClock) NewTimer(d time.Duration) Timer {
	return SystemTimer{time.NewTimer(d)}
}

type SystemTimer struct {
	*time.Timer
}

var _ Timer = SystemTimer{}

func (t SystemTimer) C() <-chan time.Time {
	return t.Timer.C
}

func (t SystemTimer) Reset(d time.Duration) bool {
	return t.Timer.Reset(d)
}

func (t SystemTimer) Stop() bool {
	return t.Timer.Stop()
}
