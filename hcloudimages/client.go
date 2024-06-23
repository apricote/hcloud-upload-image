package hcloudimages

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"golang.org/x/crypto/ssh"

	"github.com/apricote/hcloud-upload-image/hcloudimages/contextlogger"
	"github.com/apricote/hcloud-upload-image/hcloudimages/internal/actionutil"
	"github.com/apricote/hcloud-upload-image/hcloudimages/internal/control"
	"github.com/apricote/hcloud-upload-image/hcloudimages/internal/labelutil"
	"github.com/apricote/hcloud-upload-image/hcloudimages/internal/randomid"
	"github.com/apricote/hcloud-upload-image/hcloudimages/internal/sshkey"
	"github.com/apricote/hcloud-upload-image/hcloudimages/internal/sshsession"
)

const (
	CreatedByLabel = "apricote.de/created-by"
	CreatedByValue = "hcloud-upload-image"

	resourcePrefix = "hcloud-upload-image-"
)

var (
	DefaultLabels = map[string]string{
		CreatedByLabel: CreatedByValue,
	}

	serverTypePerArchitecture = map[hcloud.Architecture]*hcloud.ServerType{
		hcloud.ArchitectureX86: {Name: "cx22"},
		hcloud.ArchitectureARM: {Name: "cax11"},
	}

	defaultImage      = &hcloud.Image{Name: "ubuntu-24.04"}
	defaultLocation   = &hcloud.Location{Name: "fsn1"}
	defaultRescueType = hcloud.ServerRescueTypeLinux64

	defaultSSHDialTimeout = 1 * time.Minute
)

type UploadOptions struct {
	// ImageURL must be publicly available. The instance will download the image from this endpoint.
	ImageURL *url.URL

	// ImageReader
	ImageReader io.Reader

	// ImageCompression describes the compression of the referenced image file. It defaults to [CompressionNone]. If
	// set to anything else, the file will be decompressed before written to the disk.
	ImageCompression Compression

	// Possible future additions:
	// ImageSignatureVerification
	// ImageLocalPath
	// ImageType (RawDiskImage, ISO, qcow2, ...)

	// Architecture should match the architecture of the Image. This decides if the Snapshot can later be
	// used with [hcloud.ArchitectureX86] or [hcloud.ArchitectureARM] servers.
	//
	// Internally this decides what server type is used for the temporary server.
	//
	// Optional if [UploadOptions.ServerType] is set.
	Architecture hcloud.Architecture

	// ServerType can be optionally set to override the default server type for the architecture.
	// Situations where this makes sense:
	//
	//   - Your image is larger than the root disk of the default server types.
	//   - The default server type is no longer available, or not temporarily out of stock.
	ServerType *hcloud.ServerType

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
	CompressionXZ   Compression = "xz"

	// Possible future additions:
	// zip,zstd
)

// NewClient instantiates a new client. It requires a working [*hcloud.Client] to interact with the Hetzner Cloud API.
func NewClient(c *hcloud.Client) *Client {
	return &Client{
		c: c,
	}
}

type Client struct {
	c *hcloud.Client
}

// Upload the specified image into a snapshot on Hetzner Cloud.
//
// As the Hetzner Cloud API has no direct way to upload images, we create a temporary server,
// overwrite the root disk and take a snapshot of that disk instead.
//
// The temporary server costs money. If the upload fails, we might be unable to delete the server. Check out
// CleanupTempResources for a helper in this case.
func (s *Client) Upload(ctx context.Context, options UploadOptions) (*hcloud.Image, error) {
	logger := contextlogger.From(ctx).With(
		"library", "hcloudimages",
		"method", "upload",
	)

	id, err := randomid.Generate()
	if err != nil {
		return nil, err
	}
	logger = logger.With("run-id", id)
	// For simplicity, we use the name random name for SSH Key + Server
	resourceName := resourcePrefix + id
	labels := labelutil.Merge(DefaultLabels, options.Labels)

	// 1. Create SSH Key
	logger.InfoContext(ctx, "# Step 1: Generating SSH Key")
	publicKey, privateKey, err := sshkey.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate temporary ssh key pair: %w", err)
	}

	key, _, err := s.c.SSHKey.Create(ctx, hcloud.SSHKeyCreateOpts{
		Name:      resourceName,
		PublicKey: string(publicKey),
		Labels:    labels,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to submit temporary ssh key to API: %w", err)
	}
	logger.DebugContext(ctx, "Uploaded ssh key", "ssh-key-id", key.ID)
	defer func() {
		// Cleanup SSH Key
		if options.DebugSkipResourceCleanup {
			logger.InfoContext(ctx, "Cleanup: Skipping cleanup of temporary ssh key")
			return
		}

		logger.InfoContext(ctx, "Cleanup: Deleting temporary ssh key")

		_, err := s.c.SSHKey.Delete(ctx, key)
		if err != nil {
			logger.WarnContext(ctx, "Cleanup: ssh key could not be deleted", "error", err)
			// TODO
		}
	}()

	// 2. Create Server
	logger.InfoContext(ctx, "# Step 2: Creating Server")
	var serverType *hcloud.ServerType
	if options.ServerType != nil {
		serverType = options.ServerType
	} else {
		var ok bool
		serverType, ok = serverTypePerArchitecture[options.Architecture]
		if !ok {
			return nil, fmt.Errorf("unknown architecture %q, valid options: %q, %q", options.Architecture, hcloud.ArchitectureX86, hcloud.ArchitectureARM)
		}
	}

	logger.DebugContext(ctx, "creating server with config",
		"image", defaultImage.Name,
		"location", defaultLocation.Name,
		"serverType", serverType.Name,
	)
	serverCreateResult, _, err := s.c.Server.Create(ctx, hcloud.ServerCreateOpts{
		Name:       resourceName,
		ServerType: serverType,

		// Not used, but without this the user receives an email with a password for every created server
		SSHKeys: []*hcloud.SSHKey{key},

		// We need to enable rescue system first
		StartAfterCreate: hcloud.Ptr(false),
		// Image will never be booted, we only boot into rescue system
		Image:    defaultImage,
		Location: defaultLocation,
		Labels:   labels,
	})
	if err != nil {
		return nil, fmt.Errorf("creating the temporary server failed: %w", err)
	}
	logger = logger.With("server", serverCreateResult.Server.ID)
	logger.DebugContext(ctx, "Created Server")

	logger.DebugContext(ctx, "waiting on actions")
	err = s.c.Action.WaitFor(ctx, append(serverCreateResult.NextActions, serverCreateResult.Action)...)
	if err != nil {
		return nil, fmt.Errorf("creating the temporary server failed: %w", err)
	}
	logger.DebugContext(ctx, "actions finished")

	server := serverCreateResult.Server
	defer func() {
		// Cleanup Server
		if options.DebugSkipResourceCleanup {
			logger.InfoContext(ctx, "Cleanup: Skipping cleanup of temporary server")
			return
		}

		logger.InfoContext(ctx, "Cleanup: Deleting temporary server")

		_, _, err := s.c.Server.DeleteWithResult(ctx, server)
		if err != nil {
			logger.WarnContext(ctx, "Cleanup: server could not be deleted", "error", err)
		}
	}()

	// 3. Activate Rescue System
	logger.InfoContext(ctx, "# Step 3: Activating Rescue System")
	enableRescueResult, _, err := s.c.Server.EnableRescue(ctx, server, hcloud.ServerEnableRescueOpts{
		Type:    defaultRescueType,
		SSHKeys: []*hcloud.SSHKey{key},
	})
	if err != nil {
		return nil, fmt.Errorf("enabling the rescue system on the temporary server failed: %w", err)
	}

	logger.DebugContext(ctx, "rescue system requested, waiting on action")

	err = s.c.Action.WaitFor(ctx, enableRescueResult.Action)
	if err != nil {
		return nil, fmt.Errorf("enabling the rescue system on the temporary server failed: %w", err)
	}
	logger.DebugContext(ctx, "action finished, rescue system enabled")

	// 4. Boot Server
	logger.InfoContext(ctx, "# Step 4: Booting Server")
	powerOnAction, _, err := s.c.Server.Poweron(ctx, server)
	if err != nil {
		return nil, fmt.Errorf("starting the temporary server failed: %w", err)
	}

	logger.DebugContext(ctx, "boot requested, waiting on action")

	err = s.c.Action.WaitFor(ctx, powerOnAction)
	if err != nil {
		return nil, fmt.Errorf("starting the temporary server failed: %w", err)
	}
	logger.DebugContext(ctx, "action finished, server is booting")

	// 5. Open SSH Session
	logger.InfoContext(ctx, "# Step 5: Opening SSH Connection")
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("parsing the automatically generated temporary private key failed: %w", err)
	}

	sshClientConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// There is no way to get the host key of the rescue system beforehand
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         defaultSSHDialTimeout,
	}

	// the server needs some time until its properly started and ssh is available
	var sshClient *ssh.Client

	err = control.Retry(
		contextlogger.New(ctx, logger.With("operation", "ssh")),
		10,
		func() error {
			var err error
			logger.DebugContext(ctx, "trying to connect to server", "ip", server.PublicNet.IPv4.IP)
			sshClient, err = ssh.Dial("tcp", server.PublicNet.IPv4.IP.String()+":ssh", sshClientConfig)
			return err
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to ssh into temporary server: %w", err)
	}
	defer sshClient.Close()

	// 6. SSH On Server: Download Image, Decompress, Write to Root Disk
	logger.InfoContext(ctx, "# Step 6: Downloading image and writing to disk")
	cmd := ""
	if options.ImageURL != nil {
		cmd += fmt.Sprintf("wget --no-verbose -O - %q | ", options.ImageURL.String())
	}

	if options.ImageCompression != CompressionNone {
		switch options.ImageCompression {
		case CompressionBZ2:
			cmd += "bzip2 -cd | "
		case CompressionXZ:
			cmd += "xz -cd | "
		default:
			return nil, fmt.Errorf("unknown compression: %q", options.ImageCompression)
		}
	}

	cmd += "dd of=/dev/sda bs=4M && sync"

	// Make sure that we fail early, ie. if the image url does not work.
	// the pipefail does not work correctly without wrapping in bash.
	cmd = fmt.Sprintf("bash -c 'set -euo pipefail && %s'", cmd)
	logger.DebugContext(ctx, "running download, decompress and write to disk command", "cmd", cmd)

	output, err := sshsession.Run(sshClient, cmd, options.ImageReader)
	logger.InfoContext(ctx, "# Step 6: Finished writing image to disk")
	logger.DebugContext(ctx, string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to download and write the image: %w", err)
	}

	// 7. SSH On Server: Shutdown
	logger.InfoContext(ctx, "# Step 7: Shutting down server")
	_, err = sshsession.Run(sshClient, "shutdown now", nil)
	if err != nil {
		// TODO Verify if shutdown error, otherwise return
		logger.WarnContext(ctx, "shutdown returned error", "err", err)
	}

	// 8. Create Image from Server
	logger.InfoContext(ctx, "# Step 8: Creating Image")
	createImageResult, _, err := s.c.Server.CreateImage(ctx, server, &hcloud.ServerCreateImageOpts{
		Type:        hcloud.ImageTypeSnapshot,
		Description: options.Description,
		Labels:      labels,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}
	logger.DebugContext(ctx, "image creation requested, waiting on action")

	err = s.c.Action.WaitFor(ctx, createImageResult.Action)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}
	logger.DebugContext(ctx, "action finished, image was created")

	image := createImageResult.Image
	logger.InfoContext(ctx, "# Image was created", "image", image.ID)

	// Resource cleanup is happening in `defer`
	return image, nil
}

// CleanupTempResources tries to delete any resources that were left over from previous calls to [Client.Upload].
// Upload tries to clean up any temporary resources it created at runtime, but might fail at any point.
// You can then use this command to make sure that all temporary resources are removed from your project.
//
// This method tries to delete any server or ssh keys that match the [DefaultLabels]
func (s *Client) CleanupTempResources(ctx context.Context) error {
	logger := contextlogger.From(ctx).With(
		"library", "hcloudimages",
		"method", "cleanup",
	)

	selector := labelutil.Selector(DefaultLabels)
	logger = logger.With("selector", selector)

	logger.InfoContext(ctx, "# Cleaning up Servers")
	err := s.cleanupTempServers(ctx, logger, selector)
	if err != nil {
		return fmt.Errorf("failed to clean up all servers: %w", err)
	}
	logger.DebugContext(ctx, "cleaned up all servers")

	logger.InfoContext(ctx, "# Cleaning up SSH Keys")
	err = s.cleanupTempSSHKeys(ctx, logger, selector)
	if err != nil {
		return fmt.Errorf("failed to clean up all ssh keys: %w", err)
	}
	logger.DebugContext(ctx, "cleaned up all ssh keys")

	return nil
}

func (s *Client) cleanupTempServers(ctx context.Context, logger *slog.Logger, selector string) error {
	servers, err := s.c.Server.AllWithOpts(ctx, hcloud.ServerListOpts{ListOpts: hcloud.ListOpts{
		LabelSelector: selector,
	}})
	if err != nil {
		return fmt.Errorf("failed to list servers: %w", err)
	}

	if len(servers) == 0 {
		logger.InfoContext(ctx, "No servers found")
		return nil
	}
	logger.InfoContext(ctx, "removing servers", "count", len(servers))

	errs := []error{}
	actions := make([]*hcloud.Action, 0, len(servers))

	for _, server := range servers {
		result, _, err := s.c.Server.DeleteWithResult(ctx, server)
		if err != nil {
			errs = append(errs, err)
			logger.WarnContext(ctx, "failed to delete server", "server", server.ID, "error", err)
			continue
		}

		actions = append(actions, result.Action)
	}

	successActions, errorActions, err := actionutil.Settle(ctx, &s.c.Action, actions...)
	if err != nil {
		return fmt.Errorf("failed to wait for server delete: %w", err)
	}

	if len(successActions) > 0 {
		ids := make([]int64, 0, len(successActions))
		for _, action := range successActions {
			for _, resource := range action.Resources {
				if resource.Type == hcloud.ActionResourceTypeServer {
					ids = append(ids, resource.ID)
				}
			}
		}

		logger.InfoContext(ctx, "successfully deleted servers", "servers", ids)
	}

	if len(errorActions) > 0 {
		for _, action := range errorActions {
			errs = append(errs, action.Error())
		}
	}

	if len(errs) > 0 {
		// The returned message contains no info about the server IDs which failed
		return fmt.Errorf("failed to delete some of the servers: %w", errors.Join(errs...))
	}

	return nil
}

func (s *Client) cleanupTempSSHKeys(ctx context.Context, logger *slog.Logger, selector string) error {
	keys, _, err := s.c.SSHKey.List(ctx, hcloud.SSHKeyListOpts{ListOpts: hcloud.ListOpts{
		LabelSelector: selector,
	}})
	if err != nil {
		return fmt.Errorf("failed to list keys: %w", err)
	}

	if len(keys) == 0 {
		logger.InfoContext(ctx, "No ssh keys found")
		return nil
	}

	errs := []error{}
	for _, key := range keys {
		_, err := s.c.SSHKey.Delete(ctx, key)
		if err != nil {
			errs = append(errs, err)
			logger.WarnContext(ctx, "failed to delete ssh key", "ssh-key", key.ID, "error", err)
			continue
		}
	}

	if len(errs) > 0 {
		// The returned message contains no info about the server IDs which failed
		return fmt.Errorf("failed to delete some of the ssh keys: %w", errors.Join(errs...))
	}

	return nil
}
