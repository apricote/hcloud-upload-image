package hcloud_upload_image

import (
	"context"
	"net/url"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type SnapshotClient interface {
	// Upload the specified image into a snapshot on Hetzner Cloud.
	//
	// As the Hetzner Cloud API has no direct way to upload images, we create a temporary server,
	// overwrite the root disk and take a snapshot of that disk instead.
	Upload(ctx context.Context, options UploadOptions) (*hcloud.Image, error)

	// Possible future additions:
	// List(ctx context.Context) []*hcloud.Image
	// Delete(ctx context.Context, image *hcloud.Image) error
	// CleanupTempResources(ctx context.Context) error
}

type UploadOptions struct {
	// ImageURL must be publicly available. The instance will download the image from this endpoint.
	ImageURL         *url.URL
	ImageCompression Compression
	// ImageSignatureVerification

	// Architecture should match the architecture of the Image. This decides if the Snapshot can later be
	// used with [hcloud.ArchitectureX86] or [hcloud.ArchitectureARM] servers.
	//
	// Internally this decides what server type is used for the temporary server.
	Architecture hcloud.Architecture

	// Description is an optional description that the resulting image (snapshot) will have. There is no way to
	// select images by its description, you should use Labels if you need  to identify your image later.
	Description *string

	// Labels will be added to the resulting image (snapshot). Use these to filter the image list if you
	// need to identify the image later on.
	//
	// We also always add a label `apricote.de/created-by=hcloud-image-upload` ([CreatedByLabel], [CreatedByValue]).
	Labels map[string]string

	// DebugSkipResourceCleanup will skip the cleanup of the temporary SSH Key and Server.
	DebugSkipResourceCleanup bool
}

type Compression string

const (
	CompressionNone Compression = ""
	CompressionBZ2  Compression = "bz2"
	// zip,xz,zstd
)
