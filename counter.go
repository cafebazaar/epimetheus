package epimetheus

import (
	"strings"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/prometheus/client_golang/prometheus"
)

type Counter struct {
	watcher *prometheus.CounterVec
	client  *statsd.Statter
	prefix  string
	labels  []string
}

type StaticCounter struct {
	Base   *Counter
	values []string
}

func NewCounter(namespace, subsystem, name string, labelNames []string, client *statsd.Statter) *Counter {
	opts := prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
	}
	vec := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(vec)

	return &Counter{
		watcher: vec,
		labels:  labelNames,
		client:  client,
		prefix:  strings.Join([]string{namespace, subsystem, name}, "."),
	}
}

func (w *Counter) Inc(labelValues ...string) {
	w.watcher.WithLabelValues(labelValues...).Inc()
	metaLabel := w.prefix + "." + strings.Join(labelValues, ".")
	(*w.client).Inc(metaLabel, 1, 1.0)
}

func (w *Counter) NewStaticCounter(labelValues ...string) *StaticCounter {
	return &StaticCounter{
		Base:   w,
		values: labelValues,
	}
}

func (sc *StaticCounter) Inc() {
	sc.Base.Inc(sc.values...)
}
