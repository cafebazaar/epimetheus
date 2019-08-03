package epimetheus

import "github.com/prometheus/client_golang/prometheus"

type Gauge struct {
	watcher *prometheus.GaugeVec
	labels  []string
}

type StaticGauge struct {
	Base   *Gauge
	values []string
}

func NewGauge(opts prometheus.GaugeOpts, labelNames []string) *Gauge {
	vec := prometheus.NewGaugeVec(opts, labelNames)
	prometheus.MustRegister(vec)

	return &Gauge{
		watcher: vec,
		labels:  labelNames,
	}
}

func (w *Gauge) Set(labelValues ...string, value float64) {
	w.watcher.WithLabelValues(labelValues...).Set(value)
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
