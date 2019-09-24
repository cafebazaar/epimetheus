package epimetheus

import (
	"github.com/cactus/go-statsd-client/statsd"
)

type TimerWithCounter struct {
	*Timer
	*Counter
}

type RunningTimerWithCounter struct {
	runningTimer *RunningTimer
	timerWithCounter *TimerWithCounter
}

func NewTimerWithCounter(namespace, subsystem, name string, labelNames []string, client *statsd.Statter) *TimerWithCounter {
	timer := NewTimer(namespace, subsystem, name+"Timer", labelNames, client)
	counter := NewCounter(namespace, subsystem, name+"Counter", labelNames, client)
	return &TimerWithCounter{
		timer,
		counter,
	}
}

func (w *TimerWithCounter) Start() *RunningTimerWithCounter {
	return &RunningTimerWithCounter{
		runningTimer: w.Timer.Start(),
		timerWithCounter: w,
	}
}

func (rt *RunningTimerWithCounter) Done(labelValues ...string) {
	rt.runningTimer.Done(labelValues...)
	rt.timerWithCounter.Counter.Inc(labelValues...)
}
