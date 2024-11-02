package cmd

import (
	"fmt"
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
	uploadFlagArchitecture = "architecture"
	uploadFlagServerType   = "server-type"
	uploadFlagDescription  = "description"
	uploadFlagLabels       = "labels"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload (--image-path=<local-path> | --image-url=<url>) --architecture=<x86|arm>",
	Short: "Upload the specified disk image into your Hetzner Cloud project.",
	Long: `This command implements a fake "upload", by going through a real server and snapshots.
This does cost a bit of money for the server.`,
	Example: `  hcloud-upload-image upload --image-path /home/you/images/custom-linux-image-x86.bz2 --architecture x86 --compression bz2 --description "My super duper custom linux"
  hcloud-upload-image upload --image-url https://examples.com/image-arm.raw --architecture arm --labels foo=bar,version=latest`,
	DisableAutoGenTag: true,

	GroupID: "primary",

	PreRun: initClient,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := contextlogger.From(ctx)

		imageURLString, _ := cmd.Flags().GetString(uploadFlagImageURL)
		imagePathString, _ := cmd.Flags().GetString(uploadFlagImagePath)
		imageCompression, _ := cmd.Flags().GetString(uploadFlagCompression)
		architecture, _ := cmd.Flags().GetString(uploadFlagArchitecture)
		serverType, _ := cmd.Flags().GetString(uploadFlagServerType)
		description, _ := cmd.Flags().GetString(uploadFlagDescription)
		labels, _ := cmd.Flags().GetStringToString(uploadFlagLabels)

		options := hcloudimages.UploadOptions{
			ImageCompression: hcloudimages.Compression(imageCompression),
			Description:      hcloud.Ptr(description),
			Labels:           labels,
		}

		if imageURLString != "" {
			imageURL, err := url.Parse(imageURLString)
			if err != nil {
				return fmt.Errorf("unable to parse url from --%s=%q: %w", uploadFlagImageURL, imageURLString, err)
			}

			options.ImageURL = imageURL
		} else if imagePathString != "" {
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

	uploadCmd.Flags().String(uploadFlagCompression, "", "Type of compression that was used on the disk image [choices: bz2, xz]")
	_ = uploadCmd.RegisterFlagCompletionFunc(
		uploadFlagCompression,
		cobra.FixedCompletions([]string{string(hcloudimages.CompressionBZ2), string(hcloudimages.CompressionXZ)}, cobra.ShellCompDirectiveNoFileComp),
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
}
