# WTX Project Overview

## What is WTX?

wtx (Worktree eXperience) is a Git worktree workspace manager that makes switching between branches feel like instant "workspace tabs" in your editor. It provides a beautiful TUI, editor integration, and safety features to make Git worktrees practical for everyday development.

## Project Status

✅ **MVP Complete** - All core features implemented

## Architecture

### Technology Stack

- **Language**: Go 1.21+
- **TUI Framework**: Bubble Tea (charmbracelet)
- **CLI Framework**: Cobra
- **Configuration**: Viper
- **Styling**: Lipgloss

### Project Structure

```
wtx/
├── cmd/wtx/              # CLI application
│   ├── main.go          # Entry point, root command
│   ├── list.go          # List command
│   ├── add.go           # Add command
│   ├── rm.go            # Remove command
│   ├── open.go          # Open command
│   ├── status.go        # Status command
│   ├── prune.go         # Prune command
│   └── config.go        # Config command
│
├── internal/
│   ├── git/             # Git operations
│   │   ├── repo.go      # Repository detection
│   │   ├── worktree.go  # Worktree management
│   │   └── status.go    # Status checking
│   │
│   ├── editor/          # Editor integration
│   │   ├── editor.go    # Interface definition
│   │   ├── adapters.go  # VS Code, Cursor, Neovim, etc.
│   │   └── detector.go  # Auto-detection logic
│   │
│   ├── metadata/        # Data persistence
│   │   ├── models.go    # Data structures
│   │   └── store.go     # Save/load logic
│   │
│   ├── config/          # Configuration
│   │   ├── config.go    # Config model
│   │   └── loader.go    # Load/save config
│   │
│   ├── tui/             # Terminal UI
│   │   ├── selector.go  # Main TUI
│   │   ├── models.go    # Data models
│   │   └── styles.go    # Styling
│   │
│   └── ports/           # Port detection
│       └── detector.go  # Check port availability
│
├── test/
│   ├── integration/     # Integration tests
│   └── testdata/        # Test fixtures
│
├── docs/                # Documentation
├── Makefile            # Build automation
├── go.mod              # Go dependencies
├── .goreleaser.yaml    # Release configuration
└── README.md           # User documentation
```

## Core Features

### 1. Interactive TUI Switcher
- Fuzzy search through worktrees
- Real-time git status display
- Keyboard-driven navigation
- Fast performance (<150ms startup)

### 1.1 Full TUI Manager (wtx --tui)
- Tab-based interface with three views:
  - **Worktrees Tab**: Select and open worktrees
  - **Manage Tab**: Create, delete, and prune worktrees
  - **Settings Tab**: Configure wtx settings interactively
- Keyboard shortcuts: 1/2/3 to switch tabs
- All worktree operations available without leaving TUI

### 1.2 TUI Settings Editor (wtx config --tui)
- Interactive settings configuration
- Option cycling for selection-based settings
- Text input for custom values
- Immediate save functionality

### 1.3 Worktree Manager TUI (wtx manage)
- Create new worktrees with name, branch, and base branch
- Delete worktrees with safety checks
- Prune stale worktrees with selection interface
- Refresh worktree list

### 2. Multi-Editor Support
- VS Code, Cursor, VSCodium
- Neovim, Vim
- Terminal fallback
- Window reuse support
- Auto-detection with priority

### 3. Safety-First Design
- Never delete dirty worktrees without confirmation
- Multiple confirmation levels
- Clear error messages
- Graceful error handling

### 4. Metadata Tracking
- Creation timestamp
- Last opened timestamp
- Custom dev commands
- Port assignments
- JSON storage in .git/

### 5. Git Integration
- Clean/dirty status
- Ahead/behind tracking
- Branch information
- Automatic git operations

### 6. Smart Cleanup
- Find stale worktrees
- Configurable age threshold
- Batch removal
- Safety checks

## Key Design Decisions

### Why Go?
- Single binary distribution
- Fast compilation and execution
- Excellent CLI/TUI libraries
- Cross-platform support
- Strong standard library

### Why Bubble Tea?
- Modern, composable TUI framework
- Elm architecture (predictable state management)
- Beautiful styling with Lipgloss
- Active community

### Why Git Worktrees?
- Native Git feature (reliable)
- No additional tools required
- Well-documented
- Already handles the hard parts

### File Locations
- **Config**: `~/.config/wtx/config.json`
- **Metadata**: `.git/wtx-meta.json` (per-repo)
- **Worktrees**: `../worktrees/` (configurable)

### Safety Philosophy
- Explicit > Implicit
- Multiple confirmation for destructive actions
- Never auto-delete uncommitted work
- Clear, actionable error messages

## User Workflows

### Golden Path (< 2 seconds)
1. User runs `wtx`
2. TUI shows all worktrees with status
3. User types to filter (fuzzy search)
4. User presses Enter
5. Opens in editor with window reuse
6. User starts working

### Creation Flow
1. User runs `wtx add feature-name`
2. Tool checks if branch exists
3. Creates worktree in configured location
4. Saves metadata
5. Prompts to open immediately
6. Opens if confirmed

### Removal Flow
1. User runs `wtx rm feature-name`
2. Tool checks for uncommitted changes
3. Prompts if dirty
4. Removes only if safe or forced
5. Cleans up metadata
6. Confirms removal

## Performance Targets

| Operation           | Target  | Status |
|---------------------|---------|--------|
| TUI startup         | <150ms  | ✅     |
| List worktrees      | <100ms  | ✅     |
| Switch workspace    | <2s     | ✅     |
| Create worktree     | <5s     | ✅     |
| Status check        | <200ms  | ✅     |

## Testing Strategy

### Unit Tests
- All packages have test coverage
- Table-driven tests for variations
- Mock external dependencies
- Target: 80%+ coverage

### Integration Tests
- End-to-end workflows
- Temporary git repositories
- All commands tested
- Shell script harness

### Manual Testing
- Multiple platforms (macOS, Linux)
- Different editors
- Various git states
- Edge cases

## Build & Release

### Development
```bash
make dev      # Build and run
make test     # Run tests
make fmt      # Format code
```

### Release Process
1. Tag version: `git tag v1.0.0`
2. Push tag: `git push origin v1.0.0`
3. GitHub Actions runs GoReleaser
4. Binaries published to GitHub Releases
5. Homebrew formula updated (future)

### Distribution
- GitHub Releases (all platforms)
- Homebrew (macOS/Linux) - planned
- Go install: `go install github.com/darkLord19/wtx@latest`

## Future Roadmap

### v1.1 - Enhanced Features
- [ ] Dev server management (start/stop/logs)
- [ ] Full port conflict resolution
- [ ] JSON output mode for scripting
- [ ] Shell completions (bash, zsh, fish)
- [ ] Worktree templates
- [ ] Custom keybindings

### v1.2 - Integration
- [ ] VS Code extension
- [ ] Raycast extension
- [ ] GitHub CLI plugin
- [ ] Git hooks integration

### v2.0 - Advanced
- [ ] Team workspace sharing
- [ ] Docker workspace isolation
- [ ] Multi-repo support
- [ ] Workspace presets
- [ ] Analytics/insights

## Known Limitations

### Current
- Windows support is partial (editor detection)
- No Windows Terminal integration yet
- Dev server management not implemented
- No workspace templates yet

### By Design
- Requires Git 2.x
- Works only in Git repositories
- One active worktree at a time per editor
- Metadata not synced across machines

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development setup
- Code style guidelines
- Testing requirements
- Pull request process

## Dependencies

### Direct
- github.com/charmbracelet/bubbletea - TUI framework
- github.com/charmbracelet/bubbles - TUI components
- github.com/charmbracelet/lipgloss - Styling
- github.com/spf13/cobra - CLI framework
- github.com/spf13/viper - Configuration

### Build Tools
- Go 1.21+
- Make
- GoReleaser (for releases)

## Metrics & Success

### Goals (3 months)
- 100+ GitHub stars
- 10+ contributors
- Featured in weekly newsletters
- Positive feedback from users

### Tracking
- GitHub stars/forks
- Download counts
- Issue activity
- Community discussions

## Support & Community

- 📖 Documentation: README.md + docs/
- 🐛 Bug Reports: GitHub Issues
- 💡 Feature Requests: GitHub Issues
- 💬 Discussions: GitHub Discussions
- 🤝 Contributing: CONTRIBUTING.md

## License

MIT License - Open source, free to use and modify

## Credits

Built with ❤️ using excellent open source libraries from the Go community, especially the Charm tools (Bubble Tea, Lipgloss) which make beautiful TUI apps possible.

---

**Project Motto**: "Make Git worktrees feel like Cmd+Tab for your branches"
