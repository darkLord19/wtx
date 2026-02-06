package main

import (
	"fmt"
	"time"

	"github.com/darkLord19/wtx/internal/metadata"
	"github.com/darkLord19/wtx/internal/validation"
	"github.com/spf13/cobra"
)

var (
	baseBranch string
)

var addCmd = &cobra.Command{
	Use:   "add <name> [branch]",
	Short: "Create a new worktree",
	Long:  "Create a new worktree with the specified name and optionally a branch name",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		branch := name
		if len(args) == 2 {
			branch = args[1]
		}

		validator := validation.NewWorktreeValidator()
		if err := validator.ValidateName(name); err != nil {
			return err
		}
		if err := validator.ValidateBranchName(branch); err != nil {
			return err
		}

		fmt.Printf("Creating worktree '%s' for branch '%s'...\n", name, branch)

		// Create worktree
		path, err := gitMgr.Add(name, branch, baseBranch)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Created worktree: %s\n", name)
		fmt.Printf("  Path: %s\n", path)
		fmt.Printf("  Branch: %s\n", branch)

		// Save metadata
		meta := &metadata.WorktreeMetadata{
			Name:       name,
			Path:       path,
			Branch:     branch,
			CreatedAt:  time.Now(),
			LastOpened: time.Now(),
		}
		metaStore.Add(meta)
		if err := metaStore.Save(); err != nil {
			fmt.Printf("Warning: failed to save metadata: %v\n", err)
		}

		// Ask if user wants to open now
		fmt.Print("\nOpen in editor now? [Y/n]: ")
		var response string
		// We explicitly ignore the error from Scanln here because if the user just hits Enter,
		// Scanln returns an error (unexpected newline) but we want to treat that as "empty input"
		// which results in the default behavior (opening the editor).
		_, _ = fmt.Scanln(&response)

		if response == "" || response == "y" || response == "Y" {
			ed, err := edDetector.GetPreferred()
			if err != nil {
				return err
			}

			fmt.Printf("Opening in %s...\n", ed.Name())
			if err := ed.Open(path, cfg.ReuseWindow); err != nil {
				return fmt.Errorf("failed to open editor: %w", err)
			}
		}

		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&baseBranch, "from", "f", "main", "Base branch to create from")
}
