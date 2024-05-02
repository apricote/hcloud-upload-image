package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"

	"github.com/apricote/hcloud-upload-image/hcloudimages"
	"github.com/apricote/hcloud-upload-image/hcloudimages/backoff"
	"github.com/apricote/hcloud-upload-image/hcloudimages/contextlogger"
)

// The pre-authenticated client. Set in the root command PersistentPreRun
var client hcloudimages.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "hcloud-upload-image",
	Long:         `Manage custom OS images on Hetzner Cloud.`,
	SilenceUsage: true,

	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()

		// Add logger to command context
		logger := slog.Default()
		ctx = contextlogger.New(ctx, logger)
		cmd.SetContext(ctx)

		client = newClient(ctx)
	},
}

func newClient(ctx context.Context) hcloudimages.Client {
	logger := contextlogger.From(ctx)
	// Build hcloud-go client
	if os.Getenv("HCLOUD_TOKEN") == "" {
		logger.ErrorContext(ctx, "You need to set the HCLOUD_TOKEN environment variable to your Hetzner Cloud API Token.")
		os.Exit(1)
	}

	opts := []hcloud.ClientOption{
		hcloud.WithToken(os.Getenv("HCLOUD_TOKEN")),
		hcloud.WithApplication("hcloud-image", ""),
		hcloud.WithPollBackoffFunc(backoff.ExponentialBackoffWithLimit(2, 1*time.Second, 30*time.Second)),
	}

	if os.Getenv("HCLOUD_DEBUG") != "" {
		opts = append(opts, hcloud.WithDebugWriter(os.Stderr))
	}

	return hcloudimages.NewClient(hcloud.NewClient(opts...))
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetErrPrefix("\033[1;31mError:")
	rootCmd.SetFlagErrorFunc(func(command *cobra.Command, err error) error {
		return fmt.Errorf("fooo")
	})
}
