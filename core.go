package epimetheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

type Epimetheus struct {
	config *viper.Viper
	timers map[string]*Timer
}

func NewEpimetheusWatcher(config *viper.Viper) *Epimetheus {
	return &Epimetheus{
		config: config,
	}
}

func (e *Epimetheus) Listen() {
	port := e.config.GetInt("prometheus.port")
	server := NewServer(port)
	server.Serve()
}

func (e *Epimetheus) NewCommunicationTimer() *Timer {
	if e.timers["Communications"] != nil {
		return e.timers["Communications"]
	}
	// making the labels default to increase simplicity
	labels := [...]string{"service", "method", "status"}
	ct := e.NewTimer("Communications", labels[:])
	e.timers["Communications"] = ct
	return ct
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

func CommParams(service string, method string, status string) []string {
	return map[string]string{
		"service": service,
		"method":  method,
		"status":  status,
	}
}
