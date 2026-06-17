package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/apricote/hcloud-upload-image/hcloudimages/v2"
	"github.com/apricote/hcloud-upload-image/hcloudimages/v2/contextlogger"
)

const (
	writeFlagImageURL    = "image-url"
	writeFlagImagePath   = "image-path"
	writeFlagCompression = "compression"
	writeFlagFormat      = "format"
	writeFlagServer      = "server"
)

func registerWriteOptions(cmd *cobra.Command) {
	cmd.Flags().String(writeFlagImageURL, "", "Remote URL of the disk image")
	cmd.Flags().String(writeFlagImagePath, "", "Local path to the disk image")
	cmd.MarkFlagsMutuallyExclusive(writeFlagImageURL, writeFlagImagePath)
	cmd.MarkFlagsOneRequired(writeFlagImageURL, writeFlagImagePath)

	cmd.Flags().String(writeFlagCompression, "", "Type of compression that was used on the disk image [choices: bz2, xz, zstd]")
	_ = cmd.RegisterFlagCompletionFunc(
		writeFlagCompression,
		cobra.FixedCompletions([]string{string(hcloudimages.CompressionBZ2), string(hcloudimages.CompressionXZ), string(hcloudimages.CompressionZSTD)}, cobra.ShellCompDirectiveNoFileComp),
	)

	cmd.Flags().String(writeFlagFormat, "", "Format of the disk image. [default: raw, choices: qcow2]")
	_ = cmd.RegisterFlagCompletionFunc(
		writeFlagFormat,
		cobra.FixedCompletions([]string{string(hcloudimages.FormatQCOW2)}, cobra.ShellCompDirectiveNoFileComp),
	)
}

func parseAndValidateWriteOptions(ctx context.Context, flags *pflag.FlagSet) (hcloudimages.WriteOptions, error) {
	logger := contextlogger.From(ctx)

	imageURLString, _ := flags.GetString(writeFlagImageURL)
	imagePathString, _ := flags.GetString(writeFlagImagePath)
	imageCompression, _ := flags.GetString(writeFlagCompression)
	imageFormat, _ := flags.GetString(writeFlagFormat)

	options := hcloudimages.WriteOptions{
		ImageCompression: hcloudimages.Compression(imageCompression),
		ImageFormat:      hcloudimages.Format(imageFormat),
	}

	if imageURLString != "" {
		imageURL, err := url.Parse(imageURLString)
		if err != nil {
			return hcloudimages.WriteOptions{}, fmt.Errorf("unable to parse url from --%s=%q: %w", writeFlagImageURL, imageURLString, err)
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
		_ = resp.Body.Close()

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
			return hcloudimages.WriteOptions{}, fmt.Errorf("unable to read file from --%s=%q: %w", writeFlagImagePath, imagePathString, err)
		}

		options.ImageReader = imageFile
	}

	return options, nil
}

//go:embed write-to-disk.md
var writeToDiskLongDescription string

// writeToDiskCmd represents the write-to-disk command
var writeToDiskCmd = &cobra.Command{
	Use:   "write-to-disk (--image-path=<local-path> | --image-url=<url>) --server <id-or-name>",
	Short: "Write the specified disk image to the root disk of the specified server.",
	Long:  writeToDiskLongDescription,
	Example: `  hcloud-upload-image write-to-disk --image-path /home/you/images/custom-linux-image-x86.bz2 --compression bz2 --server my-server
  hcloud-upload-image write-to-disk --image-url https://examples.com/image-arm.raw --server my-arm-server
  hcloud-upload-image write-to-disk --image-url https://examples.com/image-x86.qcow2 --format qcow2 --server my-x86-server`,
	DisableAutoGenTag: true,

	GroupID: "primary",

	PreRun: initClient,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := contextlogger.From(ctx)

		options, err := parseAndValidateWriteOptions(ctx, cmd.Flags())
		if err != nil {
			return err
		}

		serverIDOrName, _ := cmd.Flags().GetString(writeFlagServer)
		options.Server, _, err = hcloudclient.Server.Get(ctx, serverIDOrName)
		if err != nil {
			return fmt.Errorf("could not get server %q: %w", serverIDOrName, err)
		}
		if options.Server == nil {
			return fmt.Errorf("server %q not found", serverIDOrName)
		}

		err = client.WriteToDisk(ctx, options)
		if err != nil {
			return fmt.Errorf("failed to write the image: %w", err)
		}

		logger.InfoContext(ctx, "Successfully wrote the image!")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(writeToDiskCmd)

	registerWriteOptions(writeToDiskCmd)

	writeToDiskCmd.Flags().String(writeFlagServer, "", "ID or name of target server")
	_ = writeToDiskCmd.MarkFlagRequired(writeFlagServer)
	_ = writeToDiskCmd.RegisterFlagCompletionFunc(
		writeFlagServer,
		func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			servers, err := hcloudclient.Server.AllWithOpts(cmd.Context(), hcloud.ServerListOpts{})
			if err != nil {
				return []cobra.Completion{}, cobra.ShellCompDirectiveError
			}

			serverNames := make([]string, len(servers))
			for i, server := range servers {
				serverNames[i] = server.Name
			}

			return serverNames, cobra.ShellCompDirectiveNoFileComp
		},
	)
}
