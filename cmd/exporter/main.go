package main

import (
	"flag"
	"log/slog"
	"net/http"
	"runtime"
	"time"

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
	flag.Parse()

	if GitCommit == "" {
		slog.Info("expected value to be set during build",
			"GitCommit", GitCommit,
		)
	}
	if OSVersion == "" {
		slog.Info("expected value to be set during build",
			"OSVersion", OSVersion,
		)
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewExporterCollector(OSVersion, GoVersion, GitCommit, StartTime))
	registry.MustRegister(collector.NewStatusCollector())

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(root))
	mux.Handle("/healthz", http.HandlerFunc(healthz))
	mux.Handle("/robots.txt", http.HandlerFunc(robots))

	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	slog.Info("Server starting",
		"endpoint", *endpoint,
		"metricsPath", *metricsPath,
	)
	slog.Error("server error",
		"err", http.ListenAndServe(*endpoint, mux),
	)
}
