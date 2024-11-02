package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"

	"github.com/apricote/hcloud-upload-image/hcloudimages"
	"github.com/apricote/hcloud-upload-image/hcloudimages/backoff"
	"github.com/apricote/hcloud-upload-image/hcloudimages/contextlogger"
	"github.com/apricote/hcloud-upload-image/internal/ui"
	"github.com/apricote/hcloud-upload-image/internal/version"
)

const (
	flagVerbose = "verbose"
)

var (
	// 1 activates slog debug output
	// 2 activates hcloud-go debug output
	verbose int
)

// The pre-authenticated client. Set in the root command PersistentPreRun
var client *hcloudimages.Client

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:               "hcloud-upload-image",
	Short:             `Manage custom OS images on Hetzner Cloud.`,
	Long:              `Manage custom OS images on Hetzner Cloud.`,
	SilenceUsage:      true,
	DisableAutoGenTag: true,

	Version: version.Version,

	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		ctx := cmd.Context()

		slog.SetDefault(initLogger())

		// Add logger to command context
		logger := slog.Default()
		ctx = contextlogger.New(ctx, logger)
		cmd.SetContext(ctx)
	},
}

func initLogger() *slog.Logger {
	logLevel := slog.LevelInfo
	if verbose >= 1 {
		logLevel = slog.LevelDebug
	}

	return slog.New(ui.NewHandler(os.Stdout, &ui.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove attributes that are unnecessary for the cli context
			if a.Key == "library" || a.Key == "method" {
				return slog.Attr{}
			}

			return a
		},
	}))

}

func initClient(cmd *cobra.Command, _ []string) {
	if client != nil {
		// Only init if not set.
		// Theoretically this is not safe against data races and should use [sync.Once], but :shrug:
		return
	}

	ctx := cmd.Context()

	logger := contextlogger.From(ctx)
	// Build hcloud-go client
	if os.Getenv("HCLOUD_TOKEN") == "" {
		logger.ErrorContext(ctx, "You need to set the HCLOUD_TOKEN environment variable to your Hetzner Cloud API Token.")
		os.Exit(1)
	}

	opts := []hcloud.ClientOption{
		hcloud.WithToken(os.Getenv("HCLOUD_TOKEN")),
		hcloud.WithApplication("hcloud-upload-image", version.Version),
		hcloud.WithPollBackoffFunc(backoff.ExponentialBackoffWithLimit(2, 1*time.Second, 30*time.Second)),
	}

	if os.Getenv("HCLOUD_DEBUG") != "" || verbose >= 2 {
		opts = append(opts, hcloud.WithDebugWriter(os.Stderr))
	}

	client = hcloudimages.NewClient(hcloud.NewClient(opts...))
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.SetErrPrefix("\033[1;31mError:")

	RootCmd.PersistentFlags().CountVarP(&verbose, flagVerbose, "v", "verbose debug output, can be specified up to 2 times")

	RootCmd.AddGroup(&cobra.Group{
		ID:    "primary",
		Title: "Primary Commands:",
	})
}
