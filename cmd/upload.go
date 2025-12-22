package cmd

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"

	"github.com/apricote/hcloud-upload-image/hcloudimages"
	"github.com/apricote/hcloud-upload-image/hcloudimages/contextlogger"
)

const (
	uploadFlagImageURL     = "image-url"
	uploadFlagImagePath    = "image-path"
	uploadFlagCompression  = "compression"
	uploadFlagFormat       = "format"
	uploadFlagArchitecture = "architecture"
	uploadFlagServerType   = "server-type"
	uploadFlagDescription  = "description"
	uploadFlagLabels       = "labels"
	uploadFlagLocation     = "location"
)

//go:embed upload.md
var longDescription string

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload (--image-path=<local-path> | --image-url=<url>) --architecture=<x86|arm>",
	Short: "Upload the specified disk image into your Hetzner Cloud project.",
	Long:  longDescription,
	Example: `  hcloud-upload-image upload --image-path /home/you/images/custom-linux-image-x86.bz2 --architecture x86 --compression bz2 --description "My super duper custom linux"
  hcloud-upload-image upload --image-url https://examples.com/image-arm.raw --architecture arm --labels foo=bar,version=latest
  hcloud-upload-image upload --image-url https://examples.com/image-x86.qcow2 --architecture x86 --format qcow2`,
	DisableAutoGenTag: true,

	GroupID: "primary",

	PreRun: initClient,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := contextlogger.From(ctx)

		imageURLString, _ := cmd.Flags().GetString(uploadFlagImageURL)
		imagePathString, _ := cmd.Flags().GetString(uploadFlagImagePath)
		imageCompression, _ := cmd.Flags().GetString(uploadFlagCompression)
		imageFormat, _ := cmd.Flags().GetString(uploadFlagFormat)
		architecture, _ := cmd.Flags().GetString(uploadFlagArchitecture)
		serverType, _ := cmd.Flags().GetString(uploadFlagServerType)
		description, _ := cmd.Flags().GetString(uploadFlagDescription)
		labels, _ := cmd.Flags().GetStringToString(uploadFlagLabels)
		location, _ := cmd.Flags().GetString(uploadFlagLocation)

		options := hcloudimages.UploadOptions{
			ImageCompression: hcloudimages.Compression(imageCompression),
			ImageFormat:      hcloudimages.Format(imageFormat),
			Description:      hcloud.Ptr(description),
			Labels:           labels,
		}

		if imageURLString != "" {
			imageURL, err := url.Parse(imageURLString)
			if err != nil {
				return fmt.Errorf("unable to parse url from --%s=%q: %w", uploadFlagImageURL, imageURLString, err)
			}

			// Check for image size
			resp, err := http.Head(imageURL.String())
			switch {
			case err != nil:
				logger.DebugContext(ctx, "failed to check for file size, error on request", "err", err)
			case resp.ContentLength == -1:
				logger.DebugContext(ctx, "failed to check for file size, server did not set the Content-Length", "err", err)
			default:
				options.ImageSize = resp.ContentLength
			}

			options.ImageURL = imageURL
		} else if imagePathString != "" {
			stat, err := os.Stat(imagePathString)
			if err != nil {
				logger.DebugContext(ctx, "failed to check for file size, error on stat", "err", err)
			} else {
				options.ImageSize = stat.Size()
			}

			imageFile, err := os.Open(imagePathString)
			if err != nil {
				return fmt.Errorf("unable to read file from --%s=%q: %w", uploadFlagImagePath, imagePathString, err)
			}

			options.ImageReader = imageFile
		}

		if architecture != "" {
			options.Architecture = hcloud.Architecture(architecture)
		} else if serverType != "" {
			options.ServerType = &hcloud.ServerType{Name: serverType}
		}

		if location != "" {
			options.Location = &hcloud.Location{Name: location}
		}

		image, err := client.Upload(ctx, options)
		if err != nil {
			return fmt.Errorf("failed to upload the image: %w", err)
		}

		logger.InfoContext(ctx, "Successfully uploaded the image!", "image", image.ID)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().String(uploadFlagImageURL, "", "Remote URL of the disk image that should be uploaded")
	uploadCmd.Flags().String(uploadFlagImagePath, "", "Local path to the disk image that should be uploaded")
	uploadCmd.MarkFlagsMutuallyExclusive(uploadFlagImageURL, uploadFlagImagePath)
	uploadCmd.MarkFlagsOneRequired(uploadFlagImageURL, uploadFlagImagePath)

	uploadCmd.Flags().String(uploadFlagCompression, "", "Type of compression that was used on the disk image [choices: bz2, xz, zstd]")
	_ = uploadCmd.RegisterFlagCompletionFunc(
		uploadFlagCompression,
		cobra.FixedCompletions([]string{string(hcloudimages.CompressionBZ2), string(hcloudimages.CompressionXZ), string(hcloudimages.CompressionZSTD)}, cobra.ShellCompDirectiveNoFileComp),
	)

	uploadCmd.Flags().String(uploadFlagFormat, "", "Format of the image. [default: raw, choices: qcow2]")
	_ = uploadCmd.RegisterFlagCompletionFunc(
		uploadFlagFormat,
		cobra.FixedCompletions([]string{string(hcloudimages.FormatQCOW2)}, cobra.ShellCompDirectiveNoFileComp),
	)

	uploadCmd.Flags().String(uploadFlagArchitecture, "", "CPU architecture of the disk image [choices: x86, arm]")
	_ = uploadCmd.RegisterFlagCompletionFunc(
		uploadFlagArchitecture,
		cobra.FixedCompletions([]string{string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)}, cobra.ShellCompDirectiveNoFileComp),
	)

	uploadCmd.Flags().String(uploadFlagServerType, "", "Explicitly use this server type to generate the image. Mutually exclusive with --architecture.")

	// Only one of them needs to be set
	uploadCmd.MarkFlagsOneRequired(uploadFlagArchitecture, uploadFlagServerType)
	uploadCmd.MarkFlagsMutuallyExclusive(uploadFlagArchitecture, uploadFlagServerType)

	uploadCmd.Flags().String(uploadFlagDescription, "", "Description for the resulting image")

	uploadCmd.Flags().StringToString(uploadFlagLabels, map[string]string{}, "Labels for the resulting image")

	uploadCmd.Flags().String(uploadFlagLocation, "", "Datacenter location for the temporary server [default: fsn1, choices: fsn1, nbg1, hel1, ash, hil, sin]")
	_ = uploadCmd.RegisterFlagCompletionFunc(
		uploadFlagLocation,
		cobra.FixedCompletions([]string{"fsn1", "nbg1", "hel1", "ash", "hil", "sin"}, cobra.ShellCompDirectiveNoFileComp),
	)
}
