package epimetheus

import (
	"strings"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/prometheus/client_golang/prometheus"
)

// Counter keeps the contents of underlying counter, including labels
type Counter struct {
	watcher *prometheus.CounterVec
	client  *statsd.Statter
	prefix  string
	labels  []string
}

// StaticCounter keeps the contents of underlying counter, excluding labels
type StaticCounter struct {
	Base   *Counter
	values []string
}

// newCounter creates a prometheus.CounterVec and register it only if isPrometheusEnabled is true otherwise it keeps
// watcher unregistered to avoid multiple register error in development setups.
func newCounter(namespace, subsystem, name string, labelNames []string, client *statsd.Statter, isPrometheusEnabled bool) *Counter {
	opts := prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
	}
	vec := prometheus.NewCounterVec(opts, labelNames)
	if isPrometheusEnabled {
		prometheus.MustRegister(vec)
	}

	return &Counter{
		watcher: vec,
		labels:  labelNames,
		client:  client,
		prefix:  strings.Join([]string{namespace, subsystem, name}, "."),
	}
}

// Inc increments the value of current Counter
func (w *Counter) Inc(labelValues ...string) {
	w.watcher.WithLabelValues(labelValues...).Inc()
	metaLabel := w.prefix + "." + strings.Join(labelValues, ".")
	(*w.client).Inc(metaLabel, 1, 1.0)
}

func (w *Counter) newStaticCounter(labelValues ...string) *StaticCounter {
	return &StaticCounter{
		Base:   w,
		values: labelValues,
	}
}

// Inc increments the value of current Counter
func (sc *StaticCounter) Inc() {
	sc.Base.Inc(sc.values...)
}
