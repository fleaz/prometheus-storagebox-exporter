package hetzner

import (
	"context"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var (
	Client *hcloud.Client
)

func GetBoxes() ([]*hcloud.StorageBox, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	boxes, err := Client.StorageBox.All(ctx)
	if err != nil {
		return nil, err
	}
	return boxes, nil
}
