# wtx - Git Worktree Workspace Manager

**"Cmd+Tab for your Git branches"**

`wtx` makes Git worktrees feel like instant "workspace tabs" across editors with zero friction. Switch between worktrees in under 2 seconds, open them in your favorite editor, and manage multiple development environments in parallel.

## Features

- 🚀 **Zero-friction switching** - Interactive TUI to switch worktrees in <2 seconds
- 🎯 **Editor-native feel** - Opens in VS Code, Cursor, Neovim, or your preferred editor
- 🛡️ **Safe by default** - Never lose uncommitted work with smart safety checks
- ⚡ **Parallel-first** - Built for running multiple dev environments simultaneously
- ⌨️ **Keyboard-driven** - Everything accessible without touching the mouse
- 📊 **Smart status** - See clean/dirty status, ahead/behind commits at a glance

## Quick Start

### Installation

```bash
# Clone and install
git clone https://github.com/darkLord19/wtx.git
cd wtx
make install

# Or build locally
make build
./bin/wtx
```

### Basic Usage

```bash
# Interactive worktree switcher (main command)
wtx

# Full TUI with tabs (worktrees, manage, settings)
wtx --tui

# Create a new worktree
wtx add feature-auth

# List all worktrees
wtx list

# Open specific worktree
wtx open feature-auth

# Remove a worktree (with safety checks)
wtx rm feature-auth

# Show detailed status
wtx status feature-auth

# Clean up stale worktrees
wtx prune

# View/edit configuration
wtx config

# Interactive settings editor
wtx config --tui

# Worktree management TUI
wtx manage

# Run setup wizard
wtx setup
```

## How It Works

### The Golden Path

1. Run `wtx` in any git repository
2. Fuzzy search for the worktree you want
3. Press Enter
4. Opens in your editor with window reuse
5. Start working immediately

**Target time: < 2 seconds** ⏱️

### Creating Worktrees

```bash
# Create from current branch (main)
wtx add feature-payments

# Create from specific branch
wtx add hotfix-bug --from develop

# Creates worktree at: ../worktrees/feature-payments
```

### Safe Removal

```bash
wtx rm old-feature
```

If the worktree has uncommitted changes:
```
⚠  Worktree 'old-feature' has uncommitted changes

Options:
  c - Cancel
  f - Force delete (lose changes)

Your choice [c/f]:
```

## Editor Support

wtx automatically detects and supports:

- ✅ **VS Code** (`code -r`)
- ✅ **Cursor** (`cursor -r`)
- ✅ **VSCodium** (`codium -r`)
- ✅ **Neovim** (`nvim`)
- ✅ **Vim** (`vim`)
- ✅ **Terminal** (fallback)

### Editor Selection Priority

1. User config (`~/.config/wtx/config.json`)
2. `$EDITOR` environment variable
3. Auto-detect installed editors
4. Terminal fallback

## First-Run Setup

When you run `wtx` for the first time, an interactive setup wizard will guide you through:

1. **Editor Selection** - Choose your preferred editor or enter a custom command
2. **Worktree Directory** - Where to create new worktrees
3. **Window Reuse** - Whether to reuse existing editor windows

You can re-run the setup wizard anytime with:
```bash
wtx setup
```

## Configuration

Config file: `~/.config/wtx/config.json`

```json
{
  "editor": "cursor",
  "reuse_window": true,
  "worktree_dir": "../worktrees",
  "auto_start_dev": false,
  "custom_commands": {}
}
```

### Options

- **editor** - Override editor selection (`vscode`, `cursor`, `neovim`, etc.)
- **reuse_window** - Reuse existing editor window (default: true)
- **worktree_dir** - Where to create worktrees (default: `../worktrees`)
- **auto_start_dev** - Auto-start dev servers (future feature)
- **custom_commands** - Per-worktree custom commands

## TUI Interface

### Quick Selector (default)

```
┌──────────────────────────────────────┐
│ Workspace Manager (your-repo)       │
├──────────────────────────────────────┤
│ main            ● clean              │
│ feature-auth    ✗ dirty  ↑2         │
│ bugfix-otp      ● clean              │
│ experiment      ● clean              │
├──────────────────────────────────────┤
│ Press enter to open • q/esc to quit │
└──────────────────────────────────────┘
```

### Full TUI Manager

Launch with `wtx --tui` or `wtx -t` for the full TUI experience with tabs:

```
┌─────────────────────────────────────────────────────┐
│ [1] Worktrees   [2] Manage   [3] Settings          │
├─────────────────────────────────────────────────────┤
│                                                     │
│  Worktrees Tab: Select and open worktrees          │
│  Manage Tab: Create, delete, prune worktrees       │
│  Settings Tab: Configure wtx settings              │
│                                                     │
└─────────────────────────────────────────────────────┘
```

**Keyboard shortcuts:**
- `1`, `2`, `3` - Switch between tabs
- In Manage tab: `c` create, `d` delete, `p` prune, `r` refresh
- In Settings tab: `↑/↓` navigate, `enter` edit, `s` save

### Standalone Commands

```bash
# Launch worktree management TUI
wtx manage

# Launch settings TUI
wtx config --tui
```

### Status Indicators

| Symbol | Meaning               |
|--------|-----------------------|
| ●      | Clean working tree    |
| ✗      | Uncommitted changes   |
| ↑N     | N commits ahead       |
| ↓N     | N commits behind      |
| ⭐     | Main worktree         |
| :3000  | Dev server on port    |

## Advanced Usage

### Cleanup Stale Worktrees

```bash
# Find worktrees not opened in 30 days
wtx prune

# Custom time period
wtx prune --days 60
```

### Detailed Status

```bash
wtx status feature-auth
```

Output:
```
Worktree: feature-auth
─────────────────────────────────────
Path:     /path/to/worktrees/feature-auth
Branch:   feature-auth
Type:     Linked worktree

Git Status:
  ● Working tree clean
  ↑ 2 commit(s) ahead of upstream

Metadata:
  Created:     2024-01-15 10:30
  Last opened: 2024-01-20 14:22
```

## Why Git Worktrees?

Git worktrees let you have multiple branches checked out simultaneously. This is perfect for:

- 🔄 Switching contexts without stashing
- 🐛 Quick bug fixes while working on features
- 👀 Reviewing PRs alongside your work
- 🧪 Running tests on one branch while developing on another
- 📦 Comparing implementations side-by-side

### The Problem wtx Solves

Git worktrees are powerful but clunky:

```bash
# Without wtx 😞
cd ..
git worktree add ../feature-auth feature-auth
cd ../feature-auth
code .

# With wtx 😊
wtx
# Type "feature" → Enter
# Done!
```

## Metadata & Tracking

wtx stores metadata at `.git/wtx-meta.json`:

```json
{
  "repo_path": "/path/to/repo",
  "worktrees": {
    "feature-auth": {
      "name": "feature-auth",
      "path": "/path/to/worktrees/feature-auth",
      "branch": "feature-auth",
      "created_at": "2024-01-15T10:30:00Z",
      "last_opened": "2024-01-20T14:22:00Z",
      "dev_command": "npm run dev",
      "ports": [3000]
    }
  }
}
```

This enables:
- Smart cleanup suggestions
- Usage tracking
- Future features (dev server management, etc.)

## Development

### Requirements

- Go 1.21 or later
- Git 2.x

### Building

```bash
# Download dependencies
make deps

# Build
make build

# Run tests
make test

# Install locally
make install
```

### Project Structure

```
wtx/
├── cmd/wtx/           # CLI entry point and commands
├── internal/
│   ├── git/           # Git operations (worktrees, status)
│   ├── editor/        # Editor adapters and detection
│   ├── metadata/      # Metadata storage
│   ├── config/        # Configuration management
│   ├── tui/           # Terminal UI
│   └── ports/         # Port detection
├── test/              # Integration tests
├── Makefile
└── README.md
```

## Roadmap

### v1.0 (MVP) ✅
- [x] Interactive TUI switcher
- [x] Multi-editor support with window reuse
- [x] Safe create/delete operations
- [x] Git status indicators
- [x] Metadata persistence
- [x] Configuration system

### v1.1 (Planned)
- [ ] Dev server management
- [ ] Port conflict detection and resolution
- [ ] JSON output mode for scripting
- [ ] Shell completion (bash, zsh, fish)
- [ ] Worktree templates

### v2.0 (Future)
- [ ] VS Code extension
- [ ] Raycast extension
- [ ] GitHub CLI integration
- [ ] Team workspace sharing
- [ ] Docker workspace isolation

## FAQ

**Q: How is this different from just using `git worktree`?**  
A: wtx adds a beautiful TUI, editor integration, safety checks, metadata tracking, and makes worktrees feel like instant workspace tabs. Git worktree is powerful but raw.

**Q: Will this work with my existing worktrees?**  
A: Yes! wtx detects all existing worktrees and works seamlessly with them.

**Q: Can I use this with GitHub flow / GitLab flow?**  
A: Absolutely. wtx is workflow-agnostic and enhances any Git workflow.

**Q: Does it work on Windows?**  
A: Partially. The CLI works, but editor detection may need adjustments. Windows Terminal support is coming.

**Q: Can I use it without the TUI?**  
A: Yes! All commands work non-interactively: `wtx open feature-auth`, `wtx add new-feature`, etc.

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Credits

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling

## Support

- 🐛 [Report bugs](https://github.com/darkLord19/wtx/issues)
- 💡 [Request features](https://github.com/darkLord19/wtx/issues)
- 📖 [Documentation](https://github.com/darkLord19/wtx/wiki)
- 💬 [Discussions](https://github.com/darkLord19/wtx/discussions)

---

**Made with ❤️ for developers who love Git worktrees**

⭐ Star this repo if you find it useful!
