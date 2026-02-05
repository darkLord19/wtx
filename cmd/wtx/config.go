package main

import (
	"fmt"

	"github.com/darkLord19/wtx/internal/tui"
	"github.com/spf13/cobra"
)

var (
	configTUI bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or edit configuration",
	Long:  "Display current configuration or launch TUI to edit settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Launch TUI settings editor
		if configTUI {
			return tui.RunSettings(cfg, edDetector)
		}

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

			fmt.Println("\nTip: Use 'wtx config --tui' to edit settings interactively")

			return nil
		}

		// TODO: Implement config set
		fmt.Println("Setting config values not yet implemented")
		fmt.Println("Use 'wtx config --tui' to edit settings interactively")

		return nil
	},
}

func init() {
	configCmd.Flags().BoolVarP(&configTUI, "tui", "t", false, "Launch TUI settings editor")
}
