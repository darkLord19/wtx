package git

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseWorktreeList(t *testing.T) {
	// Note: The 'bare' line in the last worktree is currently ignored by the implementation
	// because it splits by space and expects 2 parts. 'bare' has only 1 part.
	// This test reflects the current behavior.
	output := `worktree /path/to/main
HEAD abc1234
branch refs/heads/main

worktree /path/to/feature
HEAD def5678
branch refs/heads/feature

worktree /path/to/bare
bare
`

	worktrees, err := parseWorktreeList(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(worktrees) != 3 {
		t.Errorf("expected 3 worktrees, got %d", len(worktrees))
	}

	// Helper to check worktree fields
	check := func(i int, wt Worktree, path, name, head, branch string, isMain bool) {
		t.Helper()
		if wt.Path != path {
			t.Errorf("worktree[%d].Path = %q, want %q", i, wt.Path, path)
		}
		if wt.Name != name {
			t.Errorf("worktree[%d].Name = %q, want %q", i, wt.Name, name)
		}
		if wt.Head != head {
			t.Errorf("worktree[%d].Head = %q, want %q", i, wt.Head, head)
		}
		if wt.Branch != branch {
			t.Errorf("worktree[%d].Branch = %q, want %q", i, wt.Branch, branch)
		}
		if wt.IsMain != isMain {
			t.Errorf("worktree[%d].IsMain = %v, want %v", i, wt.IsMain, isMain)
		}
	}

	// First worktree is marked as main because it's the first one and no other was marked main.
	check(0, worktrees[0], "/path/to/main", "main", "abc1234", "main", true)
	check(1, worktrees[1], "/path/to/feature", "feature", "def5678", "feature", false)
	// Third worktree: 'bare' is ignored, so IsMain is false.
	check(2, worktrees[2], "/path/to/bare", "bare", "", "", false)
}

func generateWorktreeOutput(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("worktree /path/to/worktree-%d\n", i))
		sb.WriteString(fmt.Sprintf("HEAD head-%d\n", i))
		sb.WriteString(fmt.Sprintf("branch refs/heads/branch-%d\n", i))
		sb.WriteString("\n")
	}
	return sb.String()
}

func BenchmarkParseWorktreeList(b *testing.B) {
	// Use a large input for benchmarking
	input := generateWorktreeOutput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parseWorktreeList(input)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
