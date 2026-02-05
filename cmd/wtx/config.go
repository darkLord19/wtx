package main

import (
	"fmt"
	"strconv"

	"github.com/darkLord19/wtx/internal/tui"
	"github.com/spf13/cobra"
)

var (
	configTUI bool
)

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "View or edit configuration",
	Long:  "Display current configuration, set values, or launch TUI to edit settings",
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

		key := args[0]

		switch key {
		case "editor":
			if len(args) == 4 && args[1] == "custom_command" {
				// wtx config editor custom_command <name> <cmd>
				name := args[2]
				command := args[3]

				if cfg.CustomCommands == nil {
					cfg.CustomCommands = make(map[string]string)
				}
				cfg.CustomCommands[name] = command
				if err := cfg.Save(); err != nil {
					return fmt.Errorf("failed to save config: %w", err)
				}
				fmt.Printf("Set custom command '%s' to '%s'\n", name, command)
				return nil
			} else if len(args) == 2 {
				// wtx config editor <value>
				val := args[1]
				// Validate/Warn
				validEditors := map[string]bool{
					"vscode": true, "cursor": true, "vscodium": true,
					"neovim": true, "vim": true, "terminal": true,
				}
				if !validEditors[val] {
					fmt.Printf("Warning: '%s' is not a known editor type. Known types: vscode, cursor, vscodium, neovim, vim, terminal\n", val)
				}

				cfg.Editor = val
				if err := cfg.Save(); err != nil {
					return fmt.Errorf("failed to save config: %w", err)
				}
				fmt.Printf("Set editor to '%s'\n", val)
				return nil
			} else {
				return fmt.Errorf("invalid arguments for editor config.\nUsage: \n  wtx config editor <value>\n  wtx config editor custom_command <name> <command>")
			}

		case "reuse_window":
			if len(args) != 2 {
				return fmt.Errorf("usage: wtx config reuse_window <true|false>")
			}
			val, err := strconv.ParseBool(args[1])
			if err != nil {
				return fmt.Errorf("invalid boolean value: %s", args[1])
			}
			cfg.ReuseWindow = val
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Printf("Set reuse_window to %v\n", val)
			return nil

		case "worktree_dir":
			if len(args) != 2 {
				return fmt.Errorf("usage: wtx config worktree_dir <path>")
			}
			val := args[1]
			if val == "" {
				return fmt.Errorf("worktree_dir cannot be empty")
			}
			cfg.WorktreeDir = val
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Printf("Set worktree_dir to '%s'\n", val)
			return nil

		case "auto_start_dev":
			if len(args) != 2 {
				return fmt.Errorf("usage: wtx config auto_start_dev <true|false>")
			}
			val, err := strconv.ParseBool(args[1])
			if err != nil {
				return fmt.Errorf("invalid boolean value: %s", args[1])
			}
			cfg.AutoStartDev = val
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Printf("Set auto_start_dev to %v\n", val)
			return nil

		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}
	},
}

func init() {
	configCmd.Flags().BoolVarP(&configTUI, "tui", "t", false, "Launch TUI settings editor")
}
