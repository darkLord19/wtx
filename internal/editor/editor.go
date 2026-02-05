package editor

import "fmt"

// Editor represents a code editor
type Editor interface {
	Name() string
	Installed() bool
	Open(path string, reuseWindow bool) error
}

// EditorType represents different editor types
type EditorType string

const (
	VSCode   EditorType = "vscode"
	Cursor   EditorType = "cursor"
	VSCodium EditorType = "vscodium"
	Neovim   EditorType = "neovim"
	Vim      EditorType = "vim"
	Terminal EditorType = "terminal"
)

// New creates a new editor instance
func New(editorType EditorType) (Editor, error) {
	switch editorType {
	case VSCode:
		return &VSCodeEditor{}, nil
	case Cursor:
		return &CursorEditor{}, nil
	case VSCodium:
		return &VSCodiumEditor{}, nil
	case Neovim:
		return &NeovimEditor{}, nil
	case Vim:
		return &VimEditor{}, nil
	case Terminal:
		return &TerminalEditor{}, nil
	default:
		return nil, fmt.Errorf("unknown editor type: %s", editorType)
	}
}
