# Contributing to Anti-Gravity Quota

Thank you for your interest in contributing to **Anti-Gravity Quota**! This guide will help you get started with our development process.

## 🚀 Getting Started

1. **Fork the repository** on GitHub.
2. **Clone your fork** locally.
3. **Install dependencies**: `go mod tidy`.
4. **Run tests** to ensure everything is working: `make test`.

## 🛠 Development Workflow

We follow a strict **GitFlow** and **Jira**-driven development process.

### 1. Branching Model
- `master`: Production-ready code (tags only).
- `dev`: Integration branch for features.
- `features/AGQ-XXX-description`: Feature branches created from `dev`.

### 2. Commit Message Standards
We use **Conventional Commits** to support automatic changelog generation.
Format: `<type>(<scope>): <description> (AGQ-XXX)`

**Types:**
- `feat`: New feature for the user.
- `fix`: Bug fix for the user.
- `docs`: Documentation changes.
- `style`: Formatting, missing semi-colons, etc; no production code change.
- `refactor`: Refactoring production code, eg. renaming a variable.
- `perf`: Code change that improves performance.
- `test`: Adding missing tests, refactoring tests; no production code change.
- `chore`: Updating grunt tasks etc; no production code change.

**Example:**
`feat(ui): add compact mode support (AGQ-40)`

### 3. Pull Request Process
- Ensure all tests pass: `make test`.
- Ensure code is linted: `make lint`.
- Targeted branch should be `dev`.
- Include the Jira issue link in your description.

## 🏗 Coding Standards

- **Error Handling**: Always wrap errors with context using `fmt.Errorf("...: %w", err)`.
- **Formatting**: Always run `go fmt` before committing.
- **Concurrency**: Use `golang.org/x/sync/errgroup` for parallel tasks.
- **UI**: Use `github.com/jedib0t/go-pretty/v6/table` for table displays.

## 📦 Building and Testing

```bash
# Build the binary
make build

# Run unit tests
make test

# Run linter
make lint

# Generate changelog
make changelog
```

## ❓ Need Help?

If you have questions, feel free to open an issue or reach out to the project maintainers.

Happy coding! 🐉⚡
