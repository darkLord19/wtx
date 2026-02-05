package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <name>",
	Short: "Open a specific worktree",
	Long:  "Open a specific worktree in the configured editor",
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

		// Open in editor
		ed, err := edDetector.GetPreferred()
		if err != nil {
			return err
		}

		fmt.Printf("Opening %s in %s...\n", name, ed.Name())

		if err := ed.Open(targetPath, cfg.ReuseWindow); err != nil {
			return fmt.Errorf("failed to open editor: %w", err)
		}

		// Update metadata
		metaStore.Touch(name)
		return metaStore.Save()
	},
}
