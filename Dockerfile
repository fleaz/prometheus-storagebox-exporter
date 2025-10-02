FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN apk --no-cache add ca-certificates
WORKDIR /go/src/prometheus-storagebox-exporter
COPY . /go/src/prometheus-storagebox-exporter
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags '-extldflags "-static"' -o prometheus-storagebox-exporter

FROM scratch
LABEL org.opencontainers.image.title=prometheus-storagebox-exporter
LABEL org.opencontainers.image.description="Prometheus Exporter for Hetzner's Storagebox Service"
LABEL org.opencontainers.image.authors="Felix Breidenstein <mail@felixbreidenstein.de>"
LABEL org.opencontainers.image.url=https://github.com/fleaz/prometheus-storagebox-exporter
LABEL org.opencontainers.image.licenses=MIT

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/prometheus-storagebox-exporter/prometheus-storagebox-exporter /
CMD ["./prometheus-storagebox-exporter"]
