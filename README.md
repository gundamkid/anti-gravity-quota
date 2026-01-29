# 🚀 AG-Quota

A lightweight CLI tool to monitor your Anti-Gravity (Claude Code) AI model quota and usage in real-time.

## Features

- ✅ **Check Quota** - View quota for all AI models (Claude 4 Opus/Sonnet, Gemini 3 Flash/Pro)
- � **Multiple Accounts** - Support for multiple Google accounts with easy switching
- �🔐 **Secure Auth** - Google OAuth2 with PKCE flow
- 📊 **Pretty Output** - Colored tables with visual progress bars
- 🔄 **Auto-Refresh** - Automatic token refresh when expired
- 📝 **JSON Output** - Machine-readable format for scripting
- ⚡ **Fast & Simple** - Single binary, no dependencies

## Installation

```bash
# Build from source
go build -o ag-quota ./cmd/ag-quota

# Or install directly
go install github.com/gundamkid/anti-gravity-quota/cmd/ag-quota@latest
```

## Quick Start

```bash
# Login with Google account
ag-quota login

# Check quota for default account
ag-quota

# List all saved accounts
ag-quota accounts list

# Switch default account
ag-quota accounts switch user@gmail.com

# Check quota for all accounts at once
ag-quota quota --all
```

## Usage

### Check Quota

Display quota information for all AI models:

```bash
$ ag-quota

  ✨ Anti-Gravity Quota Status

  Account: user@example.com
  Project: my-project-123456
  Fetched: 2026-01-26 07:45:23 UTC

  ┌────────────────────────┬───────┬──────────┬──────────┐
  │ Model                  │ Quota │ Reset In │ Status   │
  ├────────────────────────┼───────┼──────────┼──────────┤
  │ Claude 4 Opus          │  85%  │ 4h 23m   │ ✓ OK     │
  │ Claude 4 Sonnet        │ 100%  │ 5h 0m    │ ✓ OK     │
  │ Gemini 3 Flash         │   0%  │ 2h 15m   │ ✗ EMPTY  │
  │ Gemini 3 Pro           │  50%  │ 3h 45m   │ ✓ OK     │
  └────────────────────────┴───────┴──────────┴──────────┘

  ⭐ Default Model: Claude 4 Sonnet
```

**Status Indicators:**
- ✓ OK (green) - Quota above 10%
- ⚠ LOW (yellow) - Quota at or below 10%
- ✗ EMPTY (red) - Quota exhausted

### Multi-Account Support

Check quota for a specific account:

```bash
ag-quota quota --account user@gmail.com
```

Check quota for **all** saved accounts at once:

```bash
ag-quota quota --all
```

### Account Management

Manage your saved Google accounts:

```bash
# List all saved accounts
ag-quota accounts list

# Set the default account
ag-quota accounts default user@gmail.com

# Quickly switch between accounts (alias for default)
ag-quota accounts switch another@gmail.com

# Remove an account
ag-quota accounts remove old@user.com
```

### JSON Output

Get machine-readable output for scripting:

```bash
$ ag-quota --json
# or
$ ag-quota quota --json

{
  "Email": "user@example.com",
  "ProjectID": "my-project-123456",
  "Models": [
    {
      "ModelID": "claude-sonnet-4-5",
      "DisplayName": "Claude 4 Sonnet",
      "Label": "Claude 4 Sonnet",
      "Provider": "claude",
      "RemainingFraction": 0.85,
      "ResetTime": "2026-01-26T12:00:00Z",
      "IsExhausted": false
    }
  ],
  "DefaultModelID": "claude-sonnet-4-5",
  "FetchedAt": "2026-01-26T07:45:23Z"
}
```

### Login

Authenticate with your Google account:

```bash
$ ag-quota login

Starting authentication flow...

Opening browser for authentication...
If browser doesn't open, visit this URL:
https://accounts.google.com/o/oauth2/v2/auth?...

Login successful!
Logged in as: user@example.com
```

### Check Status

View current authentication status:

```bash
$ ag-quota status

Authentication Status
====================

✓ Logged in as: user@example.com
✓ Token valid for: 58m

Config directory: /home/user/.config/ag-quota
```

### Logout

Clear stored authentication tokens:

```bash
$ ag-quota logout

✓ Logged out successfully
```

## Commands

| Command | Description | Flags |
|---------|-------------|-------|
| `ag-quota` | Check quota (default account) | `--json, -j` |
| `ag-quota quota` | Check quota | `--account, --all, --json` |
| `ag-quota accounts list` | List all saved accounts | |
| `ag-quota accounts default` | Set the default account | |
| `ag-quota accounts switch` | Alias for `accounts default` | |
| `ag-quota accounts remove` | Remove a saved account | |
| `ag-quota login` | Authenticate with Google account | |
| `ag-quota status` | Show authentication status | |
| `ag-quota logout` | Clear stored tokens | |
| `ag-quota --help` | Show help message | |
| `ag-quota --version` | Show version information | |

## Configuration

Configuration files are stored in the following locations:

- **Linux/WSL**: `~/.config/ag-quota/`
  - Token: `~/.config/ag-quota/token.json`
- **macOS**: `~/.config/ag-quota/` (or `~/Library/Application Support/ag-quota/`)
  - Token: `~/.config/ag-quota/token.json`

**Token File Permissions:**
- Token files are stored with `0600` permissions (owner read/write only)
- Never commit token files to version control

**Environment Variables:**
- `XDG_CONFIG_HOME` - Override default config directory (Linux/macOS)

## How It Works

1. **Authentication**: Uses Google OAuth2 with PKCE flow for secure authentication
2. **Token Storage**: Stores OAuth tokens locally with automatic refresh
3. **API Integration**: Calls Google Cloud Code API endpoints:
   - `POST /v1internal:loadCodeAssist` - Retrieves project information
   - `POST /v1internal:fetchAvailableModels` - Fetches model quotas
4. **Display**: Formats and displays quota information with visual indicators

## API Details

**Base URL:** `https://cloudcode-pa.googleapis.com`

**Authentication:**
- OAuth2 with PKCE (Proof Key for Code Exchange)
- Client ID: Public Google Cloud Code extension client
- Scopes: `openid email profile cloud-platform`

**Rate Limiting:**
- Automatic retry with exponential backoff
- Max 3 retry attempts for failed requests

## Requirements

- Go 1.21+ (for building from source)
- Internet connection
- Google account with Anti-Gravity (Claude Code) access
## Troubleshooting

### "Not logged in" error

Run `ag-quota login` to authenticate with your Google account.

### Token expired

Tokens are automatically refreshed. If you see auth errors, try:
```bash
ag-quota logout
ag-quota login
```

### API errors or rate limiting

The tool implements automatic retry with exponential backoff. If you consistently see errors:
- Check your internet connection
- Verify your Google account has Anti-Gravity access
- Wait a few minutes if rate limited

### Browser doesn't open during login

If the browser doesn't open automatically, copy and paste the URL shown in the terminal.

## Development

### Building

```bash
# Build for current platform
make build
# or
go build -o ag-quota ./cmd/ag-quota
```

### Testing

```bash
# Run all tests
make test
# or
go test -v ./...
```

### Project Structure

```
anti-gravity-quota/
├── cmd/ag-quota/       # CLI entry point
├── internal/
│   ├── api/           # Cloud Code API client
│   ├── auth/          # OAuth2 authentication
│   ├── config/        # Configuration management
│   ├── models/        # Data models
│   └── ui/            # Display formatting
├── docs/              # Documentation
│   ├── build-flow.md  # Build and CI/CD flow
│   └── technical.md   # API & Technical details
└── README.md
```

## Credits

Inspired by [antigravity-usage](https://github.com/skainguyen1412/antigravity-usage) by skainguyen1412.

Built with:
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [color](https://github.com/fatih/color) - Colored terminal output
- [oauth2](https://golang.org/x/oauth2) - OAuth2 client
- [go-pretty](https://github.com/jedib0t/go-pretty) - Table output

## License

MIT License - see LICENSE file for details
