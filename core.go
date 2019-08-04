package epimetheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

type Epimetheus struct {
	config        *viper.Viper
	CommTimer     *Timer
	FunctionTimer *Timer
	CacheRate     *Counter
}

func NewEpimetheusWatcher(config *viper.Viper) *Epimetheus {
	e := &Epimetheus{
		config: config,
	}
	ctLabels := [...]string{"service", "method", "status"}
	ct := e.NewTimer("Communications", ctLabels[:])
	ptLabels := [...]string{"funcName"}
	pt := e.NewTimer("Functions", ptLabels[:])
	crLabels := [...]string{"cacheName", "status"}
	cr := e.NewCounter("Caches", crLabels[:])
	return &Epimetheus{
		config:        config,
		CommTimer:     ct,
		FunctionTimer: pt,
		CacheRate:     cr,
	}
}

func (e *Epimetheus) Listen() {
	port := e.config.GetInt("prometheus.port")
	server := NewServer(port)
	server.Serve()
}

func (e *Epimetheus) NewTimer(name string, labelNames []string) *Timer {
	opts := prometheus.HistogramOpts{
		Namespace: e.config.GetString("prometheus.namespace"),
		Subsystem: e.config.GetString("prometheus.system-name"),
		Name:      name,
	}
	return NewTimer(opts, labelNames)
}

func (e *Epimetheus) NewCounter(name string, labelNames []string) *Counter {
	opts := prometheus.CounterOpts{
		Namespace: e.config.GetString("prometheus.namespace"),
		Subsystem: e.config.GetString("prometheus.system-name"),
		Name:      name,
	}
	return NewCounter(opts, labelNames)
}

func (e *Epimetheus) NewGauge(name string, labelNames []string) *Gauge {
	opts := prometheus.GaugeOpts{
		Namespace: e.config.GetString("prometheus.namespace"),
		Subsystem: e.config.GetString("prometheus.system-name"),
		Name:      name,
	}
	return NewGauge(opts, labelNames)
}
