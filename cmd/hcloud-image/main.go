package main

import (
	"log/slog"
	"os"

	"github.com/apricote/hcloud-upload-image/cmd/hcloud-image/cmd"
	"github.com/apricote/hcloud-upload-image/cmd/hcloud-image/ui"
)

func init() {
	slog.SetDefault(slog.New(ui.NewHandler(os.Stdout, &ui.HandlerOptions{Level: slog.LevelDebug})))
}

func main() {
	cmd.Execute()
}
