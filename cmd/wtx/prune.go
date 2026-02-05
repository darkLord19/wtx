package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	staleDays int
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Clean up stale worktrees",
	Long:  "Remove worktrees that haven't been opened in a specified number of days (default: 30)",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get stale worktrees
		staleNames := metaStore.GetStale(staleDays)

		if len(staleNames) == 0 {
			fmt.Println("No stale worktrees found")
			return nil
		}

		// Check which ones are clean
		cleanStale := []string{}
		worktrees, err := gitMgr.List()
		if err != nil {
			return err
		}

		for _, name := range staleNames {
			for _, wt := range worktrees {
				if wt.Name == name {
					clean, _ := gitMgr.IsClean(wt.Path)
					if clean {
						cleanStale = append(cleanStale, name)
					}
					break
				}
			}
		}

		if len(cleanStale) == 0 {
			fmt.Printf("No clean stale worktrees found (>%d days old)\n", staleDays)
			return nil
		}

		fmt.Printf("Stale worktrees (clean, >%d days old):\n\n", staleDays)
		for _, name := range cleanStale {
			meta, _ := metaStore.Get(name)
			if meta != nil {
				fmt.Printf("  • %s (last opened: %s)\n", name, meta.LastOpened.Format("2006-01-02"))
			} else {
				fmt.Printf("  • %s\n", name)
			}
		}

		fmt.Print("\nDelete all? [y/N]: ")
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			fmt.Println("Cancelled")
			return nil
		}

		if response != "y" && response != "Y" {
			fmt.Println("Cancelled")
			return nil
		}

		// Delete them
		removed := 0
		for _, name := range cleanStale {
			if err := gitMgr.Remove(name, false); err != nil {
				fmt.Printf("⚠  Failed to remove %s: %v\n", name, err)
				continue
			}
			metaStore.Remove(name)
			removed++
			fmt.Printf("✓ Removed %s\n", name)
		}

		if err := metaStore.Save(); err != nil {
			fmt.Printf("Warning: failed to update metadata: %v\n", err)
		}

		fmt.Printf("\n✓ Removed %d worktree(s)\n", removed)
		return nil
	},
}

func init() {
	pruneCmd.Flags().IntVarP(&staleDays, "days", "d", 30, "Number of days to consider a worktree stale")
}
