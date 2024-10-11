# Prometheus Exporter for Vultr Status

[![build](https://github.com/DazWilkin/vultr-status-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/DazWilkin/vultr-status-exporter/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/DazWilkin/vultr-status-exporter.svg)](https://pkg.go.dev/github.com/DazWilkin/vultr-status-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/DazWilkin/vultr-status-exporter)](https://goreportcard.com/report/github.com/DazWilkin/vultr-status-exporter)

A Prometheus Exporter for [Vultr Server Status JSON Endpoints](https://www.vultr.com/docs/vultr-server-status-json-endpoints)

## Metrics

Metrics are prefixed `vultr_status_`

|Name|Type|Description|
|----|----|-----------|
|`vultr_exporter_build_info`|Counter|A metric with a constant '1' value labeled by OS version, Go version, and the Git commit of the exporter|
|`vultr_exporter_start_time`|Gauge|Exporter start time in Unix epoch seconds|
|`vultr_status_infrastructure`|Gauge|Vultr Infrastructure status|
|`vultr_status_service_alert`|Counter|Vultr Service Alerts|

## Installation

### Go

```bash
go install github.com/DazWilkin/vultr-status-exporter@latest
```

### Container

```bash
PORT="8080"

podman run \
--interactive --tty \
--publish=${PORT}:${PORT} \
ghcr.io/dazwilkin/vultr-status-exporter:e02426de88ba1ce737b4ce12ad26fa46bbe8efa8
```

### Kubernetes

```bash
NAMESPACE="exporter"

kubectl create namespace ${NAMESPACE}

# Apply Deployment|Service
kubectl apply \
--filename=./kubernetes.yaml \
--namespace=${NAMESPACE}

# Get Service NodePort
NODE_PORT=$(\
  kubectl get service/vultr-status-exporter \
  --namespace=${NAMESPACE} \
  --output=jsonpath="{.spec.ports[?(@.name==\"metrics\")].nodePort}") && echo ${NODE_PORT}
```

Browse: `localhost:${NODE_PORT}`

## Sigstore

`vultr-status-exporter`` container images are being signed by Sigstore and may be verified:

```bash
cosign verify \
--key=./cosign.pub \
ghcr.io/dazwilkin/vultr-status-exporter:e02426de88ba1ce737b4ce12ad26fa46bbe8efa8
```

> **NOTE** `cosign.pub` may be downloaded [here](./cosign.pub)

To install `cosign`, e.g.:

```bash
go install github.com/sigstore/cosign/cmd/cosign@latest
```

## Similar Exporters

+ [Prometheus Exporter for Azure](https://github.com/DazWilkin/azure-exporter)
+ [Prometheus Exporter for Fly.io](https://github.com/DazWilkin/fly-exporter)
+ [Prometheus Exporter for GCP](https://github.com/DazWilkin/gcp-exporter)
+ [Prometheus Exporter for Koyeb](https://github.com/DazWilkin/koyeb-exporter)
+ [Prometheus Exporter for Linode](https://github.com/DazWilkin/linode-exporter)
+ [Prometheus Exporter for Vultr](https://github.com/DazWilkin/vultr-exporter)

<hr/>
<br/>
<a href="https://www.buymeacoffee.com/dazwilkin" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>