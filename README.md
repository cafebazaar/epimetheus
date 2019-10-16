# Epimetheus

Epimetheus is a lightweight wrpper around [Epimetheus Go client](https://github.com/prometheus/client_golang) and [Statsd Go client](github.com/cactus/go-statsd-client) which makes measuring communication, functions, background jobs and ... easier simultaneously with both Prometheus and Statsd.

## Getting Started

### Installing

```console
go get -u git.cafebazaar.ir/bazaar/search/epimetheus
```

### Initialize server

```go
epimetheusServer := epimetheus.NewEpimetheusServer(config)
go epimetheusServer.Listen()

```

### Predefined metrics

```go
epi := epimetheus.NewEpimetheus(config)
commTimer := epi.CommTimer
go epimetheusServer.Listen()
```

### Measuring duration and count
```go
epi := epimetheus.NewEpimetheus(config)
timerWithCounter := epi.NewTimerWithCounter("req1", string[]{"label1"})
rtc := timerWithCounter.Start()
// Do some work here
rtc.Done("dispatch")
```


### Measuring duration
```go
epi := epimetheus.NewEpimetheus(config)
timer := epi.NewTimer("req1", string[]{"label1"})
t := timer.Start()
// Do some work here
t.Done("dispatch")
```

### Measuring count
```go
epi := epimetheus.NewEpimetheus(config)
counter := epi.NewCounter("req1", string[]{"label1"})
// Do some work here
c.Inc("dispatch")
```

## Configuration

Epimetheus uses Viper as it's config engine. Template should be something like this:
```yaml
    stats:
      prometheus:
        enabled: true
        port: 1234
      statsd:
        enabled: true
        port: 5678
        host: "w.x.y.z"
      namespace: search
      system-name: octopus
```

## Documentation

Documents are available at [https://godoc.org/git.cafebazaar.ir/bazaar/search/epimetheus](https://godoc.org/git.cafebazaar.ir/bazaar/search/epimetheus)

## Built With

* [Prometheus Go clinet](github.com/prometheus/client_golang) - The underlying library for Prometheus
* [Statsd Go client](github.com/cactus/go-statsd-client) - The underlying library for Statsd

## Contributing

Please read [CONTRIBUTING.md](https://git.cafebazaar.ir/bazaar/search/epimetheus/blob/master/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://git.cafebazaar.ir/bazaar/search/epimetheus/tags). 

## Roadmap
    - Improve documentation
    - Add tests

## Authors

* **Ramtin Rostami** - *Initial work* - [rrostami](https://github.com/rrostami)
* **Pedram Teymoori** - *Initial work* - [pedramteymoori](https://github.com/pedramteymoori)
* **Parsa abdollahi** - *Initial work* - []()

See also the list of [contributors](https://git.cafebazaar.ir/bazaar/search/epimetheus/-/graphs/master) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* [Sepehr nematollahi](https://www.behance.net/sseeppeehhrr) Epimetheus typography designer

Made with <span class="heart">❤</span> in Cafebazaar search