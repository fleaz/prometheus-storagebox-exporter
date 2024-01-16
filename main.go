// Copyright 2024 Felix Breidenstein <mail@felixbreidenstein.de>
//
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type APIBoxList []struct {
	Box struct {
		ID int `json:"id"`
	} `json:"storagebox"`
}

type APIBoxDetail struct {
	Box Storagebox `json:"storagebox"`
}

type APIError struct {
	Error struct {
		Status int    `json:"status"`
		Code   string `json:"code"`
	} `json:"error"`
}

type Storagebox struct {
	ID                   int     `json:"id"`
	Login                string  `json:"login"`
	Name                 string  `json:"name"`
	Product              string  `json:"product"`
	Cancelled            bool    `json:"cancelled"`
	Locked               bool    `json:"locked"`
	LinkedServer         int     `json:"linked_server"`
	PaidUntil            string  `json:"paid_until"`
	DiskQuota            float64 `json:"disk_quota"`
	DiskUsage            float64 `json:"disk_usage"`
	DiskUsageData        float64 `json:"disk_usage_data"`
	DiskUsageSnapshots   float64 `json:"disk_usage_snapshots"`
	Webdav               bool    `json:"webdav"`
	Samba                bool    `json:"samba"`
	SSH                  bool    `json:"ssh"`
	BackupService        bool    `json:"backup_service"`
	ExternalReachability bool    `json:"external_reachability"`
	Zfs                  bool    `json:"zfs"`
	Server               string  `json:"server"`
}

var (
	hetznerUsername string
	hetznerPassword string
	boxes           []Storagebox
	labels          = []string{"id", "name", "product", "server"}
	diskQuota       = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "storagebox",
			Name:      "disk_quota",
			Help:      "Total diskspace in MB",
		},
		labels,
	)
	diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "storagebox",
			Name:      "disk_usage",
			Help:      "Total used diskspace in MB",
		},
		labels,
	)
	diskUsageData = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "storagebox",
			Name:      "disk_usage_data",
			Help:      "Used diskspace by files in MB",
		},
		labels,
	)
	diskUsageSnapshots = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "storagebox",
			Name:      "disk_usage_snapshots",
			Help:      "Used diskspace by snapshots in MB",
		},
		labels,
	)
)

func updateBoxes() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://robot-ws.your-server.de/storagebox", nil)
	req.SetBasicAuth(hetznerUsername, hetznerPassword)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	bodyText, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		var apiErr APIError
		err = json.Unmarshal(bodyText, &apiErr)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("API Error: %d - %s", apiErr.Error.Status, apiErr.Error.Code)
		return
	}

	var apiResponse APIBoxList
	err = json.Unmarshal(bodyText, &apiResponse)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range apiResponse {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://robot-ws.your-server.de/storagebox/%d", entry.Box.ID), nil)
		req.SetBasicAuth(hetznerUsername, hetznerPassword)
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		bodyText, err := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			var apiErr APIError
			err = json.Unmarshal(bodyText, &apiErr)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("API Error: %d - %s", apiErr.Error.Status, apiErr.Error.Code)
			return
		}

		var box APIBoxDetail
		err = json.Unmarshal(bodyText, &box)
		if err != nil {
			log.Fatal(err)
		}
		boxes = append(boxes, box.Box)
	}
}

func updateMetrics() {
	for {
		updateBoxes()
		for _, box := range boxes {
			diskQuota.With(prometheus.Labels{
				"id":      strconv.Itoa(box.ID),
				"name":    box.Name,
				"product": box.Product,
				"server":  box.Server,
			}).Set(box.DiskQuota)

			diskUsage.With(prometheus.Labels{
				"id":      strconv.Itoa(box.ID),
				"name":    box.Name,
				"product": box.Product,
				"server":  box.Server,
			}).Set(box.DiskUsage)

			diskUsageData.With(prometheus.Labels{
				"id":      strconv.Itoa(box.ID),
				"name":    box.Name,
				"product": box.Product,
				"server":  box.Server,
			}).Set(box.DiskUsageData)

			diskUsageSnapshots.With(prometheus.Labels{
				"id":      strconv.Itoa(box.ID),
				"name":    box.Name,
				"product": box.Product,
				"server":  box.Server,
			}).Set(box.DiskUsageSnapshots)

		}

		// Try to avoid rate limiting
		// Limit is 200req / 1h
		time.Sleep(60 * time.Second)
	}
}

const (
	listenAddr = ":9509"
)

func main() {
	hetznerUsername = os.Getenv("HETZNER_USER")
	hetznerPassword = os.Getenv("HETZNER_PASS")

	if hetznerUsername == "" || hetznerPassword == "" {
		log.Fatal("Please provide HETZNER_USER and HETZNER_PASS as environment variables")
	}

	prometheus.MustRegister(diskQuota)
	prometheus.MustRegister(diskUsage)
	prometheus.MustRegister(diskUsageData)
	prometheus.MustRegister(diskUsageSnapshots)

	go updateMetrics()

	fmt.Printf("Listening on %q", listenAddr)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(listenAddr, nil)
}
