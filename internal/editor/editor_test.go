package editor

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/darkLord19/wtx/internal/config"
)

// TestHelperProcess isn't a real test. It's used as a helper process to mock exec.Command
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Print arguments to stdout so we can verify them
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stdout, "no args\n")
		os.Exit(0)
	}
	fmt.Fprintf(os.Stdout, "%v\n", args)
	os.Exit(0)
}

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func mockLookPathSuccess(file string) (string, error) {
	return "/bin/" + file, nil
}

func mockLookPathFail(file string) (string, error) {
	return "", fmt.Errorf("executable file not found in $PATH")
}

func TestInstalled(t *testing.T) {
	// Restore original functions after test
	defer func() {
		execLookPath = exec.LookPath
	}()

	// Test all types
	allTypes := []EditorType{VSCode, Cursor, VSCodium, Neovim, Vim, Terminal}

	for _, edType := range allTypes {
		t.Run(string(edType), func(t *testing.T) {
			execLookPath = mockLookPathSuccess

			editor, err := New(edType)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			if !editor.Installed() {
				t.Error("Installed() = false, want true")
			}

			if editor.Name() == "" {
				t.Error("Name() is empty")
			}
		})
	}

	t.Run("VSCode not installed", func(t *testing.T) {
		execLookPath = mockLookPathFail
		editor, _ := New(VSCode)
		if editor.Installed() {
			t.Error("Installed() = true, want false")
		}
	})
}

func TestOpen(t *testing.T) {
	defer func() {
		execCommand = exec.Command
	}()
	execCommand = mockExecCommand

	allTypes := []EditorType{VSCode, Cursor, VSCodium, Neovim, Vim, Terminal}

	// Use a temp dir that exists for Terminal Chdir test
	tmpDir := t.TempDir()

	for _, edType := range allTypes {
		t.Run(string(edType), func(t *testing.T) {
			editor, _ := New(edType)

			// Open
			if err := editor.Open(tmpDir, false); err != nil {
				t.Errorf("Open() failed for %s: %v", edType, err)
			}

			// Open with reuse (if applicable)
			if err := editor.Open(tmpDir, true); err != nil {
				t.Errorf("Open(reuse=true) failed for %s: %v", edType, err)
			}
		})
	}
}

func TestDetector(t *testing.T) {
	defer func() {
		execLookPath = exec.LookPath
	}()

	// Mock config
	cfg := config.Default()
	d := NewDetector(cfg)

	t.Run("Configured Editor", func(t *testing.T) {
		cfg.Editor = "vscode"
		execLookPath = mockLookPathSuccess

		ed, err := d.GetPreferred()
		if err != nil {
			t.Fatalf("GetPreferred() error = %v", err)
		}
		if ed.Name() != "Visual Studio Code" {
			t.Errorf("GetPreferred() = %v, want Visual Studio Code", ed.Name())
		}
	})

	t.Run("Fallback to Terminal", func(t *testing.T) {
		cfg.Editor = ""
		execLookPath = mockLookPathFail
		// Unset EDITOR env var just in case
		os.Unsetenv("EDITOR")

		ed, err := d.GetPreferred()
		if err != nil {
			t.Fatalf("GetPreferred() error = %v", err)
		}
		if ed.Name() != "Terminal" {
			t.Errorf("GetPreferred() = %v, want Terminal", ed.Name())
		}
	})

	t.Run("Env VAR EDITOR", func(t *testing.T) {
		cfg.Editor = ""
		os.Setenv("EDITOR", "vim")
		execLookPath = mockLookPathSuccess

		ed, err := d.GetPreferred()
		if err != nil {
			t.Fatalf("GetPreferred() error = %v", err)
		}
		if ed.Name() != "Vim" {
			t.Errorf("GetPreferred() = %v, want Vim", ed.Name())
		}
		os.Unsetenv("EDITOR")
	})

	t.Run("DetectAll", func(t *testing.T) {
		execLookPath = mockLookPathSuccess
		editors := d.DetectAll()
		if len(editors) == 0 {
			t.Error("DetectAll() found no editors")
		}
		// Expect at least VSCode, Cursor, VSCodium, Neovim, Vim
		// Since we mocked success, all should be there
		if len(editors) < 5 {
			t.Errorf("DetectAll() found %d editors, want >= 5", len(editors))
		}
	})
}

func TestNew_Unknown(t *testing.T) {
	_, err := New("unknown")
	if err == nil {
		t.Error("New(unknown) should fail")
	}
}
