ARG GOLANG_VERSION=1.24.3

ARG TARGETOS
ARG TARGETARCH

ARG COMMIT
ARG VERSION

FROM --platform=${TARGETARCH} docker.io/golang:${GOLANG_VERSION} AS build

WORKDIR /vultr-status-exporter

COPY go.* ./

COPY api api
COPY cmd cmd
COPY collector ./collector

ARG TARGETOS
ARG TARGETARCH

ARG VERSION
ARG COMMIT

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /go/bin/vultr-status-exporter \
    ./cmd/exporter/...

FROM --platform=${TARGETARCH} gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.description="Prometheus Exporter for Vultr Status"
LABEL org.opencontainers.image.source=https://github.com/DazWilkin/vultr-status-exporter

COPY --from=build /go/bin/vultr-status-exporter /

EXPOSE 9402

ENTRYPOINT ["/vultr-status-exporter"]
CMD ["--endpoint=:8080","--path=/metrics"]
