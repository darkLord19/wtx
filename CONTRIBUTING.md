# Contributing to wtx

Thank you for your interest in contributing to wtx! This document provides guidelines and instructions for contributing.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/wtx.git
   cd wtx
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/darkLord19/wtx.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git 2.x
- Make (optional but recommended)

### Install Dependencies

```bash
make deps
```

### Build and Run

```bash
# Build
make build

# Run
./bin/wtx

# Or build and run in one step
make dev
```

### Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/git/
```

## Code Style

- Follow standard Go formatting (`go fmt`)
- Run `make fmt` before committing
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

## Making Changes

1. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** with clear commit messages:
   ```bash
   git commit -m "Add feature: description of feature"
   ```

3. **Keep commits atomic** - one logical change per commit

4. **Write tests** for new functionality

5. **Update documentation** if needed

## Commit Message Guidelines

Use conventional commits format:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test changes
- `refactor:` - Code refactoring
- `chore:` - Build process or auxiliary tool changes

Examples:
```
feat: add port conflict detection
fix: handle missing upstream branch gracefully
docs: update installation instructions
test: add tests for metadata store
```

## Pull Request Process

1. **Update your branch** with latest upstream:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request** on GitHub

4. **Fill out the PR template** completely

5. **Wait for review** - maintainers will review your PR

6. **Address feedback** if requested

7. **Squash commits** if asked (we prefer clean history)

## PR Checklist

- [ ] Tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] New features have tests
- [ ] Documentation is updated
- [ ] Commit messages follow guidelines
- [ ] Branch is up to date with main
- [ ] No merge conflicts

## Testing Guidelines

### Unit Tests

- Test files should be named `*_test.go`
- Place tests in the same package as the code
- Use table-driven tests for multiple scenarios
- Mock external dependencies

Example:
```go
func TestWorktreeList(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    []Worktree
        wantErr bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Integration Tests

- Place in `test/integration/`
- Test end-to-end workflows
- Clean up test artifacts

## Adding New Features

### New Commands

1. Create a new file in `cmd/wtx/` (e.g., `newcmd.go`)
2. Implement the command using Cobra
3. Register it in `main.go`'s `init()` function
4. Add tests
5. Update README with usage

### New Editor Support

1. Add editor type to `internal/editor/editor.go`
2. Create adapter in `internal/editor/adapters.go`
3. Add to detection priority in `detector.go`
4. Test on target platform
5. Update README

## Documentation

- Update README.md for user-facing changes
- Add inline comments for complex logic
- Update command help text (`--help`)
- Consider adding examples to docs/

## Release Process

(For maintainers)

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create and push tag:
   ```bash
   git tag -a v1.x.x -m "Release v1.x.x"
   git push upstream v1.x.x
   ```
4. GitHub Actions will automatically build and release

## Getting Help

- 💬 [GitHub Discussions](https://github.com/darkLord19/wtx/discussions) for questions
- 🐛 [GitHub Issues](https://github.com/darkLord19/wtx/issues) for bugs
- 📖 Read the [README](README.md) and [documentation](docs/)

## Code of Conduct

Be respectful and inclusive. We're all here to build something useful together.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to wtx! 🎉
