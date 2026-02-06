# WTX Quick Start Guide

**Goal**: Be productive with wtx in 5 minutes ⏱️

---

## ✅ Prerequisites Checklist

Before starting, make sure you have:

- [ ] **Git 2.x installed** - Check with `git --version`
- [ ] **Go 1.21+ installed** - Check with `go version`  
- [ ] **A Git repository** to work in

Don't have these? See [Installation Help](#installation-help) below.

---

## 📦 Step 1: Install wtx (2 minutes)

### Option A: Install from Source (Recommended)

```bash
# Clone the repository
git clone https://github.com/darkLord19/wtx.git
cd wtx

# Install to $GOPATH/bin
make install

# Verify installation
wtx --help
```

### Option B: Install with Go

```bash
# Direct install
go install github.com/darkLord19/wtx/cmd/wtx@latest

# Verify installation
wtx --help
```

### Option C: Build Locally

```bash
# Clone and build
git clone https://github.com/darkLord19/wtx.git
cd wtx
make build

# Run from bin directory
./bin/wtx --help

# Optional: Add to PATH
export PATH="$PATH:$(pwd)/bin"
```

**✓ Success**: Running `wtx --help` shows usage information.

---

## 🚀 Step 2: First Run (1 minute)

Navigate to any Git repository and run wtx:

```bash
cd ~/projects/my-app
wtx
```

**What happens**:
1. ✨ **Setup Wizard** launches automatically
2. 📝 You'll configure a few settings
3. 🎉 wtx is ready to use!

### Setup Wizard Steps

The wizard will ask you to configure:

#### 1️⃣ **Choose Your Editor**

```
┌────────────────────────────────────┐
│ Select your preferred editor:      │
│                                    │
│ ▸ (auto-detect)                   │
│   visual studio code               │
│   cursor                           │
│   neovim                           │
│   vim                              │
│   (custom)                         │
└────────────────────────────────────┘
```

**Recommendation**: Choose "(auto-detect)" unless you have a specific preference.

#### 2️⃣ **Worktree Directory**

```
┌────────────────────────────────────┐
│ Where should worktrees be created? │
│                                    │
│ ../worktrees                       │
│                                    │
│ (relative to repository root)      │
└────────────────────────────────────┘
```

**Recommendation**: Keep the default `../worktrees`.

**What this means**: If your repo is at `~/projects/my-app`, worktrees will be created at `~/projects/worktrees/`.

#### 3️⃣ **Window Reuse**

```
┌────────────────────────────────────┐
│ Reuse existing editor window?      │
│                                    │
│ [ Yes ]  [ No ]                    │
│                                    │
│ If yes, opening a worktree will    │
│ reuse your current editor window   │
└────────────────────────────────────┘
```

**Recommendation**: Choose "Yes" for seamless workflow.

---

## 🎯 Step 3: Create Your First Worktree (1 minute)

```bash
# Create a worktree for a new feature
wtx add my-first-feature
```

**What happens**:
```
Creating worktree 'my-first-feature' for branch 'my-first-feature'...
✓ Created worktree: my-first-feature
  Path: /Users/you/projects/worktrees/my-first-feature
  Branch: my-first-feature

Open in editor now? [Y/n]:
```

Press **Enter** (or type `y`) to open in your editor!

### Understanding Worktrees

A worktree is like having a second copy of your repository, but:
- ✅ Shares the same .git history (disk efficient)
- ✅ Can be on a different branch
- ✅ Won't affect your main working directory
- ✅ Can run different dev servers simultaneously

---

## 🔄 Step 4: Switch Between Worktrees (1 minute)

```bash
# Launch the interactive switcher
wtx
```

You'll see something like this:

```
┌─────────────────────────────────────────────┐
│ Workspace Manager (my-app)                  │
├─────────────────────────────────────────────┤
│ main                ● clean                 │
│ my-first-feature    ● clean                 │
├─────────────────────────────────────────────┤
│ Press enter to open • q to quit             │
└─────────────────────────────────────────────┘
```

**Navigation**:
- `↑/↓` or `j/k` - Move selection
- `/` - Search/filter
- `Enter` - Open selected worktree
- `q` or `Esc` - Quit

### Quick Commands

If you know the worktree name:

```bash
# Open specific worktree directly
wtx open my-first-feature

# List all worktrees
wtx list

# Show detailed status
wtx status my-first-feature

# Remove a worktree
wtx rm old-feature
```

---

## 🎓 Step 5: Your First Workflow (bonus)

Let's simulate a real workflow:

```bash
# 1. Create feature branch
wtx add feat-user-login

# 2. Make some changes
cd /path/to/worktrees/feat-user-login
echo "console.log('login')" > login.js
git add login.js
git commit -m "feat: add login functionality"

# 3. Urgent bug comes in
wtx add hotfix-critical --from main

# 4. Fix the bug
cd /path/to/worktrees/hotfix-critical
# ... make fixes ...
git commit -am "fix: critical bug"

# 5. Back to feature
wtx open feat-user-login

# 6. Clean up when done
wtx rm hotfix-critical
```

**Key insight**: No stashing, no context switching headaches!

---

## 🎨 Explore the Full TUI

For the complete experience with tabs:

```bash
wtx --tui
# or
wtx -t
```

This shows three tabs:

**[1] Worktrees** - Browse and open  
**[2] Manage** - Create, delete, prune  
**[3] Settings** - Configure wtx

**Shortcuts**:
- `1`, `2`, `3` - Switch tabs
- `?` - Show help
- `q` - Quit

---

## 📊 Verify Everything Works

Run this quick check:

```bash
# Should show your worktrees
wtx list

# Should show configuration
wtx config

# Should work without errors
wtx
```

**Expected output**:
- `list`: Shows at least your main worktree
- `config`: Shows your editor and settings
- Interactive mode launches successfully

---

## 🎯 What's Next?

### Learn More
- **Full guide**: [README.md](../README.md)
- **Common workflows**: [docs/WORKFLOWS.md](WORKFLOWS.md)
- **FAQ**: [docs/FAQ.md](FAQ.md)

### Customize
```bash
# Edit settings interactively
wtx config --tui

# Or manually edit
vim ~/.config/wtx/config.json
```

### Get Productive
```bash
# Set up an alias
echo "alias w='wtx'" >> ~/.zshrc
source ~/.zshrc

# Now just type 'w' to switch worktrees!
w
```

---

## 💡 Pro Tips

### Tip 1: Quick Switching
```bash
# Instead of:
wtx
# (select worktree)
# (press enter)

# Do this:
wtx open feat-login
```

### Tip 2: Branch from Remote
```bash
# Review a PR without leaving your work
wtx add review-pr-123 --from origin/pull/123/head
```

### Tip 3: Keep It Clean
```bash
# Weekly cleanup
wtx prune --days 7

# Or be aggressive
wtx prune --days 3
```

### Tip 4: See Status at a Glance
```bash
wtx list
```
Shows clean/dirty status and ahead/behind commits.

---

## 🆘 Troubleshooting

### "not a git repository"
**Solution**: Run wtx from inside a Git repository.
```bash
cd ~/projects/your-repo
wtx
```

### "git is not installed"
**Solution**: Install Git for your platform.
```bash
# macOS
brew install git

# Ubuntu/Debian
sudo apt install git

# Windows
# Download from https://git-scm.com/
```

### "command not found: wtx"
**Solution**: Add to PATH or use full path.
```bash
# Check where it was installed
which wtx

# If not found, add Go bin to PATH
export PATH="$PATH:$HOME/go/bin"
```

### Editor doesn't open
**Solutions**:
1. Check if editor is in PATH: `which code`
2. Set manually: `wtx config editor code`
3. Try auto-detect: `wtx config editor ""`

### Want to start over?
```bash
# Remove config and try again
rm ~/.config/wtx/config.json
wtx  # Will run setup wizard again
```

---

## 📚 Installation Help

### Install Git

**macOS**:
```bash
brew install git
# or
xcode-select --install
```

**Ubuntu/Debian**:
```bash
sudo apt update
sudo apt install git
```

**Fedora**:
```bash
sudo dnf install git
```

**Windows**:
Download from [https://git-scm.com/](https://git-scm.com/)

### Install Go

**macOS**:
```bash
brew install go
```

**Ubuntu/Debian**:
```bash
# Add Go PPA
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go
```

**Any Platform**:
Download from [https://go.dev/dl/](https://go.dev/dl/)

---

## 🎉 Success!

You're now ready to use wtx! 

**Remember**:
- Use `wtx` for interactive switching
- Use `wtx add <n>` to create worktrees
- Use `wtx list` to see all worktrees
- Use `wtx prune` to clean up old ones

**Get Help**:
- Press `?` in the TUI for keyboard shortcuts
- Run `wtx <cmd> --help` for command help
- Check [FAQ.md](FAQ.md) for common questions
- Open [GitHub Issues](https://github.com/darkLord19/wtx/issues) for bugs

---

**Happy coding!** 🚀

*Made with ❤️ for developers who love Git worktrees*
