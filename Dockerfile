ARG GOLANG_VERSION=1.21.5

ARG GOOS=linux
ARG GOARCH=amd64

ARG COMMIT
ARG VERSION

FROM docker.io/golang:${GOLANG_VERSION} as build

WORKDIR /vultr-status-exporter

COPY api api
COPY cmd cmd
COPY collector ./collector

ARG GOOS
ARG GOARCH

ARG VERSION
ARG COMMIT

RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /go/bin/vultr-status-exporter \
    ./cmd/exporter/main.go

FROM gcr.io/distroless/static-debian11:latest

LABEL org.opencontainers.image.description "Prometheus Exporter for Vultr Status"
LABEL org.opencontainers.image.source https://github.com/DazWilkin/vultr-status-exporter

COPY --from=build /go/bin/vultr-status-exporter /

EXPOSE 9402

ENTRYPOINT ["/vultr-status-exporter"]
CMD ["--endpoint=:8080","--path=/metrics"]