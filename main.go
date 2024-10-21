package main

import (
	"crash_exporter/websocket"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
)

type Flags struct {
	Listen     string
	Host       string
	Token      string
	LogLevel   string
	LokiServer string
}

func parseFlags() Flags {
	var flags Flags
	flag.StringVar(&flags.Listen, "listen", "0.0.0.0:9091", "Listen address")
	flag.StringVar(&flags.Host, "host", "192.168.66.1:9999", "ShellCrash address (Required)")
	flag.StringVar(&flags.Token, "token", "", "Crash token")
	flag.StringVar(&flags.LogLevel, "log-level", "debug", "LogLevel: debug|info|warning|error. Default: info.")
	flag.StringVar(&flags.LokiServer, "loki-server", "http://192.168.66.2:3100", "Loki server address")
	flag.Parse()
	return flags
}

func dynamicLabelProvider(entry *logrus.Entry) lokirus.Labels {
	labels := make(lokirus.Labels, len(entry.Data))
	for k, v := range entry.Data {
		labels[k] = fmt.Sprintf("%+v", v)
	}
	return labels
}

func getLoggerWithLoki(lokiAddr string) *logrus.Logger {
	// Configure the Loki hook
	opts := lokirus.NewLokiHookOptions().
		// Grafana doesn't have a "panic" level, but it does have a "critical" level
		// https://grafana.com/docs/grafana/latest/explore/logs-integration/
		WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
		WithStaticLabels(lokirus.Labels{
			"service_name": "crash-exporter",
		}).
		WithDynamicLabelProvider(dynamicLabelProvider).
		WithFormatter(&logrus.JSONFormatter{})

	hook := lokirus.NewLokiHookWithOpts(
		lokiAddr,
		opts,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel)

	// Configure the logger
	logger := logrus.New()
	logger.AddHook(hook)
	return logger
}

func main() {
	args := parseFlags()
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	logger := getLoggerWithLoki(args.LokiServer)

	crash, err := websocket.NewCrash(args.Host, args.Token, reg, logger)
	if err != nil {
		log.Fatalln(err)
	}
	if err = crash.Registers(args.LogLevel); err != nil {
		log.Fatalln(err)
	}

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		Registry: reg, EnableOpenMetrics: true}))

	logger.Infof("Exporter started at %s/metrics", args.Listen)
	server := &http.Server{
		Addr:              args.Listen,
		ReadHeaderTimeout: 1 * time.Second,
	}
	if err = server.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
