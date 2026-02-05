package main

import (
	"github.com/darkLord19/wtx/internal/tui"
	"github.com/spf13/cobra"
)

var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Launch worktree management TUI",
	Long:  "Interactive TUI for creating, deleting, and pruning worktrees",
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.RunWorktreeManager(gitMgr, metaStore)
	},
}

func init() {
	// Register the manage command in main.go init
}
