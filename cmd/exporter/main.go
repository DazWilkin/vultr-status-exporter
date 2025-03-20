package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/DazWilkin/go-probe/probe"
	"github.com/DazWilkin/vultr-status-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// GitCommit is the git commit value and is expected to be set during build
	GitCommit string
	// GoVersion is the Golang runtime version
	GoVersion = runtime.Version()
	// OSVersion is the OS version (uname --kernel-release) and is expected to be set during build
	OSVersion string
	// StartTime is the start time of the exporter represented as a UNIX epoch
	StartTime = time.Now().Unix()
)
var (
	endpoint    = flag.String("endpoint", ":8080", "The endpoint of the HTTP server")
	metricsPath = flag.String("path", "/metrics", "The path on which Prometheus metrics will be served")
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	flag.Parse()

	if GitCommit == "" {
		logger.Info("expected value to be set during build",
			"GitCommit", GitCommit,
		)
	}
	if OSVersion == "" {
		logger.Info("expected value to be set during build",
			"OSVersion", OSVersion,
		)
	}

	p := probe.New("liveness", logger)
	healthz := p.Handler(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan probe.Status)
	go p.Updater(ctx, ch, nil)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewExporterCollector(OSVersion, GoVersion, GitCommit, StartTime))
	registry.MustRegister(collector.NewStatusCollector(ch, logger))

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(root))
	mux.Handle("/healthz", http.HandlerFunc(healthz))
	mux.Handle("/robots.txt", http.HandlerFunc(robots))

	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	logger.Info("Server starting",
		"endpoint", *endpoint,
		"metricsPath", *metricsPath,
	)
	logger.Error("server error",
		"err", http.ListenAndServe(*endpoint, mux),
	)
}
