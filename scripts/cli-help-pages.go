package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/apricote/hcloud-upload-image/cmd"
)

func run() error {
	// Define the directory where the docs will be generated
	dir := "docs/reference/cli"

	// Ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating docs directory: %v", err)
	}

	// Generate the docs
	if err := doc.GenMarkdownTree(cmd.RootCmd, dir); err != nil {
		return fmt.Errorf("error generating docs: %v", err)
	}

	fmt.Println("Docs generated successfully in", dir)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
