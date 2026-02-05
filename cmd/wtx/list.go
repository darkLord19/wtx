package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees",
	Long:  "Display all worktrees with their status and branch information",
	RunE: func(cmd *cobra.Command, args []string) error {
		worktrees, err := gitMgr.List()
		if err != nil {
			return err
		}

		if len(worktrees) == 0 {
			fmt.Println("No worktrees found")
			return nil
		}

		fmt.Printf("%-20s %-30s %-10s %s\n", "NAME", "BRANCH", "STATUS", "PATH")
		fmt.Println("────────────────────────────────────────────────────────────────────────────")

		for _, wt := range worktrees {
			status, _ := gitMgr.GetStatus(wt.Path)

			statusStr := "●"
			statusText := "clean"
			if status != nil && !status.Clean {
				statusStr = "✗"
				statusText = "dirty"
			}

			mainIndicator := ""
			if wt.IsMain {
				mainIndicator = " ⭐"
			}

			fmt.Printf("%-20s %-30s %s %-8s %s%s\n",
				wt.Name,
				wt.Branch,
				statusStr,
				statusText,
				wt.Path,
				mainIndicator,
			)
		}

		return nil
	},
}
