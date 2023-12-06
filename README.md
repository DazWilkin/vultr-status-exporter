# Prometheus Exporter for Vultr Status

[![build](https://github.com/DazWilkin/vultr-status-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/DazWilkin/vultr-status-exporter/actions/workflows/build.yml)

A Prometheus Exporter for [Vultr Server Status JSON Endpoints](https://www.vultr.com/docs/vultr-server-status-json-endpoints)

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
ghcr.io/dazwilkin/vultr-status-exporter:1234567890123456789012345678901234567890
```

### Kubernetes

```bash
NAMESPACE="vultr-status-exporter"

kubectl create namespace ${NAMESPACE}

kubectl apply \
--filename=./kubernetes.yaml \
--namespace=${NAMESPACE}
```

## Sigstore

`vultr-status-exporter`` container images are being signed by Sigstore and may be verified:

```bash
cosign verify \
--key=./cosign.pub \
ghcr.io/dazwilkin/vultr-status-exporter:1234567890123456789012345678901234567890
```

> **NOTE** `cosign.pub` may be downloaded [here]()

To install `cosign`, e.g.:

```bash
go install github.com/sigstore/cosign/cmd/cosign@latest
```