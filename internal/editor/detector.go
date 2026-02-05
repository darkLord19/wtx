package editor

import (
	"os"

	"github.com/darkLord19/wtx/internal/config"
)

// Detector handles editor detection and selection
type Detector struct {
	config *config.Config
}

// NewDetector creates a new editor detector
func NewDetector(cfg *config.Config) *Detector {
	return &Detector{config: cfg}
}

// GetPreferred returns the preferred editor based on config and detection
func (d *Detector) GetPreferred() (Editor, error) {
	// 1. User config
	if d.config.Editor != "" {
		editor, err := New(EditorType(d.config.Editor))
		if err == nil && editor.Installed() {
			return editor, nil
		}
	}

	// 2. $EDITOR environment variable
	if envEditor := os.Getenv("EDITOR"); envEditor != "" {
		editorMap := map[string]EditorType{
			"code":   VSCode,
			"cursor": Cursor,
			"codium": VSCodium,
			"nvim":   Neovim,
			"vim":    Vim,
		}
		if edType, ok := editorMap[envEditor]; ok {
			editor, err := New(edType)
			if err == nil && editor.Installed() {
				return editor, nil
			}
		}
	}

	// 3. Auto-detect installed editors
	priority := []EditorType{Cursor, VSCode, VSCodium, Neovim, Vim}
	for _, edType := range priority {
		editor, err := New(edType)
		if err == nil && editor.Installed() {
			return editor, nil
		}
	}

	// 4. Fallback to terminal
	return &TerminalEditor{}, nil
}

// DetectAll returns all installed editors
func (d *Detector) DetectAll() []Editor {
	var editors []Editor

	allTypes := []EditorType{VSCode, Cursor, VSCodium, Neovim, Vim}
	for _, edType := range allTypes {
		editor, err := New(edType)
		if err == nil && editor.Installed() {
			editors = append(editors, editor)
		}
	}

	return editors
}
