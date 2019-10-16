package epimetheus

import (
	"fmt"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Epimetheus is the main struct of this library
//
// You should have exactly one instance of it with `NewEpimetheusServer`, and then call `listen()` method of it
//
//  - config: Instance of viper
//  - CommTimer: An instance of `TimerWithCounter` which is used for your application communications
//  - FunctionTimer: An instance of `TimerWithCounter` which is used for your application important functions
//  - CacheRate: An instance of `Counter` which is used in conjunction with `Mnemosyne` project for caching
//  - BGWorker: An instance of `Counter` which is used for your application background workers
type Epimetheus struct {
	config        *viper.Viper
	CommTimer     *TimerWithCounter
	FunctionTimer *TimerWithCounter
	CacheRate     *Counter
	BGWorker      *Counter
}

// NewEpimetheus Creates a new instance of `Epimetheus` with prepared metrics in advance
func NewEpimetheus(config *viper.Viper) *Epimetheus {
	e := &Epimetheus{
		config: config,
	}
	ctLabels := [...]string{"service", "method", "status"}
	e.CommTimer = e.NewTimerWithCounter("Communications", ctLabels[:])
	ptLabels := [...]string{"funcName"}
	e.FunctionTimer = e.NewTimerWithCounter("Functions", ptLabels[:])
	crLabels := [...]string{"cacheName", "status"}
	e.CacheRate = e.NewCounter("Caches", crLabels[:])
	bgLabels := [...]string{"type", "status"}
	e.BGWorker = e.NewCounter("BGWorker", bgLabels[:])
	return e
}

// NewEpimetheusServer Creates a new instance of `Epimetheus` without prepared metrics in advance
func NewEpimetheusServer(config *viper.Viper) *Epimetheus {
	return &Epimetheus{
		config: config,
	}
}

// Listen makes epimetheus start on specified port if enabled, then return the server in order to stop later
func (e *Epimetheus) Listen() *Server {
	if !e.config.GetBool("stats.prometheus.enabled") {
		return nil
	}
	port := e.config.GetInt("stats.prometheus.port")
	server := newServer(port)
	logrus.Infof("Epimetheus is Listening on port %d", port)
	server.serve()
	return server
}

func (e *Epimetheus) makeClient() *statsd.Statter {
	if !e.config.GetBool("stats.statsd.enabled") {
		client, _ := statsd.NewNoopClient()
		return &client
	}
	port := e.config.GetInt("stats.statsd.port")
	host := e.config.GetString("stats.statsd.host")
	addr := fmt.Sprintf("%s:%d", host, port)
	logrus.Infof("Statsd is sending to %s", addr)
	client, err := statsd.NewBufferedClient(addr, "", 500*time.Millisecond, 0)
	if err != nil {
		logrus.Error("Failed to start Statsd Client")
	}
	return &client
}

// NewTimer creates an instance of `Timer` with specified configs
func (e *Epimetheus) NewTimer(name string, labelNames []string) *Timer {
	namespace := e.config.GetString("stats.namespace")
	subsystem := e.config.GetString("stats.system-name")
	client := e.makeClient()
	return newTimer(namespace, subsystem, name, labelNames, client)
}

// NewCounter creates an instance of `Counter` with specified configs
func (e *Epimetheus) NewCounter(name string, labelNames []string) *Counter {
	namespace := e.config.GetString("stats.namespace")
	subsystem := e.config.GetString("stats.system-name")
	client := e.makeClient()
	return newCounter(namespace, subsystem, name, labelNames, client)
}

// NewGauge creates an instance of `Gauge` with specified configs
func (e *Epimetheus) NewGauge(name string, labelNames []string) *Gauge {
	namespace := e.config.GetString("stats.namespace")
	subsystem := e.config.GetString("stats.system-name")
	client := e.makeClient()
	return newGauge(namespace, subsystem, name, labelNames, client)
}

// NewTimerWithCounter creates an instance of `TimerWithCounter` with specified configs
func (e *Epimetheus) NewTimerWithCounter(name string, labelNames []string) *TimerWithCounter {
	namespace := e.config.GetString("stats.namespace")
	subsystem := e.config.GetString("stats.system-name")
	client := e.makeClient()
	return newTimerWithCounter(namespace, subsystem, name, labelNames, client)
}
