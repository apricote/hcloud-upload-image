package cmd

import (
	_ "embed"
	"fmt"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"

	"github.com/apricote/hcloud-upload-image/hcloudimages/v2"
	"github.com/apricote/hcloud-upload-image/hcloudimages/v2/contextlogger"
)

const (
	uploadFlagArchitecture = "architecture"
	uploadFlagServerType   = "server-type"
	uploadFlagDescription  = "description"
	uploadFlagLabels       = "labels"
	uploadFlagLocation     = "location"
)

//go:embed upload.md
var uploadLongDescription string

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload (--image-path=<local-path> | --image-url=<url>) --architecture=<x86|arm>",
	Short: "Upload the specified disk image into your Hetzner Cloud project.",
	Long:  uploadLongDescription,
	Example: `  hcloud-upload-image upload --image-path /home/you/images/custom-linux-image-x86.bz2 --architecture x86 --compression bz2 --description "My super duper custom linux"
  hcloud-upload-image upload --image-url https://examples.com/image-arm.raw --architecture arm --labels foo=bar,version=latest
  hcloud-upload-image upload --image-url https://examples.com/image-x86.qcow2 --architecture x86 --format qcow2`,
	DisableAutoGenTag: true,

	GroupID: "primary",

	PreRun: initClient,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := contextlogger.From(ctx)

		writeOptions, err := parseAndValidateWriteOptions(ctx, cmd.Flags())
		if err != nil {
			return err
		}

		architecture, _ := cmd.Flags().GetString(uploadFlagArchitecture)
		serverType, _ := cmd.Flags().GetString(uploadFlagServerType)
		description, _ := cmd.Flags().GetString(uploadFlagDescription)
		labels, _ := cmd.Flags().GetStringToString(uploadFlagLabels)
		location, _ := cmd.Flags().GetString(uploadFlagLocation)

		options := hcloudimages.UploadOptions{
			WriteOptions: writeOptions,
			Description:  hcloud.Ptr(description),
			Labels:       labels,
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

	registerWriteOptions(uploadCmd)

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
