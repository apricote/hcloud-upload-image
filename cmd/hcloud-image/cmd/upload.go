package cmd

import (
	"fmt"
	"net/url"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"

	hcloud_upload_image "github.com/apricote/hcloud-upload-image"
	"github.com/apricote/hcloud-upload-image/contextlogger"
)

const (
	uploadFlagImageURL     = "image-url"
	uploadFlagCompression  = "compression"
	uploadFlagArchitecture = "architecture"
	uploadFlagDescription  = "description"
	uploadFlagLabels       = "labels"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload the specified disk image into your Hetzner Cloud project.",
	Long: `This command implements a fake "upload", by going through a real server and snapshots.
This does cost a bit of money for the server.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := contextlogger.From(ctx)

		imageURLString, _ := cmd.Flags().GetString(uploadFlagImageURL)
		imageCompression, _ := cmd.Flags().GetString(uploadFlagCompression)
		architecture, _ := cmd.Flags().GetString(uploadFlagArchitecture)
		description, _ := cmd.Flags().GetString(uploadFlagDescription)
		labels, _ := cmd.Flags().GetStringToString(uploadFlagLabels)

		imageURL, err := url.Parse(imageURLString)
		if err != nil {
			return fmt.Errorf("unable to parse url from --%s=%q: %w", uploadFlagImageURL, imageURLString, err)
		}

		image, err := client.Upload(ctx, hcloud_upload_image.UploadOptions{
			ImageURL:         imageURL,
			ImageCompression: hcloud_upload_image.Compression(imageCompression),
			Architecture:     hcloud.Architecture(architecture),
			Description:      hcloud.Ptr(description),
			Labels:           labels,
		})
		if err != nil {
			return fmt.Errorf("failed to upload the image: %w", err)
		}

		logger.InfoContext(ctx, "Successfully uploaded the image!", "image", image.ID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().String(uploadFlagImageURL, "", "Remote URL of the disk image that should be uploaded (required)")
	_ = uploadCmd.MarkFlagRequired(uploadFlagImageURL)

	uploadCmd.Flags().String(uploadFlagCompression, "", "Type of compression that was used on the disk image")
	_ = uploadCmd.RegisterFlagCompletionFunc(
		uploadFlagCompression,
		cobra.FixedCompletions([]string{string(hcloud_upload_image.CompressionBZ2)}, cobra.ShellCompDirectiveNoFileComp),
	)

	uploadCmd.Flags().String(uploadFlagArchitecture, "", "CPU Architecture of the disk image. Choices: x86|arm")
	_ = uploadCmd.RegisterFlagCompletionFunc(
		uploadFlagArchitecture,
		cobra.FixedCompletions([]string{string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)}, cobra.ShellCompDirectiveNoFileComp),
	)
	_ = uploadCmd.MarkFlagRequired(uploadFlagArchitecture)

	uploadCmd.Flags().String(uploadFlagDescription, "", "Description for the resulting Image")

	uploadCmd.Flags().StringToString(uploadFlagLabels, map[string]string{}, "Labels for the resulting Image")
}
