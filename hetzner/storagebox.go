package hetzner

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"

	"github.com/imroc/req/v3"
)

var (
	Client *req.Client
)

type APIReponse struct {
	StorageBoxes []StorageBox `json:"storage_boxes"`
	Meta         ResponseMeta `json:"meta"`
}

type ResponseMeta struct {
	Pagination struct {
		Page         int `json:"page"`
		PerPage      int `json:"per_page"`
		PreviousPage int `json:"previous_page"`
		NextPage     int `json:"next_page"`
		LastPage     int `json:"last_page"`
		TotalEntries int `json:"total_entries"`
	} `json:"pagination"`
}
type StorageBox struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Status         string `json:"status"`
	Name           string `json:"name"`
	StorageBoxType struct {
		ID                     int    `json:"id"`
		Name                   string `json:"name"`
		Description            string `json:"description"`
		SnapshotLimit          int    `json:"snapshot_limit"`
		AutomaticSnapshotLimit int    `json:"automatic_snapshot_limit"`
		SubaccountsLimit       int    `json:"subaccounts_limit"`
		Size                   int    `json:"size"`
		Prices                 []struct {
			Location    string `json:"location"`
			PriceHourly struct {
				Net   string `json:"net"`
				Gross string `json:"gross"`
			} `json:"price_hourly"`
			PriceMonthly struct {
				Net   string `json:"net"`
				Gross string `json:"gross"`
			} `json:"price_monthly"`
			SetupFee struct {
				Net   string `json:"net"`
				Gross string `json:"gross"`
			} `json:"setup_fee"`
		} `json:"prices"`
		Deprecation struct {
			UnavailableAfter time.Time `json:"unavailable_after"`
			Announced        time.Time `json:"announced"`
		} `json:"deprecation"`
	} `json:"storage_box_type"`
	Location struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Country     string  `json:"country"`
		City        string  `json:"city"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		NetworkZone string  `json:"network_zone"`
	} `json:"location"`
	AccessSettings struct {
		ReachableExternally bool `json:"reachable_externally"`
		SambaEnabled        bool `json:"samba_enabled"`
		SSHEnabled          bool `json:"ssh_enabled"`
		WebdavEnabled       bool `json:"webdav_enabled"`
		ZfsEnabled          bool `json:"zfs_enabled"`
	} `json:"access_settings"`
	Server string `json:"server"`
	System string `json:"system"`
	Stats  struct {
		Size          int `json:"size"`
		SizeData      int `json:"size_data"`
		SizeSnapshots int `json:"size_snapshots"`
	} `json:"stats"`
	Labels struct {
		Environment  string `json:"environment"`
		ExampleComMy string `json:"example.com/my"`
		JustAKey     string `json:"just-a-key"`
	} `json:"labels"`
	Protection struct {
		Delete bool `json:"delete"`
	} `json:"protection"`
	SnapshotPlan struct {
		MaxSnapshots int         `json:"max_snapshots"`
		Minute       interface{} `json:"minute"`
		Hour         interface{} `json:"hour"`
		DayOfWeek    interface{} `json:"day_of_week"`
		DayOfMonth   interface{} `json:"day_of_month"`
	} `json:"snapshot_plan"`
	Created time.Time `json:"created"`
}

func GetBoxes() ([]StorageBox, error) {

	resp, err := Client.R().Get("https://api.hetzner.com/v1/storage_boxes")
	if err != nil {
		return nil, err
	}

	bodyText, err := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Printf("API Error %d: %s\n", resp.StatusCode, bodyText)
		return nil, errors.New("HTTP Error")
	}

	var apiResponse APIReponse
	err = json.Unmarshal(bodyText, &apiResponse)
	if err != nil {
		log.Fatal(err)
	}

	return apiResponse.StorageBoxes, nil
}
