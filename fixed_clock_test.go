package handwound

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_FixedClock(t *testing.T) {
	t.Run("Does not move unbid", func(t *testing.T) {
		clock := NewFixedClock(time.Now())
		t1 := clock.Now()
		<-time.After(time.Second)
		t2 := clock.Now()
		require.Equal(t, t1, t2)
	})

	t.Run("Moves when bid", func(t *testing.T) {
		clock := NewFixedClock(time.Now())
		timer := clock.NewTimer(1 * time.Millisecond)
		wg := sync.WaitGroup{}
		done := false
		wg.Add(1)
		go func() {
			<-timer.C()
			done = true
			wg.Done()
		}()
		<-time.After(100 * time.Millisecond)
		require.False(t, done)
		clock.Advance(2 * time.Millisecond)
		wg.Wait()
		require.True(t, done)
	})

	t.Run("Won't refire unbid", func(t *testing.T) {

		clock := NewFixedClock(time.Now())
		timer := clock.NewTimer(1 * time.Millisecond)

		wg := sync.WaitGroup{}
		done := false
		wg.Add(1)
		go func() {
			<-timer.C()
			select {
			case <-timer.C():
			case <-time.After(100 * time.Millisecond):
				done = true
			}
			wg.Done()
		}()

		clock.Advance(2 * time.Millisecond)
		clock.Advance(2 * time.Millisecond)

		wg.Wait()
		require.True(t, done)
	})

	t.Run("Re-fires after reset", func(t *testing.T) {
		clock := NewFixedClock(time.Now())
		timer := clock.NewTimer(1 * time.Millisecond)

		wg := &sync.WaitGroup{}
		done := false
		wg.Add(1)
		go func() {
			<-timer.C()
			timer.Reset(1 * time.Millisecond)
			<-timer.C()
			done = true
			wg.Done()
		}()

		clock.Advance(2 * time.Millisecond)
		require.False(t, done)
		clock.Advance(2 * time.Millisecond)
		wg.Wait()
		require.True(t, done)
	})

	t.Run("Reset will re-fire if the clock is in the far future.", func(t *testing.T) {
		clock := NewFixedClock(time.Now())
		timer := clock.NewTimer(1 * time.Millisecond)

		wg := &sync.WaitGroup{}
		done := false
		wg.Add(1)
		go func() {
			for i := 10; i > 0; i-- {
				<-timer.C()
				timer.Reset(1 * time.Millisecond)
				log.Printf("triggers remaining %d\n", i-1)
			}
			done = true
			wg.Done()
		}()
		clock.Advance(20 * time.Millisecond)
		wg.Wait()
		require.True(t, done)
	})
}
