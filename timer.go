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

func NewTimer(opts prometheus.HistogramOpts, labels []string) *Timer {
	vec := prometheus.NewHistogramVec(opts, labels)
	prometheus.MustRegister(vec)

	return &Timer{
		watcher: vec,
		labels:  labels,
	}
}

func (w *Timer) RunWithError(work func() error, labels map[string]string) error {
	start := time.Now()
	err := work()
	end := time.Now()
	if w.hasLabel("ok") {
		labels["ok"] = fmt.Sprint(err == nil)
	}
	w.register(start, end, labels)
	return err
}

func (w *Timer) RunVoid(work func(), labels map[string]string) {
	start := time.Now()
	work()
	end := time.Now()
	w.register(start, end, labels)
}

func (w *Timer) Start() *RunningTimer {
	return &RunningTimer{
		Base:  w,
		start: time.Now(),
	}
}

func (rt *RunningTimer) Done(labels map[string]string) {
	rt.end = time.Now()
	rt.Base.register(rt.start, rt.end, labels)
}

func (w *Timer) register(start time.Time, end time.Time, labels map[string]string) {
	values := make([]string, len(w.labels))
	for i, label := range w.labels {
		values[i] = labels[label]
	}
	duration := end.Sub(start)
	w.watcher.WithLabelValues(values...).Observe(duration.Seconds())
}

func (w *Timer) hasLabel(label string) bool {
	for _, l := range w.labels {
		if l == label {
			return true
		}
	}
	return false
}
