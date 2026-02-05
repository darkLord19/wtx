package main

import (
	"fmt"

	"github.com/darkLord19/wtx/internal/editor"
	"github.com/darkLord19/wtx/internal/tui"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run the setup wizard",
	Long:  "Launch the interactive setup wizard to configure wtx settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		completed, err := tui.RunSetup(cfg, edDetector)
		if err != nil {
			return fmt.Errorf("setup wizard failed: %w", err)
		}
		if !completed {
			fmt.Println("Setup cancelled.")
			return nil
		}
		// Reload editor detector with new config
		edDetector = editor.NewDetector(cfg)
		fmt.Println("Setup complete! Your settings have been saved.")
		return nil
	},
}

func init() {
	// Register the setup command in main.go init
}
