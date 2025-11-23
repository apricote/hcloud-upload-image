# hcloud-upload-image

<p align="center">
  Quickly upload any raw disk images into your <a href="https://hetzner.com/cloud" target="_blank">Hetzner Cloud</a> projects!
</p>

<p align="center">
  <a href="https://apricote.github.io/hcloud-upload-image" target="_blank"><img src="https://img.shields.io/badge/Documentation-brightgreen?style=flat-square" alt="Badge: Documentation"/></a>
  <a href="https://github.com/apricote/hcloud-upload-image/releases" target="_blank"><img src="https://img.shields.io/github/v/release/apricote/hcloud-upload-image?sort=semver&display_name=release&style=flat-square&color=green" alt="Badge: Stable Release"/></a>
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="Badge: License MIT"/>
</p>


## About

The [Hetzner Cloud API](https://docs.hetzner.cloud/) does not support uploading disk images directly and only provides a limited set of default images. The only option for custom disk images is to take a snapshot of an existing server’s root disk. These snapshots can then be used to create new servers.

To create a completely custom disk image, users need to follow these steps:

1. Create a server with the correct server type
2. Enable the rescue system for the server
3. Boot the server
4. Download the disk image from within the rescue system
5. Write the disk image to the server’s root disk
6. Shut down the server
7. Take a snapshot of the server’s root disk
8. Delete the server

This is a frustratingly long process. Many users have automated it with [Packer](https://www.packer.io/) and [`packer-plugin-hcloud`](https://github.com/hetznercloud/packer-plugin-hcloud/), but Packer introduces additional complexity that can be difficult to manage.

This repository provides a simple CLI tool and Go library to streamline the process.

## Getting Started

### CLI

#### Binary

We provide pre-built `deb`, `rpm` and `apk` packages. Alternatively we also provide the binaries directly.

Check out the [GitHub release artifacts](https://github.com/apricote/hcloud-upload-image/releases/latest) for all of these files and archives.

#### Arch Linux

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

#### Nix/NixOS

To run directly without installing (assumes flakes are enabled):

```shell
# Run the application directly
nix run github:apricote/hcloud-upload-image

# Start a shell with `hcloud-upload-image` in $PATH
nix shell github:apricote/hcloud-upload-image
```

To install on your system (assumes flakes are enabled):

```nix
# flake.nix
{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.hcloud-upload-image.url = "github:apricote/hcloud-upload-image";
  inputs.hcloud-upload-image.inputs.nixpkgs.follows = "nixpkgs";

  outputs = { self, nixpkgs, hcloud-upload-image }: {
    nixosConfigurations.my-system = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        ./configuration.nix
        {
          environment.systemPackages = [
            hcloud-upload-image.packages.x86_64-linux.default
          ];
        }
      ];
    };
  };
}
```

To install on your system (using a non-flake version manager):

```shell
# Using npins
npins add github apricote hcloud-upload-image

# Using niv
niv add apricote/hcloud-upload-image
```

Then in your Nix expressions:

```nix
let
  sources = import ./npins;             # For npins
  # sources = import ./nix/sources.nix; # For niv
in
(pkgs.callPackage sources.hcloud-upload-image {})
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
