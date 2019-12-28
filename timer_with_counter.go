package epimetheus

import (
	"time"

	"github.com/cactus/go-statsd-client/statsd"
)

// TimerWithCounter consists of a Timer and a Counter in order to measure count and duration simultaneously
type TimerWithCounter struct {
	*Timer
	*Counter
}

func newTimerWithCounter(namespace, subsystem, name string, labelNames []string, client *statsd.Statter) *TimerWithCounter {
	timer := newTimer(namespace, subsystem, name+"Timer", labelNames, client)
	counter := newCounter(namespace, subsystem, name+"Counter", labelNames, client)
	return &TimerWithCounter{
		timer,
		counter,
	}
}

// Start creates an instance of `RunningTimerWithCounter` and returns it
func (tc *TimerWithCounter) Start() time.Time {
	return tc.Timer.Start()
}

// Done marks the related timer as done and increments the related counter too
func (tc *TimerWithCounter) Done(start time.Time, labelValues ...string) {
	tc.Timer.Done(start, labelValues...)
	tc.Counter.Inc(labelValues...)
}
