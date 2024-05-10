package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/apricote/hcloud-upload-image/hcloudimages/contextlogger"
)

// cleanupCmd represents the cleanup command
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove any temporary resources that were left over",
	Long: `If the upload fails at any point, there might still exist a server or
ssh key in your Hetzner Cloud project. This command cleans up any resources
that match the label "apricote.de/created-by=hcloud-upload-image".

If you want to see a preview of what would be removed, you can use the official hcloud CLI and run:

$ hcloud server list -l apricote.de/created-by=hcloud-upload-image
$ hcloud ssh-key list -l apricote.de/created-by=hcloud-upload-image

This command does not handle any parallel executions of hcloud-upload-image
and will remove in-use resources if called at the same time.`,

	GroupID: "primary",

	PreRun: initClient,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := contextlogger.From(ctx)

		err := client.CleanupTempResources(ctx)
		if err != nil {
			return fmt.Errorf("failed to clean up temporary resources: %w", err)
		}

		logger.InfoContext(ctx, "Successfully cleaned up all temporary resources!")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
}
