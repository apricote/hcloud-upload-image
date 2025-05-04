# `hcloud-upload-image`

<p align="center">
  Quickly upload any raw disk images into your <a href="https://hetzner.com/cloud" target="_blank">Hetzner Cloud</a> projects!
</p>

<p align="center">
  <a href="https://apricote.github.io/hcloud-upload-image" target="_blank"><img src="https://img.shields.io/badge/Documentation-brightgreen?style=flat-square" alt="Badge: Documentation"/></a>
  <a href="https://github.com/apricote/hcloud-upload-image/releases" target="_blank"><img src="https://img.shields.io/github/v/release/apricote/hcloud-upload-image?sort=semver&display_name=release&style=flat-square&color=green" alt="Badge: Stable Release"/></a>
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="Badge: License MIT"/>
</p>


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

We provide pre-built `deb`, `rpm` and `apk` packages. Alternatively we also provide the binaries directly.

Check out the [GitHub release artifacts](https://github.com/apricote/hcloud-upload-image/releases/latest) for all of these files and archives.

##### Arch Linux

You can get [`hcloud-upload-image-bin`](https://aur.archlinux.org/packages/hcloud-upload-image-bin) from the AUR.

Use your preferred wrapper to install:

```shell
yay -S hcloud-upload-image-bin
```

#### `go install`

If you already have a recent Go toolchain installed, you can build & install the binary from source:

```shell
go install github.com/apricote/hcloud-upload-image@latest
```

#### Docker

There is a docker image published at `ghcr.io/apricote/hcloud-upload-image`.

```shell
docker run --rm -e HCLOUD_TOKEN="<your token>" ghcr.io/apricote/hcloud-upload-image:latest <command>
```

#### Usage

```shell
export HCLOUD_TOKEN="<your token>"
hcloud-upload-image upload \
  --image-url "https://example.com/disk-image-x86.raw.bz2" \
  --architecture x86 \
  --compression bz2
```

To learn more, you can use the embedded help output or check out the [CLI help pages in this repository](docs/reference/cli/hcloud-upload-image.md).:

```shell
hcloud-upload-image --help
hcloud-upload-image upload --help
hcloud-upload-image cleanup --help
```

### Go Library

The functionality to upload images is also exposed in the library `hcloudimages`! Check out the [reference documentation](https://pkg.go.dev/github.com/apricote/hcloud-upload-image/hcloudimages) for more details.

#### Install

```shell
go get github.com/apricote/hcloud-upload-image/hcloudimages
```

#### Usages

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
