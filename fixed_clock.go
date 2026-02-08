package handwound

import (
	"sync"
	"time"
)

type FixedClock struct {
	mux      sync.Mutex
	currentT time.Time
	timers   []*ReactiveTimer
}

func NewFixedClock(initialT time.Time) *FixedClock {
	return &FixedClock{
		currentT: initialT,
		timers:   make([]*ReactiveTimer, 0, 8),
	}
}

func (c *FixedClock) Advance(d time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.currentT = c.currentT.Add(d)
	for _, timer := range c.timers {
		go func() {
			timer.mux.Lock()
			defer timer.mux.Unlock()

			timer.trigger()
		}()
	}
}

func (timer *ReactiveTimer) trigger() {
	if timer.stopped {
		return
	}
	if timer.expireAt.Before(timer.clock.Now()) {
		timer.stopped = true
		timer.localTime = timer.expireAt
		go func(t time.Time) { timer.channel <- t }(timer.expireAt)
	}
}

var _ Clock = &FixedClock{}

func (c *FixedClock) Now() time.Time {
	return c.currentT
}

func (c *FixedClock) Sleep(d time.Duration) {
	<-c.NewTimer(d).C()
}

func (c *FixedClock) NewTimer(d time.Duration) Timer {
	c.mux.Lock()
	defer c.mux.Unlock()
	now := c.Now()
	r := &ReactiveTimer{
		mux:       sync.Mutex{},
		clock:     c,
		stopped:   false,
		expireAt:  now.Add(d),
		localTime: now,
		channel:   make(chan time.Time),
	}
	c.timers = append(c.timers, r)
	return r
}

func (c *FixedClock) After(d time.Duration) <-chan time.Time {
	return c.NewTimer(d).C()
}

func (c *FixedClock) AfterFunc(d time.Duration, f func()) Timer {
	timer := c.NewTimer(d)
	go func() {
		for {
			_ = <-timer.C()
			f()
		}
	}()
	return timer
}

type ReactiveTimer struct {
	mux       sync.Mutex
	clock     *FixedClock
	expireAt  time.Time
	localTime time.Time
	stopped   bool
	channel   chan time.Time
}

var _ Timer = &ReactiveTimer{}

func (r *ReactiveTimer) C() <-chan time.Time {
	return r.channel
}

func (r *ReactiveTimer) Reset(d time.Duration) bool {
	r.mux.Lock()
	defer r.mux.Unlock()
	active := !r.stopped && r.localTime.Before(r.expireAt)

	r.stopped = false
	r.expireAt = r.localTime.Add(d)

	r.trigger()
	return active
}

func (r *ReactiveTimer) Stop() bool {
	r.mux.Lock()
	defer r.mux.Unlock()
	if !r.stopped {
		r.localTime = r.clock.Now()
		r.stopped = true
		return true
	} else {
		return false
	}
}

type ProxyTimer struct {
	inner *ReactiveTimer
}

var _ Timer = ProxyTimer{}

func (t ProxyTimer) C() <-chan time.Time {
	return nil
}

func (t ProxyTimer) Reset(d time.Duration) bool {
	return t.inner.Reset(d)
}

func (t ProxyTimer) Stop() bool {
	return t.inner.Stop()
}
