package main

import (
	"log/slog"
	"net/http"
	"text/template"
)

const (
	robotsTemplate string = `
User-agent: *
Disallow: /
`
	rootTemplate string = `
{{- define "content" }}
<!DOCTYPE html>
<html lang="en-US">
<head>
	<meta name="description" content="Prometheus Exporter for Vultr Status service">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Prometheus Exporter for Vultr Status service</title>
	<style>
	body { font-family: Verdana; }
	</style>
</head>
<body>
	<h2>Prometheus Exporter for Vultr Status service</h2>
	<hr/>
	<ul>
	<li><a href="{{ .MetricsPath }}">metrics</a></li>
	<li><a href="/healthz">healthz</a></li>
	</ul>
</body>
</html>
{{- end}}
`
)

type Content struct {
	MetricsPath string
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		slog.Error("unable to write",
			"err", err,
		)
	}
}

func robots(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(robotsTemplate)); err != nil {
		slog.Error("unable to write",
			"err", err,
		)
	}
}

func root(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	t := template.Must(template.New("content").Parse(rootTemplate))
	if err := t.ExecuteTemplate(w, "content", Content{MetricsPath: *metricsPath}); err != nil {
		slog.Error("unable to execute template",
			"err", err,
		)
	}
}
