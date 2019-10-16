package epimetheus

import (
	"fmt"
	"strings"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/prometheus/client_golang/prometheus"
)

// Timer keeps the contents of underlying timer, including labels
type Timer struct {
	watcher *prometheus.HistogramVec
	client  *statsd.Statter
	prefix  string
	labels  []string
}

// RunningTimer keeps track of timer times
type RunningTimer struct {
	Base  *Timer
	start time.Time
	end   time.Time
}

func newTimer(namespace, subsystem, name string, labelNames []string, client *statsd.Statter) *Timer {
	opts := prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
	}
	vec := prometheus.NewHistogramVec(opts, labelNames)
	prometheus.MustRegister(vec)
	return &Timer{
		watcher: vec,
		labels:  labelNames,
		client:  client,
		prefix:  strings.Join([]string{namespace, subsystem, name}, "."),
	}
}

// RunWithError gets a function with error and lables and run it with measurment
func (w *Timer) RunWithError(work func() error, labelValues ...string) error {
	start := time.Now()
	err := work()
	end := time.Now()
	if indx := w.hasLabel("ok"); indx > 0 {
		labelValues[indx] = fmt.Sprint(err == nil)
	}
	w.register(start, end, labelValues)
	return err
}

// RunVoid gets a void function and lables and run it with measurment
func (w *Timer) RunVoid(work func(), labelValues ...string) {
	start := time.Now()
	work()
	end := time.Now()
	w.register(start, end, labelValues)
}

// Start begins a RunningTimer and then returns it
func (w *Timer) Start() *RunningTimer {
	return &RunningTimer{
		Base:  w,
		start: time.Now(),
	}
}

// Done finalize the current RunningTimer
func (rt *RunningTimer) Done(labelValues ...string) {
	rt.end = time.Now()
	rt.Base.register(rt.start, rt.end, labelValues)
}

func (w *Timer) register(start time.Time, end time.Time, labelValues []string) {
	duration := end.Sub(start)
	w.watcher.WithLabelValues(labelValues...).Observe(duration.Seconds())
	metaLabel := w.prefix + "." + strings.Join(labelValues, ".")
	(*w.client).Timing(metaLabel, int64(duration/time.Millisecond), 1.0)
	(*w.client).Inc(metaLabel, 1, 1.0)
}

func (w *Timer) hasLabel(label string) int {
	for i, l := range w.labels {
		if l == label {
			return i
		}
	}
	return -1
}
