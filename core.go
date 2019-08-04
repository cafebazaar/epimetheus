package epimetheus

import (
	"fmt"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Epimetheus struct {
	config        *viper.Viper
	CommTimer     *Timer
	FunctionTimer *Timer
	CacheRate     *Counter
}

func NewEpimetheus(config *viper.Viper) *Epimetheus {
	return &Epimetheus{
		config: config,
	}
}

func (e *Epimetheus) InitWatchers() {
	if e.CommTimer != nil {
		return
	}
	ctLabels := [...]string{"service", "method", "status"}
	e.CommTimer = e.NewTimer("Communications", ctLabels[:])
	ptLabels := [...]string{"funcName"}
	e.FunctionTimer = e.NewTimer("Functions", ptLabels[:])
	crLabels := [...]string{"cacheName", "status"}
	e.CacheRate = e.NewCounter("Caches", crLabels[:])
}

func (e *Epimetheus) Listen() {
	port := e.config.GetInt("prometheus.port")
	server := NewServer(port)
	logrus.Debugf("Epimetheus is Listening on port %d", port)
	server.Serve()
}

func (e *Epimetheus) MakeClient() *statsd.Statter {
	port := e.config.GetInt("statsd.port")
	host := e.config.GetString("statsd.host")
	addr := fmt.Sprintf("%s:%d", host, port)
	logrus.Debugf("Statsd is sending to %s", addr)
	client, err := statsd.NewBufferedClient(addr, "", 500*time.Millisecond, 0)
	if err != nil {
		logrus.Error("Failed to start Statsd Client")
	}
	return &client
}

func (e *Epimetheus) NewTimer(name string, labelNames []string) *Timer {
	namespace := e.config.GetString("prometheus.namespace")
	subsystem := e.config.GetString("prometheus.system-name")
	client := e.MakeClient()
	return NewTimer(namespace, subsystem, name, labelNames, client)
}

func (e *Epimetheus) NewCounter(name string, labelNames []string) *Counter {
	namespace := e.config.GetString("prometheus.namespace")
	subsystem := e.config.GetString("prometheus.system-name")
	client := e.MakeClient()
	return NewCounter(namespace, subsystem, name, labelNames, client)
}

func (e *Epimetheus) NewGauge(name string, labelNames []string) *Gauge {
	namespace := e.config.GetString("prometheus.namespace")
	subsystem := e.config.GetString("prometheus.system-name")
	client := e.MakeClient()
	return NewGauge(namespace, subsystem, name, labelNames, client)
}
