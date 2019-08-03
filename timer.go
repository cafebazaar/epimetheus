package epimetheus

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Timer struct {
	watcher *prometheus.HistogramVec
	labels  []string
}

type RunningTimer struct {
	Base  *Timer
	start time.Time
	end   time.Time
}

func NewTimer(opts prometheus.HistogramOpts, labelNames []string) *Timer {
	vec := prometheus.NewHistogramVec(opts, labelNames)
	prometheus.MustRegister(vec)

	return &Timer{
		watcher: vec,
		labels:  labelNames,
	}
}

func (w *Timer) RunWithError(work func() error, labelValues ...string) error {
	start := time.Now()
	err := work()
	end := time.Now()
	if w.hasLabel("ok") {
		labelValues["ok"] = fmt.Sprint(err == nil)
	}
	w.register(start, end, labelValues...)
	return err
}

func (w *Timer) RunVoid(work func(), labelValues ...string) {
	start := time.Now()
	work()
	end := time.Now()
	w.register(start, end, labelValues)
}

func (w *Timer) Start() *RunningTimer {
	return &RunningTimer{
		Base:  w,
		start: time.Now(),
	}
}

func (rt *RunningTimer) Done(labelValues ...string) {
	rt.end = time.Now()
	rt.Base.register(rt.start, rt.end, labelValues)
}

func (w *Timer) register(start time.Time, end time.Time, labelValues []string) {
	duration := end.Sub(start)
	w.watcher.WithLabelValues(labelValues...).Observe(duration.Seconds())
}

func (w *Timer) hasLabel(label string) bool {
	for _, l := range w.labels {
		if l == label {
			return true
		}
	}
	return false
}
