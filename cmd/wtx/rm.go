package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	forceRemove bool
)

var rmCmd = &cobra.Command{
	Use:   "rm <name>",
	Short: "Remove a worktree",
	Long:  "Safely remove a worktree. Prompts for confirmation if worktree has uncommitted changes.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Find the worktree
		worktrees, err := gitMgr.List()
		if err != nil {
			return err
		}

		var targetPath string
		var found bool
		for _, wt := range worktrees {
			if wt.Name == name {
				targetPath = wt.Path
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("worktree '%s' not found", name)
		}

		// Check if clean
		if !forceRemove {
			clean, err := gitMgr.IsClean(targetPath)
			if err != nil {
				return fmt.Errorf("failed to check status: %w", err)
			}

			if !clean {
				fmt.Printf("⚠  Worktree '%s' has uncommitted changes\n\n", name)
				fmt.Println("Options:")
				fmt.Println("  c - Cancel")
				fmt.Println("  f - Force delete (lose changes)")
				fmt.Print("\nYour choice [c/f]: ")

				var choice string
				if _, err := fmt.Scanln(&choice); err != nil {
					return nil // Treat input error as cancel
				}

				if choice != "f" && choice != "F" {
					fmt.Println("Cancelled")
					return nil
				}
				forceRemove = true
			}
		}

		// Remove worktree
		fmt.Printf("Removing worktree '%s'...\n", name)
		if err := gitMgr.Remove(name, forceRemove); err != nil {
			return err
		}

		// Remove from metadata
		metaStore.Remove(name)
		if err := metaStore.Save(); err != nil {
			fmt.Printf("Warning: failed to update metadata: %v\n", err)
		}

		fmt.Printf("✓ Removed worktree: %s\n", name)
		return nil
	},
}

func init() {
	rmCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Force removal even with uncommitted changes")
}
