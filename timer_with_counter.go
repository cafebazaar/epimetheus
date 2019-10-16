package epimetheus

import (
	"github.com/cactus/go-statsd-client/statsd"
)

// TimerWithCounter consists of a Timer and a Counter in order to measure count and duration simultaneously
type TimerWithCounter struct {
	*Timer
	*Counter
}

// RunningTimerWithCounter consists of a RunningTimer and a TimerWithCounter
//
// Calling `Done` on the instance which returned by `TimerWithCounter.Start` finalize the operation
type RunningTimerWithCounter struct {
	runningTimer     *RunningTimer
	timerWithCounter *TimerWithCounter
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
func (w *TimerWithCounter) Start() *RunningTimerWithCounter {
	return &RunningTimerWithCounter{
		runningTimer:     w.Timer.Start(),
		timerWithCounter: w,
	}
}

// Done marks the related timer as done and increments the related counter too
func (rt *RunningTimerWithCounter) Done(labelValues ...string) {
	rt.runningTimer.Done(labelValues...)
	rt.timerWithCounter.Counter.Inc(labelValues...)
}
