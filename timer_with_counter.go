package epimetheus

import "github.com/cactus/go-statsd-client/statsd"

type TimerWithCounter struct{
	*Timer
	*Counter
}

func NewTimerWithCounter(namespace, subsystem, name string, labelNames []string, client *statsd.Statter) *TimerWithCounter{
	timer := NewTimer(namespace, subsystem, name + "Timer", labelNames, client)
	counter := NewCounter(namespace, subsystem, name + "Counter", labelNames, client)
	return &TimerWithCounter{
		timer,
		counter,
	}
}

func (w *TimerWithCounter) Done(labelValues ...string) {
	w.Done(labelValues...)
	w.Inc(labelValues...)
}
