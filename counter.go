package epimetheus

import "github.com/prometheus/client_golang/prometheus"

type Counter struct {
	watcher *prometheus.CounterVec
	labels  []string
}

type StaticCounter struct {
	Base   *Counter
	values []string
}

func NewCounter(opts prometheus.CounterOpts, labelNames []string) *Counter {
	vec := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(vec)

	return &Counter{
		watcher: vec,
		labels:  labelNames,
	}
}

func (w *Counter) Inc(labelValues ...string) {
	w.watcher.WithLabelValues(labelValues...).Inc()
}

func (w *Counter) NewStaticCounter(labelValues ...string) *StaticCounter {
	return &StaticCounter{
		Base:   w,
		values: labelValues,
	}
}

func (sc *StaticCounter) Inc() {
	sc.Base.watcher.WithLabelValues(sc.values...).Inc()
}
