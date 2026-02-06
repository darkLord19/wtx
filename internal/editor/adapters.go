package editor

import (
	"os"
	"os/exec"
)

var (
	execLookPath = exec.LookPath
	execCommand  = exec.Command
)

// VSCodeEditor implements the Editor interface for VS Code
type VSCodeEditor struct{}

func (e *VSCodeEditor) Name() string {
	return "Visual Studio Code"
}

func (e *VSCodeEditor) Installed() bool {
	_, err := execLookPath("code")
	return err == nil
}

func (e *VSCodeEditor) Open(path string, reuseWindow bool) error {
	args := []string{path}
	if reuseWindow {
		args = append([]string{"-r"}, args...)
	}

	cmd := execCommand("code", args...)
	return cmd.Start()
}

// CursorEditor implements the Editor interface for Cursor
type CursorEditor struct{}

func (e *CursorEditor) Name() string {
	return "Cursor"
}

func (e *CursorEditor) Installed() bool {
	_, err := execLookPath("cursor")
	return err == nil
}

func (e *CursorEditor) Open(path string, reuseWindow bool) error {
	args := []string{path}
	if reuseWindow {
		args = append([]string{"-r"}, args...)
	}

	cmd := execCommand("cursor", args...)
	return cmd.Start()
}

// VSCodiumEditor implements the Editor interface for VSCodium
type VSCodiumEditor struct{}

func (e *VSCodiumEditor) Name() string {
	return "VSCodium"
}

func (e *VSCodiumEditor) Installed() bool {
	_, err := execLookPath("codium")
	return err == nil
}

func (e *VSCodiumEditor) Open(path string, reuseWindow bool) error {
	args := []string{path}
	if reuseWindow {
		args = append([]string{"-r"}, args...)
	}

	cmd := execCommand("codium", args...)
	return cmd.Start()
}

// NeovimEditor implements the Editor interface for Neovim
type NeovimEditor struct{}

func (e *NeovimEditor) Name() string {
	return "Neovim"
}

func (e *NeovimEditor) Installed() bool {
	_, err := execLookPath("nvim")
	return err == nil
}

func (e *NeovimEditor) Open(path string, reuseWindow bool) error {
	cmd := execCommand("nvim", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// VimEditor implements the Editor interface for Vim
type VimEditor struct{}

func (e *VimEditor) Name() string {
	return "Vim"
}

func (e *VimEditor) Installed() bool {
	_, err := execLookPath("vim")
	return err == nil
}

func (e *VimEditor) Open(path string, reuseWindow bool) error {
	cmd := execCommand("vim", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// TerminalEditor is a fallback that just prints the path
type TerminalEditor struct{}

func (e *TerminalEditor) Name() string {
	return "Terminal"
}

func (e *TerminalEditor) Installed() bool {
	return true
}

func (e *TerminalEditor) Open(path string, reuseWindow bool) error {
	// Just print the path and change directory
	println("\nWorktree path:", path)
	println("cd", path)
	return os.Chdir(path)
}
