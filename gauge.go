package epimetheus

import (
	"strings"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/prometheus/client_golang/prometheus"
)

// Gauge keeps the contents of underlying gauge, including labels
type Gauge struct {
	watcher *prometheus.GaugeVec
	client  *statsd.Statter
	prefix  string
	labels  []string
}

// StaticGauge keeps the contents of underlying gauge, excluding labels
type StaticGauge struct {
	Base   *Gauge
	values []string
}

// newGauge creates a prometheus.GaugeVec and register it only if isPrometheusEnabled is true otherwise it keeps
// watcher unregistered to avoid multiple register error in development setups.
func newGauge(namespace, subsystem, name string, labelNames []string, client *statsd.Statter, isPrometheusEnabled bool) *Gauge {
	opts := prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
	}
	vec := prometheus.NewGaugeVec(opts, labelNames)
	if isPrometheusEnabled {
		prometheus.MustRegister(vec)
	}

	return &Gauge{
		watcher: vec,
		labels:  labelNames,
		client:  client,
		prefix:  strings.Join([]string{namespace, subsystem, name}, "."),
	}
}

// Set sets value of the gauge
func (w *Gauge) Set(value float64, labelValues ...string) {
	w.watcher.WithLabelValues(labelValues...).Set(value)
	metaLabel := w.prefix + "." + strings.Join(labelValues, ".")
	(*w.client).Gauge(metaLabel, int64(value), 1.0)
}

func (w *Gauge) newStaticGauge(labelValues ...string) *StaticGauge {
	return &StaticGauge{
		Base:   w,
		values: labelValues,
	}
}

// Set sets value of the gauge
func (rg *StaticGauge) Set(value float64) {
	rg.Base.watcher.WithLabelValues(rg.values...).Set(value)
}
