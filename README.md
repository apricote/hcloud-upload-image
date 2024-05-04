# `hcloud-upload-image`

Quickly upload any raw disk images into your [Hetzner Cloud](https://hetzner.com/cloud) projects!

## About

The [Hetzner Cloud API](https://docs.hetzner.cloud/) does not support uploading disk images directly, and it only
provides a limited set of default images. The only option for custom disk images that users have is by taking a
"snapshot" of an existing servers root disk. These can then be used to create new servers.

To create a completely custom disk image, users have to follow these steps:

1. Create server with the correct server type
2. Enable rescue system for the server
3. Boot the server
4. Download the disk image from within the rescue system
5. Write disk image to servers root disk
6. Shut down the server
7. Take a snapshot of the servers root disk
8. Delete the server

This is an annoyingly long process. Many users have automated this with [Packer](https://www.packer.io/) &
[`packer-plugin-hcloud`](https://github.com/hetznercloud/packer-plugin-hcloud/) before, but Packer offers a lot of
additional complexity to wrap your head around.

This repository provides a simple CLI tool & Go library to do the above.

## Getting Started

### CLI

#### Binary

> TODO

#### `go install`

If you already have a recent Go toolchain installed, you can build & install the binary from source:

```shell
go install github.com/apricote/hcloud-upload-image
```

#### Usage

```shell
export HCLOUD_TOKEN="<your token>"
hcloud-upload-image upload \
  --image-url "https://example.com/disk-image-x86.raw.bz2" \
  --architecture x86 \
  --compression bz2
```

To learn more, you can use the embedded help output:

```shell
hcloud-upload-image --help
hcloud-upload-image upload --help
hcloud-upload-image cleanup --help
```

### Go Library

The functionality to upload images is also exposed as a library! Check out the [reference documentation](https://pkg.go.dev/github.com/apricote/hcloud-upload-image/hcloudimages).

To install:

```shell
go get github.com/apricote/hcloud-upload-image
```

Example Usage:

```go
package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"

	"github.com/apricote/hcloud-upload-image/hcloudimages"
)

func main() {
	client := hcloudimages.NewClient(
		hcloud.NewClient(hcloud.WithToken("<your token>")),
	)

	imageURL, err := url.Parse("https://example.com/disk-image-x86.raw.bz2")
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
```

## Contributing

If you have any questions, feedback or ideas, feel free to open an issue or pull request.

## License

This project is licensed under the MIT license, unless the file explicitly specifies another license.

## Support Disclaimer

This is not an official Hetzner Cloud product in any way and Hetzner Cloud does not provide support for this.
