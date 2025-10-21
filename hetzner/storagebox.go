package hetzner

import (
	"context"
	"log"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var (
	Client *hcloud.Client
)

func GetBoxes() ([]*hcloud.StorageBox, error) {

	boxes, err := Client.StorageBox.All(context.Background())
	if err != nil {
		log.Fatalf("error retrieving server: %s\n", err)
	}

	// bodyText, err := io.ReadAll(resp.Body)
	//
	// if resp.StatusCode != 200 {
	// 	log.Printf("API Error %d: %s\n", resp.StatusCode, bodyText)
	// 	return nil, errors.New("HTTP Error")
	// }
	//
	// var apiResponse APIReponse
	// err = json.Unmarshal(bodyText, &apiResponse)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return boxes, nil
}
