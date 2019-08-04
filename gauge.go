package epimetheus

import (
	"strings"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/prometheus/client_golang/prometheus"
)

type Gauge struct {
	watcher *prometheus.GaugeVec
	client  *statsd.Statter
	prefix  string
	labels  []string
}

type StaticGauge struct {
	Base   *Gauge
	values []string
}

func NewGauge(namespace, subsystem, name string, labelNames []string, client *statsd.Statter) *Gauge {
	opts := prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
	}
	vec := prometheus.NewGaugeVec(opts, labelNames)
	prometheus.MustRegister(vec)

	return &Gauge{
		watcher: vec,
		labels:  labelNames,
		client:  client,
		prefix:  strings.Join([]string{namespace, subsystem, name}, "."),
	}
}

func (w *Gauge) Set(value float64, labelValues ...string) {
	w.watcher.WithLabelValues(labelValues...).Set(value)
	metaLabel := w.prefix + "." + strings.Join(labelValues, ".")
	(*w.client).Gauge(metaLabel, int64(value), 1.0)
}

func (w *Gauge) NewStaticGauge(labelValues ...string) *StaticGauge {
	return &StaticGauge{
		Base:   w,
		values: labelValues,
	}
}

func (rg *StaticGauge) Set(value float64) {
	rg.Base.watcher.WithLabelValues(rg.values...).Set(value)
}
