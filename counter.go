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

func NewCounter(opts prometheus.CounterOpts, labels []string) *Counter {
	vec := prometheus.NewCounterVec(opts, labels)
	prometheus.MustRegister(vec)

	return &Counter{
		watcher: vec,
		labels:  labels,
	}
}

func (w *Counter) Inc(labels map[string]string) {
	values := make([]string, len(w.labels))
	for i, label := range w.labels {
		values[i] = labels[label]
	}
	w.watcher.WithLabelValues(values...).Inc()
}

func (w *Counter) NewStaticCounter(labels map[string]string) *StaticCounter {
	values := make([]string, len(w.labels))
	for i, label := range w.labels {
		values[i] = labels[label]
	}
	return &StaticCounter{
		Base:   w,
		values: values,
	}
}

func (sc *StaticCounter) Inc() {
	sc.Base.watcher.WithLabelValues(sc.values...).Inc()
}
