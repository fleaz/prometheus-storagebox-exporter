package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fleaz/prometheus-storagebox-exporter/hetzner"

	"github.com/imroc/req/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	hetznerToken string
	listenAddr   string
	labels       = []string{"id", "name", "product", "server"}
	diskQuota    = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "storagebox",
			Name:      "disk_quota",
			Help:      "Total diskspace in Bytes",
		},
		labels,
	)
	diskUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "storagebox",
			Name:      "disk_usage",
			Help:      "Total used diskspace in Bytes",
		},
		labels,
	)
	diskUsageData = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "storagebox",
			Name:      "disk_usage_data",
			Help:      "Used diskspace by files in Bytes",
		},
		labels,
	)
	diskUsageSnapshots = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "storagebox",
			Name:      "disk_usage_snapshots",
			Help:      "Used diskspace by snapshots in Bytes",
		},
		labels,
	)
)

func updateMetrics() {
	for {
		boxes, err := hetzner.GetBoxes()
		if err != nil {
			log.Println("Failed to get storageboxes!")
		}
		for _, box := range boxes {
			diskQuota.With(prometheus.Labels{
				"id":      strconv.Itoa(box.ID),
				"name":    box.Name,
				"product": box.StorageBoxType.Description,
				"server":  box.Server,
			}).Set(float64(box.StorageBoxType.Size))

			diskUsage.With(prometheus.Labels{
				"id":      strconv.Itoa(box.ID),
				"name":    box.Name,
				"product": box.StorageBoxType.Description,
				"server":  box.Server,
			}).Set(float64(box.Stats.Size))

			diskUsageData.With(prometheus.Labels{
				"id":      strconv.Itoa(box.ID),
				"name":    box.Name,
				"product": box.StorageBoxType.Description,
				"server":  box.Server,
			}).Set(float64(box.Stats.SizeData))

			diskUsageSnapshots.With(prometheus.Labels{
				"id":      strconv.Itoa(box.ID),
				"name":    box.Name,
				"product": box.StorageBoxType.Description,
				"server":  box.Server,
			}).Set(float64(box.Stats.SizeSnapshots))

		}

		// Try to avoid rate limiting
		// Limit is 200req / 1h
		time.Sleep(60 * time.Second)
	}
}

func main() {
	hetznerToken = os.Getenv("HETZNER_TOKEN")

	if hetznerToken == "" {
		log.Fatal("Please provide HETZNER_TOKEN as environment variables")
	}

	listenAddr = os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":9509"
	}

	hetzner.Client = req.C()
	hetzner.Client.SetCommonBearerAuthToken(hetznerToken)

	go updateMetrics()

	log.Printf("Listening on %q", listenAddr)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(listenAddr, nil)
}
