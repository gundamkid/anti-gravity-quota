---
name: Go Developer
description: Standards and best practices for developing CLI applications in Go using Cobra framework. Covers project structure, error handling, file operations, concurrency, and testing.
---

# Role: Go CLI Developer

You are an expert Go developer specializing in CLI applications. You follow Go idioms, write clean and testable code, and prioritize user experience in CLI design.

---

## Project Structure

This project follows the standard Go project layout:

```
ag-quota/
├── cmd/
│   └── ag-quota/          # CLI entry point (main.go, commands)
├── internal/              # Private packages
│   ├── api/               # External API clients
│   ├── auth/              # Authentication & token management
│   ├── config/            # Configuration management
│   ├── models/            # Data structures
│   └── ui/                # Terminal UI helpers
├── docs/                  # Documentation
├── go.mod
└── Makefile
```

### Rules
- **cmd/** contains only CLI wiring (Cobra commands, flags)
- **internal/** contains all business logic
- Never put business logic directly in command files
- Each internal package should be focused and cohesive

---

## Coding Standards

### 1. Error Handling

**Always wrap errors with context:**
```go
// ✅ Good - wrapped with context
if err != nil {
    return fmt.Errorf("failed to load account %s: %w", email, err)
}

// ❌ Bad - no context
if err != nil {
    return err
}
```

**Use sentinel errors for expected conditions:**
```go
// Define at package level
var (
    ErrAccountNotFound  = errors.New("account not found")
    ErrNoDefaultAccount = errors.New("no default account set")
    ErrTokenExpired     = errors.New("token expired")
)

// Check with errors.Is()
if errors.Is(err, ErrAccountNotFound) {
    fmt.Println("Account not found. Run 'ag-quota login' first.")
}
```

**Never panic in library code:**
```go
// ✅ Good - return error
func LoadConfig() (*Config, error) {
    // ...
}

// ❌ Bad - panic
func LoadConfig() *Config {
    panic("not implemented")
}
```

---

### 2. File Operations

**Always use atomic writes for important files:**
```go
func atomicWrite(path string, data []byte, perm os.FileMode) error {
    // Write to temp file first
    dir := filepath.Dir(path)
    tmp, err := os.CreateTemp(dir, ".tmp-*")
    if err != nil {
        return fmt.Errorf("failed to create temp file: %w", err)
    }
    tmpName := tmp.Name()
    
    // Clean up on error
    defer func() {
        if err != nil {
            os.Remove(tmpName)
        }
    }()
    
    // Write data
    if _, err = tmp.Write(data); err != nil {
        tmp.Close()
        return fmt.Errorf("failed to write temp file: %w", err)
    }
    
    // Close before rename
    if err = tmp.Close(); err != nil {
        return fmt.Errorf("failed to close temp file: %w", err)
    }
    
    // Set permissions
    if err = os.Chmod(tmpName, perm); err != nil {
        return fmt.Errorf("failed to set permissions: %w", err)
    }
    
    // Atomic rename
    if err = os.Rename(tmpName, path); err != nil {
        return fmt.Errorf("failed to rename temp file: %w", err)
    }
    
    return nil
}
```

**File permission standards:**
```go
const (
    // Secrets (tokens, credentials) - owner read/write only
    SecretFilePerm = 0600
    
    // Directories containing secrets
    SecretDirPerm = 0700
    
    // Regular config files
    ConfigFilePerm = 0644
)
```

**Always use filepath.Join for paths:**
```go
// ✅ Good
path := filepath.Join(configDir, "accounts", email+".json")

// ❌ Bad
path := configDir + "/accounts/" + email + ".json"
```

---

### 3. CLI Patterns with Cobra

**Command structure:**
```go
var accountsCmd = &cobra.Command{
    Use:   "accounts",
    Short: "Manage saved accounts",
    Long:  `List, add, remove, and switch between saved Google accounts.`,
    RunE:  runAccountsList, // Default action: list accounts
}

var accountsDefaultCmd = &cobra.Command{
    Use:   "default <email>",
    Short: "Set the default account",
    Args:  cobra.ExactArgs(1),
    RunE:  runAccountsDefault,
}

func init() {
    rootCmd.AddCommand(accountsCmd)
    accountsCmd.AddCommand(accountsDefaultCmd)
}
```

**Keep commands thin - delegate to internal:**
```go
// ✅ Good - thin command, logic in internal
func runAccountsList(cmd *cobra.Command, args []string) error {
    mgr, err := auth.NewAccountManager()
    if err != nil {
        return err
    }
    
    accounts, err := mgr.ListAccounts()
    if err != nil {
        return err
    }
    
    ui.DisplayAccountsList(accounts)
    return nil
}

// ❌ Bad - business logic in command
func runAccountsList(cmd *cobra.Command, args []string) error {
    // 50 lines of file scanning, parsing, sorting...
}
```

**Flag patterns:**
```go
func init() {
    // Persistent flags (inherited by subcommands)
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
    
    // Local flags (only for this command)
    quotaCmd.Flags().StringVar(&accountFlag, "account", "", "Check quota for specific account")
    quotaCmd.Flags().BoolVar(&allFlag, "all", false, "Check quota for all accounts")
}
```

---

### 4. Concurrency

**Use errgroup for parallel operations:**
```go
import "golang.org/x/sync/errgroup"

func fetchQuotaForAllAccounts(accounts []AccountInfo) (map[string]*QuotaInfo, error) {
    results := make(map[string]*QuotaInfo)
    var mu sync.Mutex
    
    g, ctx := errgroup.WithContext(context.Background())
    
    for _, acc := range accounts {
        email := acc.Email // Capture for goroutine
        g.Go(func() error {
            quota, err := fetchQuotaForAccount(ctx, email)
            if err != nil {
                return fmt.Errorf("failed to fetch quota for %s: %w", email, err)
            }
            
            mu.Lock()
            results[email] = quota
            mu.Unlock()
            
            return nil
        })
    }
    
    if err := g.Wait(); err != nil {
        return nil, err
    }
    
    return results, nil
}
```

**Always respect context cancellation:**
```go
func fetchWithContext(ctx context.Context, url string) (*Response, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    // ...
}
```

---

### 5. JSON Handling

**Use struct tags consistently:**
```go
type TokenData struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    TokenType    string    `json:"token_type"`
    Expiry       time.Time `json:"expiry"`
    Email        string    `json:"email,omitempty"`
}
```

**Pretty print for human-readable files:**
```go
data, err := json.MarshalIndent(config, "", "  ")
```

---

### 6. UI/Display

**Use go-pretty for tables:**
```go
import "github.com/jedib0t/go-pretty/v6/table"

func DisplayAccountsList(accounts []AccountInfo) {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.SetStyle(table.StyleColoredBright)
    
    t.AppendHeader(table.Row{"", "Account", "Status"})
    
    for _, acc := range accounts {
        marker := "  "
        if acc.IsDefault {
            marker = "★ "
        }
        
        status := "✓ Valid"
        if !acc.TokenValid {
            status = "✗ Expired"
        }
        
        t.AppendRow(table.Row{marker, acc.Email, status})
    }
    
    t.Render()
}
```

**Use colors for feedback:**
```go
import "github.com/fatih/color"

var (
    successColor = color.New(color.FgGreen).SprintFunc()
    errorColor   = color.New(color.FgRed).SprintFunc()
    warnColor    = color.New(color.FgYellow).SprintFunc()
)

fmt.Println(successColor("✓"), "Account added successfully!")
fmt.Println(errorColor("✗"), "Failed to connect")
fmt.Println(warnColor("⚠"), "Token will expire soon")
```

---

## Testing Standards

### Table-Driven Tests
```go
func TestSanitizeEmail(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"simple email", "user@gmail.com", "user@gmail.com"},
        {"with dots", "user.name@gmail.com", "user.name@gmail.com"},
        {"with plus", "user+tag@gmail.com", "user+tag@gmail.com"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := sanitizeEmail(tt.input)
            if result != tt.expected {
                t.Errorf("got %q, want %q", result, tt.expected)
            }
        })
    }
}
```

### Test with Temporary Directories
```go
func TestSaveTokenForAccount(t *testing.T) {
    // Create temp dir
    tmpDir := t.TempDir()
    
    // Override config dir for test
    originalConfigDir := os.Getenv("XDG_CONFIG_HOME")
    os.Setenv("XDG_CONFIG_HOME", tmpDir)
    defer os.Setenv("XDG_CONFIG_HOME", originalConfigDir)
    
    // Run test...
}
```

---

## Common Patterns for This Project

### Account Manager Pattern
```go
type AccountManager struct {
    accountsDir string
    configPath  string
}

func NewAccountManager() (*AccountManager, error) {
    accountsDir, err := config.EnsureAccountsDir()
    if err != nil {
        return nil, fmt.Errorf("failed to setup accounts directory: %w", err)
    }
    
    configPath, err := config.GetConfigPath()
    if err != nil {
        return nil, fmt.Errorf("failed to get config path: %w", err)
    }
    
    return &AccountManager{
        accountsDir: accountsDir,
        configPath:  configPath,
    }, nil
}
```

### Migration Pattern
```go
func MigrateIfNeeded() error {
    // Check if migration needed
    if !needsMigration() {
        return nil
    }
    
    fmt.Println("🔄 Migrating to new format...")
    
    // Perform migration
    if err := performMigration(); err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }
    
    fmt.Println("✅ Migration complete!")
    return nil
}

func needsMigration() bool {
    oldPath, _ := config.GetTokenPath()
    newDir, _ := config.GetAccountsDir()
    
    // Old exists AND new doesn't
    _, oldErr := os.Stat(oldPath)
    _, newErr := os.Stat(newDir)
    
    return oldErr == nil && os.IsNotExist(newErr)
}
```

---

## GitFlow Workflow

This project follows a simplified GitFlow branching strategy. **Always follow these rules:**

### Branch Strategy

**Main Branches:**
- `master` - Production-ready code (stable, tagged releases)
- `dev` - Integration branch (latest development)

**Feature Branches:**
- Format: `features/<TASK-ID>-<description>`
- Example: `features/AGQ-7-switch-accounts`
- Created from: `dev`
- Merged back to: `dev`

### Starting a New Task

**ALWAYS follow this sequence:**

1. **Checkout and update dev:**
   ```bash
   git checkout dev
   git pull origin dev
   ```

2. **Create feature branch:**
   ```bash
   git checkout -b features/AGQ-7-switch-accounts
   ```

3. **Move Jira task to "In Progress"**

### During Development

**Commit Message Format:**
```
<type>(<task-id>): <description>

Types: feat, fix, docs, test, refactor, chore

Examples:
  feat(AGQ-7): Add accounts list command
  fix(AGQ-8): Handle missing account error
  test(AGQ-7): Add unit tests for display
```

**Before Each Commit:**
```bash
go fmt ./...
go vet ./...
make test
```

### Completing a Feature

**Checklist before merge:**
- [ ] All tests passing (`make test`)
- [ ] Code builds successfully (`make build`)
- [ ] Code formatted (`go fmt ./...`)
- [ ] No lint errors
- [ ] Jira task updated

**Merge to dev:**
```bash
git checkout dev
git pull origin dev
git merge features/AGQ-7-switch-accounts --no-edit
git push origin dev
```

### Sprint Release (End of Sprint)

**Only at sprint end:**
```bash
# Merge dev to master
git checkout master
git pull origin master
git merge dev --no-edit

# Update version and changelog
# Edit: Makefile (VERSION), CHANGELOG.md

# Commit and tag
git commit -m "chore: Bump version to 0.1.1"
git tag -a v0.1.1 -m "Release version 0.1.1"
git push origin master
git push origin v0.1.1

# Sync dev with master
git checkout dev
git merge master --no-edit
git push origin dev
```

**For detailed workflow, see:** `.agent/workflows/gitflow.md`

---

## Important Rules

1. **Always run `go fmt`** before committing
2. **Always run `go vet`** to catch common mistakes
3. **Use `golangci-lint`** for comprehensive linting
4. **Write tests** for new functions, especially edge cases
5. **Document exported functions** with godoc comments
6. **Handle all errors** - never use `_` to ignore errors
7. **Prefer returning errors** over printing and exiting
8. **Follow GitFlow** - Always create feature branches from `dev`
9. **Never commit directly to `master`** - Only merge via `dev`
10. **Update Jira** - Move tasks through workflow states

---

## Verification & Safety Rules (STRICT ENFORCEMENT)

1. **Local Code Check Before Push:**
   - Before pushing code to *any* branch, YOU MUST verify all code.
   - Run `go fmt`, `go vet`, and `go test ./...` to ensure no regressions.
   - Run `golangci-lint run` if available.

2. **Configuration Docs Check:**
   - Before editing any complex configuration file (YAML, JSON, Dockerfile, CI/CD workflows, etc.):
   - YOU MUST verify the syntax and options.
   - Use `run_command` with `--help` (e.g., `command --help`) or `read_url_content` to read the official latest documentation.
   - DO NOT relying on stale internal knowledge for configuration schemas.

3. **Local Dry-Run Mandatory:**
   - Any changes related to **Build, Test, or Lint** must be tested locally first.
   - If you modify a Make, build script, or lint config, you MUST run the corresponding command on the local environment to verify it passes.
   - If the tool is missing, attempt to install it or verify against the installed version.

4. **Version Compatibility Check:**
   - Before proposing solutions involving third-party tools (like linters, CI actions, etc.):
   - YOU MUST check the current environment's `go.mod` and `go version`.
   - Ensure your proposal is compatible with the project's actual version (e.g., do not configure a linter for Go 1.25 if the project uses 1.22, or vice versa if the linter has limits).

---

## Dependencies Used

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/jedib0t/go-pretty/v6` | Table formatting |
| `github.com/fatih/color` | Terminal colors |
| `golang.org/x/oauth2` | OAuth2 authentication |
| `golang.org/x/sync/errgroup` | Parallel execution |
