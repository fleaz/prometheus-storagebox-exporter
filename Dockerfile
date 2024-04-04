FROM golang:1.22 AS builder
WORKDIR /go/src/prometheus-storagebox-exporter
COPY . /go/src/prometheus-storagebox-exporter
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest
LABEL org.opencontainers.image.title=prometheus-storagebox-exporter
LABEL org.opencontainers.image.description="Prometheus Exporter for Hetzner's Storagebox Service"
LABEL org.opencontainers.image.authors="Felix Breidenstein <mail@felixbreidenstein.de>"
LABEL org.opencontainers.image.url=https://github.com/fleaz/prometheus-storagebox-exporter
LABEL org.opencontainers.image.licenses=MIT

RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/prometheus-storagebox-exporter/prometheus-storagebox-exporter /
CMD ["./prometheus-storagebox-exporter"]
