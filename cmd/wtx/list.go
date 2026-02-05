package main

import (
	"fmt"
	"sync"

	"github.com/darkLord19/wtx/internal/git"
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

		// Prepare for concurrent status checks
		type statusResult struct {
			index  int
			status *git.Status
			err    error
		}

		numWorkers := 10
		if len(worktrees) < numWorkers {
			numWorkers = len(worktrees)
		}

		jobs := make(chan int, len(worktrees))
		results := make(chan statusResult, len(worktrees))
		var wg sync.WaitGroup

		// Start workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := range jobs {
					wt := worktrees[i]
					status, err := gitMgr.GetStatus(wt.Path)
					results <- statusResult{index: i, status: status, err: err}
				}
			}()
		}

		// Send jobs
		for i := range worktrees {
			jobs <- i
		}
		close(jobs)

		// Close results channel when all workers are done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Collect results
		statuses := make([]*git.Status, len(worktrees))
		statusErrors := make([]error, len(worktrees))
		for res := range results {
			statuses[res.index] = res.status
			statusErrors[res.index] = res.err
		}

		fmt.Printf("%-20s %-30s %-10s %s\n", "NAME", "BRANCH", "STATUS", "PATH")
		fmt.Println("────────────────────────────────────────────────────────────────────────────")

		for i, wt := range worktrees {
			status := statuses[i]
			err := statusErrors[i]

			statusStr := "●"
			statusText := "clean"

			if err != nil {
				statusStr = "?"
				statusText = "error"
			} else if status != nil && !status.Clean {
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
