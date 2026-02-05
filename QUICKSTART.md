# WTX - Quick Setup Guide

## 🚀 Getting Started in 5 Minutes

### Prerequisites
- Go 1.21 or later installed
- Git 2.x installed
- Terminal access

### Step 1: Extract the Project
```bash
tar -xzf wtx-complete.tar.gz
cd wtx
```

### Step 2: Install Dependencies
```bash
go mod download
```

### Step 3: Build
```bash
make build
# or
go build -o bin/wtx ./cmd/wtx
```

### Step 4: Test (Optional but Recommended)
```bash
cd /tmp
git init test-repo
cd test-repo
git commit --allow-empty -m "Initial commit"

# Now try wtx
/path/to/wtx/bin/wtx --help
```

### Step 5: Install System-Wide (Optional)
```bash
# From wtx directory
make install
# or
go install ./cmd/wtx

# Now you can use 'wtx' from anywhere
wtx --help
```

## 🎯 First Run

### In a Git Repository
```bash
cd your-git-repo

# Interactive mode
wtx

# Or create your first worktree
wtx add feature-test

# List all worktrees
wtx list
```

## 📝 Project Statistics

- **23 Go files**
- **~1,844 lines of code**
- **7 CLI commands**
- **6 internal packages**
- **Full test coverage for metadata package**

## 🏗️ Project Structure

```
wtx/
├── cmd/wtx/              # Main application (7 commands)
├── internal/
│   ├── git/              # Git operations
│   ├── editor/           # Editor integration (6 adapters)
│   ├── metadata/         # Data persistence
│   ├── config/           # Configuration
│   ├── tui/              # Terminal UI
│   └── ports/            # Port detection
├── test/                 # Tests
├── docs/                 # Documentation
├── Makefile             # Build automation
├── .goreleaser.yaml     # Release config
└── .github/workflows/   # CI/CD
```

## 🔧 Development Commands

```bash
# Build
make build

# Run tests
make test

# Format code
make fmt

# Build and run
make dev

# Clean build artifacts
make clean

# Build for all platforms
make build-all

# Show help
make help
```

## 📦 What's Included

### Core Features ✅
- [x] Interactive TUI switcher with fuzzy search
- [x] Multi-editor support (VS Code, Cursor, Neovim, Vim)
- [x] Safe worktree creation and removal
- [x] Git status indicators (clean/dirty, ahead/behind)
- [x] Metadata tracking (creation time, last opened)
- [x] Configuration system
- [x] Stale worktree cleanup

### Commands
1. `wtx` - Interactive TUI (default)
2. `wtx list` - List all worktrees
3. `wtx add <name>` - Create new worktree
4. `wtx open <name>` - Open specific worktree
5. `wtx rm <name>` - Remove worktree (with safety)
6. `wtx status <name>` - Detailed status
7. `wtx prune` - Clean stale worktrees
8. `wtx config` - View configuration

### Files & Documentation
- ✅ Complete implementation (23 Go files)
- ✅ README.md with full documentation
- ✅ CONTRIBUTING.md with guidelines
- ✅ LICENSE (MIT)
- ✅ Makefile for automation
- ✅ Tests (metadata package fully tested)
- ✅ CI/CD workflows (GitHub Actions)
- ✅ GoReleaser config
- ✅ Project overview document

## 🎨 Customization

### Config File Location
`~/.config/wtx/config.json`

### Default Config
```json
{
  "editor": "",
  "reuse_window": true,
  "worktree_dir": "../worktrees",
  "auto_start_dev": false,
  "custom_commands": {}
}
```

### To Change Editor
Edit `~/.config/wtx/config.json`:
```json
{
  "editor": "cursor"
}
```

Or set `$EDITOR` environment variable:
```bash
export EDITOR=nvim
```

## 🐛 Troubleshooting

### "git is not installed"
Install Git 2.x for your platform.

### "not a git repository"
Run wtx from within a git repository:
```bash
cd your-git-repo
wtx
```

### "no worktrees found"
This is normal for a fresh repo. Create one:
```bash
wtx add my-first-worktree
```

### Editor doesn't open
Check if editor command is in PATH:
```bash
which code    # VS Code
which cursor  # Cursor
which nvim    # Neovim
```

## 📚 Next Steps

1. **Read the README** - Full documentation
2. **Try the TUI** - Run `wtx` in a repo
3. **Create worktrees** - `wtx add feature-name`
4. **Customize config** - Edit `~/.config/wtx/config.json`
5. **Report issues** - Open GitHub issues
6. **Contribute** - See CONTRIBUTING.md

## 🚢 Deployment

### Build Release Binaries
```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Build release (creates dist/ directory)
goreleaser release --snapshot --clean
```

### GitHub Release
```bash
# Tag version
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GitHub Actions will automatically build and release
```

## 💡 Tips

- Use `wtx` (no args) for fastest switching
- Set up shell alias: `alias w=wtx`
- Keep worktrees clean with `wtx prune`
- Check status with `wtx status <name>`
- Use Ctrl+C or q to quit TUI

## 🎉 You're Ready!

The complete wtx project is now set up and ready to use. Enjoy managing your Git worktrees with ease!

---

**Questions?** Check the README.md or open an issue on GitHub.
