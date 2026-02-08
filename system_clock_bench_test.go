package handwound

import (
	"testing"
	"time"
)

func BenchmarkRawClockNow(b *testing.B) {
	for b.Loop() {
		_ = time.Now()
	}
}

func BenchmarkSystemClockNow(b *testing.B) {
	var clock Clock = SystemClock{}
	for b.Loop() {
		_ = clock.Now()
	}
}
