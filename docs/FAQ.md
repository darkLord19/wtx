# Frequently Asked Questions (FAQ)

## General Questions

### Q: How is wtx different from using `git worktree` directly?

**A:** wtx enhances git worktrees with:

1. **Beautiful TUI** - Interactive fuzzy search and navigation
2. **Editor Integration** - Opens worktrees directly in VS Code, Cursor, Neovim, etc.
3. **Safety Checks** - Prevents accidental deletion of uncommitted work
4. **Metadata Tracking** - Tracks usage statistics, last opened times
5. **Stale Cleanup** - Identifies and removes unused worktrees
6. **Multi-tab Interface** - Manage, create, delete, and configure in one place

### Q: Can I use wtx with existing worktrees?

**A:** Yes! wtx automatically detects all worktrees created by git, whether you made them with `git worktree add` or wtx. All operations are fully compatible.

### Q: Does this work with GitHub flow / GitLab flow / other workflows?

**A:** Absolutely. wtx is workflow-agnostic and enhances any Git workflow. It's particularly useful for:
- Feature branch workflows
- Hotfix workflows
- PR review workflows
- Parallel development

### Q: What happens to my work if I uninstall wtx?

**A:** Nothing changes! Your worktrees are standard git worktrees. The only thing stored by wtx is metadata in `.git/wtx-meta.json`, which can be safely deleted without affecting your work.

## Usage Questions

### Q: Can I use different editors for different worktrees?

**A:** Yes! You can override the editor per-session:
```bash
EDITOR=vim wtx open feature-auth
EDITOR=code wtx open feature-frontend
```

Or set custom commands in the config for specific worktrees.

### Q: How do I use wtx with forks and remotes?

**A:** Use the `--from` flag to specify any branch:
```bash
# From remote branch
wtx add review-pr --from origin/pull-request-123

# From specific remote
wtx add upstream-feature --from upstream/develop

# From local branch
wtx add hotfix --from main
```

### Q: Can multiple people on my team use wtx?

**A:** Yes! wtx metadata (`.git/wtx-meta.json`) is local and not committed to the repository. Each team member has their own metadata and configuration.

### Q: How do I handle merge conflicts in worktrees?

**A:** Worktrees are independent git working directories. Handle conflicts the same way:
```bash
cd /path/to/worktree
git merge main
# Resolve conflicts
git add .
git commit
```

## Performance Questions

### Q: How many worktrees can I have?

**A:** Technical limit is very high, but we recommend:
- **< 20 active worktrees** - Optimal TUI performance (<150ms)
- **20-50 worktrees** - Slight slowdown (150-300ms)
- **50+ worktrees** - Consider pruning old ones

Use `wtx prune` regularly to maintain performance.

### Q: Why is the TUI slow to start?

**Possible causes:**
1. **Too many worktrees** - Prune stale ones: `wtx prune`
2. **Large repository** - Status checks take time on huge repos
3. **Network filesystem** - Git operations slower on NFS/network drives

**Solutions:**
- Use `wtx open <name>` to skip TUI
- Reduce active worktrees to <20
- Set up indexing for large repos

### Q: Does wtx work on network drives?

**A:** Yes, but performance may be slower due to network latency. For best performance, keep repositories on local storage.

## Configuration Questions

### Q: Where is the configuration stored?

**A:** Two locations:
- **Config**: `~/.config/wtx/config.json` (user settings)
- **Metadata**: `.git/wtx-meta.json` (per-repo, local only)
- **Logs**: `~/.cache/wtx/wtx.log` (debugging)

### Q: How do I reset wtx to defaults?

**A:** Remove the config file:
```bash
rm ~/.config/wtx/config.json
wtx  # Will run setup wizard again
```

Or run the setup wizard manually:
```bash
wtx setup
```

### Q: Can I use wtx without the TUI?

**A:** Yes! All commands work non-interactively:
```bash
wtx list                    # Show all worktrees
wtx add feature-name        # Create worktree
wtx open feature-name       # Open specific worktree
wtx rm feature-name         # Remove worktree
wtx status feature-name     # Show detailed status
```

## Platform-Specific Questions

### Q: Does wtx work on Windows?

**A:** Partial support. The CLI works, but:
- ✅ Core commands (add, list, open, rm) work
- ✅ TUI works in Windows Terminal
- ⚠️ Editor detection may need manual configuration
- ⚠️ Some paths may need adjustment

Full Windows support is planned.

### Q: Does wtx work on macOS?

**A:** Yes! Fully supported on macOS 10.15+. Install via:
```bash
# From source
make install

# Or with Go
go install github.com/darkLord19/wtx/cmd/wtx@latest
```

### Q: What about Linux?

**A:** Fully supported on all major distributions (Ubuntu, Fedora, Arch, etc.).

## Troubleshooting

### Q: "not a git repository" error

**Solution:** Run wtx from inside a git repository:
```bash
cd /path/to/your/repo
wtx
```

### Q: "git is not installed" error

**Solution:** Install Git 2.x for your platform:
- **macOS**: `brew install git`
- **Ubuntu/Debian**: `sudo apt install git`
- **Fedora**: `sudo dnf install git`
- **Windows**: Download from git-scm.com

### Q: Worktree won't delete - "uncommitted changes"

**Solutions:**
1. **Commit changes**:
   ```bash
   cd /path/to/worktree
   git commit -am "WIP: save work"
   ```

2. **Stash changes**:
   ```bash
   cd /path/to/worktree
   git stash
   ```

3. **Force delete** (⚠️ loses changes):
   ```bash
   wtx rm worktree-name --force
   ```

### Q: Editor doesn't open

**Check these:**
1. Is editor in PATH? `which code` / `which cursor`
2. Is config correct? `wtx config`
3. Try setting manually: `wtx config editor code`

**For custom editors:**
```bash
wtx config editor custom_command myeditor "/usr/local/bin/myeditor --open"
```

### Q: How do I see debug logs?

**A:** Logs are stored at:
- **macOS/Linux**: `~/.cache/wtx/wtx.log`
- **Windows**: `%LOCALAPPDATA%\wtx\wtx.log`

View logs:
```bash
tail -f ~/.cache/wtx/wtx.log
```

## Advanced Usage

### Q: Can I automate wtx with scripts?

**A:** Yes! All commands support non-interactive mode:
```bash
#!/bin/bash
# Automated worktree management

# Create worktree
wtx add feature-$TICKET_ID --from develop

# Open in editor (non-interactive)
wtx open feature-$TICKET_ID

# When done, cleanup
wtx rm feature-$TICKET_ID --force
```

### Q: How do I backup wtx metadata?

**A:** Copy the metadata file:
```bash
# Backup
cp .git/wtx-meta.json .git/wtx-meta.json.backup

# Restore
cp .git/wtx-meta.json.backup .git/wtx-meta.json
```

Or commit it to a private branch (not recommended for teams).

### Q: Can I customize keyboard shortcuts?

**A:** Not yet. Custom keybindings are planned for v1.1. Current shortcuts are:
- `?` - Toggle help
- `q` / `esc` - Quit
- `1-3` - Switch tabs
- `c` - Create
- `d` - Delete
- `p` - Prune
- `r` - Refresh

### Q: How do I integrate wtx with my IDE?

**A:** For VS Code/Cursor, wtx automatically uses the `-r` flag to reuse windows. For custom integration:

**Terminal task** (VS Code):
```json
{
  "label": "wtx: Open Worktree",
  "type": "shell",
  "command": "wtx"
}
```

**Keybinding** (VS Code):
```json
{
  "key": "cmd+shift+w",
  "command": "workbench.action.tasks.runTask",
  "args": "wtx: Open Worktree"
}
```

## Getting Help

### Q: Where can I get more help?

- 📖 **Documentation**: [README.md](../README.md)
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/darkLord19/wtx/issues)
- 💡 **Feature Requests**: [GitHub Issues](https://github.com/darkLord19/wtx/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/darkLord19/wtx/discussions)
- 📧 **Email**: Check repository for maintainer contact

### Q: How can I contribute?

See [CONTRIBUTING.md](../CONTRIBUTING.md) for:
- Development setup
- Coding guidelines
- Pull request process
- Feature suggestions

### Q: Is there a Discord/Slack community?

Not yet! If there's interest, we'll create one. Join the discussion on GitHub Discussions.

---

**Don't see your question?** Open an issue or discussion on GitHub!
