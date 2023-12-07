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
ghcr.io/dazwilkin/vultr-status-exporter:6801ae47f0b019bc0c2fa260328d64789803610b
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
ghcr.io/dazwilkin/vultr-status-exporter:6801ae47f0b019bc0c2fa260328d64789803610b
```

> **NOTE** `cosign.pub` may be downloaded [here]()

To install `cosign`, e.g.:

```bash
go install github.com/sigstore/cosign/cmd/cosign@latest
```