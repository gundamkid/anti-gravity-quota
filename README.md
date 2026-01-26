# 🚀 AG-Quota

A lightweight CLI tool to check your Anti-Gravity AI model quota and usage.

## Features

- ✅ **Check Quota** - View quota for all AI models (Claude, Gemini)
- 🔐 **Cloud Login** - Authenticate via Google OAuth
- 📊 **Pretty Output** - Colored tables with quota details
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

# Check quota for all models
ag-quota

# Check authentication status
ag-quota status

# Logout
ag-quota logout
```

## Usage

### Check Quota

```bash
$ ag-quota

╔══════════════════════════════════════════════════════════╗
║           Anti-Gravity Quota Status                       ║
╠══════════════════════════════════════════════════════════╣
║  Account: user@gmail.com                                  ║
╠══════════════════════════════════════════════════════════╣
║  MODEL              │ QUOTA   │ RESET IN   │ STATUS      ║
╠══════════════════════════════════════════════════════════╣
║  Claude 4 Sonnet    │ 85%     │ 4h 23m     │ ✓ OK        ║
║  Claude 4 Opus      │ 100%    │ 5h 0m      │ ✓ OK        ║
║  Gemini 3 Flash     │ 0%      │ 2h 15m     │ ✗ EXHAUSTED ║
║  Gemini 3 Pro       │ 50%     │ 3h 45m     │ ✓ OK        ║
╚══════════════════════════════════════════════════════════╝
```

### JSON Output

```bash
ag-quota --json
```

## Commands

| Command | Description |
|---------|-------------|
| `ag-quota` | Check quota (default) |
| `ag-quota login` | Login with Google |
| `ag-quota status` | Check auth status |
| `ag-quota logout` | Clear stored tokens |
| `ag-quota --help` | Show help |

## Configuration

Config files are stored in:
- **Linux**: `~/.config/ag-quota/`
- **macOS**: `~/Library/Application Support/ag-quota/`

## Requirements

- Go 1.21+ (for building)
- Internet connection (for API calls)
- Google account with Anti-Gravity access

## Credits

Inspired by [antigravity-usage](https://github.com/skainguyen1412/antigravity-usage) by skainguyen1412.

## License

MIT
