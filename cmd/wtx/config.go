package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or edit configuration",
	Long:  "Display current configuration or set specific values",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Display current config
			fmt.Println("Current configuration:")
			fmt.Println("──────────────────────────────────")
			fmt.Printf("Editor:         %s", cfg.Editor)
			if cfg.Editor == "" {
				ed, _ := edDetector.GetPreferred()
				if ed != nil {
					fmt.Printf(" (auto-detected: %s)", ed.Name())
				}
			}
			fmt.Println()
			fmt.Printf("Reuse window:   %v\n", cfg.ReuseWindow)
			fmt.Printf("Worktree dir:   %s\n", cfg.WorktreeDir)
			fmt.Printf("Auto start dev: %v\n", cfg.AutoStartDev)

			if len(cfg.CustomCommands) > 0 {
				fmt.Println("\nCustom commands:")
				for k, v := range cfg.CustomCommands {
					fmt.Printf("  %s: %s\n", k, v)
				}
			}

			fmt.Println("\nDetected editors:")
			editors := edDetector.DetectAll()
			for _, ed := range editors {
				fmt.Printf("  • %s\n", ed.Name())
			}

			return nil
		}

		// TODO: Implement config set
		fmt.Println("Setting config values not yet implemented")
		fmt.Println("Edit ~/.config/wtx/config.json directly")

		return nil
	},
}
