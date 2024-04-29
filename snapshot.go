package hcloud_upload_image

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"golang.org/x/crypto/ssh"

	"github.com/apricote/hcloud-upload-image/util/control"
	"github.com/apricote/hcloud-upload-image/util/randomid"
	"github.com/apricote/hcloud-upload-image/util/sshkey"
	"github.com/apricote/hcloud-upload-image/util/sshsession"
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
		hcloud.ArchitectureX86: {Name: "cx11"},
		hcloud.ArchitectureARM: {Name: "cax11"},
	}

	defaultImage      = &hcloud.Image{Name: "ubuntu-22.04"}
	defaultLocation   = &hcloud.Location{Name: "fsn1"}
	defaultRescueType = hcloud.ServerRescueTypeLinux64

	defaultSSHDialTimeout = 1 * time.Minute
)

func New(client *hcloud.Client) SnapshotClient {
	return &snapshotClient{
		client: client,
	}
}

type snapshotClient struct {
	client *hcloud.Client
}

func (s snapshotClient) Upload(ctx context.Context, options UploadOptions) (*hcloud.Image, error) {
	id, err := randomid.Generate()
	if err != nil {
		return nil, err
	}
	// For simplicity, we use the name random name for SSH Key + Server
	resourceName := resourcePrefix + id

	// 1. Create SSH Key
	log.Print("Step 1: Generating SSH Key")
	publicKey, privateKey, err := sshkey.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate temporary ssh key pair: %w", err)
	}

	key, _, err := s.client.SSHKey.Create(ctx, hcloud.SSHKeyCreateOpts{
		Name:      resourceName,
		PublicKey: string(publicKey),
		Labels:    fullLabels(options.Labels),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to submit temporary ssh key to API: %w", err)
	}
	defer func() {
		// Cleanup SSH Key
		if options.DebugSkipResourceCleanup {
			return
		}

		_, err := s.client.SSHKey.Delete(ctx, key)
		if err != nil {
			// TODO
		}
	}()

	// 2. Create Server
	log.Print("Step 2: Creating Server")
	serverType, ok := serverTypePerArchitecture[options.Architecture]
	if !ok {
		return nil, fmt.Errorf("unknown architecture %q, valid options: %q, %q", options.Architecture, hcloud.ArchitectureX86, hcloud.ArchitectureARM)
	}

	serverCreateResult, _, err := s.client.Server.Create(ctx, hcloud.ServerCreateOpts{
		Name:       resourceName,
		ServerType: serverType,

		// Not used, but without this the user receives an email with a password for every created server
		SSHKeys: []*hcloud.SSHKey{key},

		// We need to enable rescue system first
		StartAfterCreate: hcloud.Ptr(false),
		// Image will never be booted, we only boot into rescue system
		Image:    defaultImage,
		Location: defaultLocation,
		Labels:   fullLabels(options.Labels),
	})
	if err != nil {
		return nil, fmt.Errorf("creating the temporary server failed: %w", err)
	}

	_, err = s.client.Action.WaitForActions(ctx, append(serverCreateResult.NextActions, serverCreateResult.Action)...)
	if err != nil {
		return nil, fmt.Errorf("creating the temporary server failed: %w", err)
	}

	server := serverCreateResult.Server
	defer func() {
		// Cleanup Server
		if options.DebugSkipResourceCleanup {
			return
		}

		_, _, err := s.client.Server.DeleteWithResult(ctx, server)
		if err != nil {
			// TODO
		}
	}()

	// 3. Activate Rescue System
	log.Print("Step 3: Activating Rescue System")
	enableRescueResult, _, err := s.client.Server.EnableRescue(ctx, server, hcloud.ServerEnableRescueOpts{
		Type:    defaultRescueType,
		SSHKeys: []*hcloud.SSHKey{key},
	})
	if err != nil {
		return nil, fmt.Errorf("enabling the rescue system on the temporary server failed: %w", err)
	}

	_, err = s.client.Action.WaitForAction(ctx, enableRescueResult.Action)
	if err != nil {
		return nil, fmt.Errorf("enabling the rescue system on the temporary server failed: %w", err)
	}

	// 4. Boot Server
	log.Print("Step 4: Booting Server")
	powerOnAction, _, err := s.client.Server.Poweron(ctx, server)
	if err != nil {
		return nil, fmt.Errorf("starting the temporary server failed: %w", err)
	}

	_, err = s.client.Action.WaitForAction(ctx, powerOnAction)
	if err != nil {
		return nil, fmt.Errorf("starting the temporary server failed: %w", err)
	}

	// 5. Open SSH Session
	log.Print("Step 5: Opening SSH Connection")
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

	err = control.Retry(ctx, 10, func() error {
		var err error
		sshClient, err = ssh.Dial("tcp", server.PublicNet.IPv4.IP.String()+":ssh", sshClientConfig)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ssh into temporary server: %w", err)
	}
	defer sshClient.Close()

	// 6. SSH On Server: Download Image, Decompress, Write to Root Disk
	log.Print("Step 6: Downloading image and writing to disk")
	decompressionCommand := ""
	if options.ImageCompression != CompressionNone {
		switch options.ImageCompression {
		case CompressionBZ2:
			decompressionCommand += "| bzip2 -cd"
		default:
			return nil, fmt.Errorf("unknown compression: %q", options.ImageCompression)
		}
	}

	output, err := sshsession.Run(sshClient, fmt.Sprintf("wget --no-verbose -O - %q %s | dd of=/dev/sda && sync", options.ImageURL.String(), decompressionCommand))
	log.Print(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to download and write the image: %w", err)
	}

	// 7. SSH On Server: Shutdown
	log.Print("Step 7: Shutting down server")
	_, err = sshsession.Run(sshClient, "shutdown now")
	if err != nil {
		// TODO Verify if shutdown error, otherwise return
		fmt.Printf("shutdown error: %+v", err)
	}

	// 8. Create Image from Server
	log.Print("Step 7: Creating Image")
	createImageResult, _, err := s.client.Server.CreateImage(ctx, server, &hcloud.ServerCreateImageOpts{
		Type:        hcloud.ImageTypeSnapshot,
		Description: options.Description,
		Labels:      fullLabels(options.Labels),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}
	_, err = s.client.Action.WaitForAction(ctx, createImageResult.Action)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	image := createImageResult.Image

	// Resource cleanup is happening in `defer`
	return image, nil
}

func fullLabels(userLabels map[string]string) map[string]string {
	if userLabels == nil {
		userLabels = make(map[string]string, len(DefaultLabels))
	}
	for k, v := range DefaultLabels {
		userLabels[k] = v
	}

	return userLabels
}
