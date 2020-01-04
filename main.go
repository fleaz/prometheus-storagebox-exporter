package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	boxes     []Storagebox
	diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "storagebox",
			Name:      "disk_usage",
			Help:      "Used diskspace in MB",
		},
		[]string{
			"id",
			"name",
			"product",
			"server",
		},
	)
)

func updateBoxes() {
	var username string = ""
	var passwd string = ""

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://robot-ws.your-server.de/storagebox", nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	var apiResponse APIBoxList
	err = json.Unmarshal(bodyText, &apiResponse)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range apiResponse {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://robot-ws.your-server.de/storagebox/%d", entry.Box.ID), nil)
		req.SetBasicAuth(username, passwd)
		resp, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		bodyText, err := ioutil.ReadAll(resp.Body)
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
			diskUsage.With(prometheus.Labels{
				"id":      strconv.Itoa(box.ID),
				"name":    box.Name,
				"product": box.Product,
				"server":  box.Server,
			}).Set(box.DiskUsage)

		}
		time.Sleep(30 * time.Second)
	}
}

const (
	listenAddr = ":2112"
)

func main() {
	prometheus.MustRegister(diskUsage)
	go updateMetrics()

	fmt.Printf("Listening on %q", listenAddr)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(listenAddr, nil)
}
