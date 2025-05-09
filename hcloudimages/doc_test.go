package hcloudimages_test

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"

	"github.com/apricote/hcloud-upload-image/hcloudimages"
)

func ExampleClient_Upload() {
	client := hcloudimages.NewClient(
		hcloud.NewClient(hcloud.WithToken("<your token>")),
	)

	imageURL, err := url.Parse("https://example.com/disk-image.raw.bz2")
	if err != nil {
		panic(err)
	}

	image, err := client.Upload(context.TODO(), hcloudimages.UploadOptions{
		ImageURL:         imageURL,
		ImageCompression: hcloudimages.CompressionBZ2,
		Architecture:     hcloud.ArchitectureX86,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Uploaded Image: %d", image.ID)
}
