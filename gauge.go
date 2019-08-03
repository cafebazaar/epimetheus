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

func NewGauge(opts prometheus.GaugeOpts, labels []string) *Gauge {
	vec := prometheus.NewGaugeVec(opts, labels)
	prometheus.MustRegister(vec)

	return &Gauge{
		watcher: vec,
		labels:  labels,
	}
}

func (w *Gauge) Set(labels map[string]string, value float64) {
	values := make([]string, len(w.labels))
	for i, label := range w.labels {
		values[i] = labels[label]
	}
	w.watcher.WithLabelValues(values...).Set(value)
}

func (w *Gauge) NewStaticGauge(labels map[string]string) *StaticGauge {
	values := make([]string, len(w.labels))
	for i, label := range w.labels {
		values[i] = labels[label]
	}
	return &StaticGauge{
		Base:   w,
		values: values,
	}
}

func (rg *StaticGauge) Set(value float64) {
	rg.Base.watcher.WithLabelValues(rg.values...).Set(value)
}
