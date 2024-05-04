// Package hcloudimages is a library to upload Disk Images into your Hetzner Cloud project.
//
// # Overview
//
// The Hetzner Cloud API does not support uploading disk images directly, and it only provides a limited set of default
// images. The only option for custom disk images that users have is by taking a "snapshot" of an existing servers root
// disk. These can then be used to create new servers.
//
// To create a completely custom disk image, users have to follow these steps:
//
//  1. Create server with the correct server type
//  2. Enable rescue system for the server
//  3. Boot the server
//  4. Download the disk image from within the rescue system
//  5. Write disk image to servers root disk
//  6. Shut down the server
//  7. Take a snapshot of the servers root disk
//  8. Delete the server
//
// This is an annoyingly long process. Many users have automated this with Packer before, but Packer offers a lot of
// additional complexity to understand.
//
// This library is a single call to do the above: [Client.Upload]
//
// # Costs
//
// The temporary server and the snapshot itself cost money. See the [Hetzner Cloud website] for up-to-date pricing
// information.
//
// Usually the upload takes no more than a few minutes of server time, so you will only be billed for the first hour
// (<1ct for most cases). If this process fails, the server might stay around until you manually delete it. In that case
// it continues to cost its hourly price. There is a utility [Client.CleanupTemporaryResources] that removes any
// leftover resources.
//
// # Logging
//
// By default, nothing is logged. As the update process takes a bit of time you might want to gain some insight into
// the process. For this we provide optional logs through [log/slog]. You can set a [log/slog.Logger] in the
// [context.Context] through [github.com/apricote/hcloud-upload-image/hcloudimages/contextlogger.New].
//
// [Hetzner Cloud website]: https://www.hetzner.com/cloud/
package hcloudimages
