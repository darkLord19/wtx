package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <name>",
	Short: "Show detailed status of a worktree",
	Long:  "Display detailed git status and metadata for a specific worktree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Find the worktree
		worktrees, err := gitMgr.List()
		if err != nil {
			return err
		}

		var target *struct {
			Path   string
			Branch string
			IsMain bool
		}

		for _, wt := range worktrees {
			if wt.Name == name {
				target = &struct {
					Path   string
					Branch string
					IsMain bool
				}{
					Path:   wt.Path,
					Branch: wt.Branch,
					IsMain: wt.IsMain,
				}
				break
			}
		}

		if target == nil {
			return fmt.Errorf("worktree '%s' not found", name)
		}

		// Get status
		status, err := gitMgr.GetStatus(target.Path)
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}

		// Get metadata
		meta, _ := metaStore.Get(name)

		// Display
		fmt.Printf("\nWorktree: %s\n", name)
		fmt.Println("─────────────────────────────────────")
		fmt.Printf("Path:     %s\n", target.Path)
		fmt.Printf("Branch:   %s\n", target.Branch)

		if target.IsMain {
			fmt.Println("Type:     Main worktree ⭐")
		} else {
			fmt.Println("Type:     Linked worktree")
		}

		fmt.Println()
		fmt.Println("Git Status:")
		if status.Clean {
			fmt.Println("  ● Working tree clean")
		} else {
			fmt.Println("  ✗ Uncommitted changes")
		}

		if status.Ahead > 0 {
			fmt.Printf("  ↑ %d commit(s) ahead of upstream\n", status.Ahead)
		}
		if status.Behind > 0 {
			fmt.Printf("  ↓ %d commit(s) behind upstream\n", status.Behind)
		}

		if meta != nil {
			fmt.Println()
			fmt.Println("Metadata:")
			fmt.Printf("  Created:     %s\n", meta.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Printf("  Last opened: %s\n", meta.LastOpened.Format("2006-01-02 15:04"))

			if meta.DevCommand != "" {
				fmt.Printf("  Dev command: %s\n", meta.DevCommand)
			}

			if len(meta.Ports) > 0 {
				fmt.Printf("  Ports:       %v\n", meta.Ports)
			}
		}

		fmt.Println()
		return nil
	},
}
