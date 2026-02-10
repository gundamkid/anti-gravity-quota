# 🚀 AG-Quota

A lightweight CLI tool to monitor your Anti-Gravity (Claude Code) AI model quota and usage in real-time.

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/gundamkid/anti-gravity-quota/blob/master/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/doc/devel/release.html)

## ✨ Features

- ✅ **Quota Monitoring** - Real-time tracking for Claude 3.5/4 and Gemini 1.5/2.5/3 models.
- 👤 **Multi-Account Support** - Manage multiple Google accounts with seamless switching.
- 🔐 **Secure OAuth2** - Enterprise-grade authentication using PKCE flow.
- 📊 **Visual Dashboard** - Beautiful terminal tables with status indicators.
- 🔄 **Real-time Notifications** - Telegram alerts when quotas change or hit critical levels.
- ⏳ **Watch Mode** - Passive monitoring with configurable refresh intervals.
- 📉 **Display Modes** - Automatic compact mode for small terminals or forced modes via flags.
- 📝 **Developer Friendly** - JSON output for easy integration and automation.

## 🚀 Installation

### Using Go
```bash
go install github.com/gundamkid/anti-gravity-quota/cmd/ag-quota@latest
```

### From Source
```bash
git clone https://github.com/gundamkid/anti-gravity-quota.git
cd anti-gravity-quota
make build
# Binary will be available as ./ag-quota
```

## 🚥 Quick Start

1. **Login** to your account:
   ```bash
   ag-quota login
   ```
2. **Check quota** immediately:
   ```bash
   ag-quota
   ```
3. **Enable Notifications** (Telegram):
   ```bash
   ag-quota config set-telegram --token "BOT_TOKEN" --chat-id "CHAT_ID"
   ```
4. **Monitor all accounts** in real-time:
   ```bash
   ag-quota quota --all --watch=5
   ```

---

## 📖 Usage Guide

### 1. Checking Quota

Display quota information for models. The tool detects your terminal size and chooses the best display mode automatically.

```bash
# Standard view (default account)
$ ag-quota

# View specific account
$ ag-quota quota --account user@gmail.com

# Aggregate view for ALL accounts
$ ag-quota quota --all

# Force display modes
$ ag-quota --compact      # Force minimal table
$ ag-quota --no-compact   # Force full detailed table
```

**Status Indicators:**
- ✅ **HEALTHY** - Above 50% remaining.
- ⚠️ **WARNING** - 21% to 50% remaining.
- ⛔ **CRITICAL** - 1% to 20% remaining.
- ❌ **EMPTY** - Quota exhausted.

### 2. Account Management

Securely manage multiple Google sessions.

```bash
# List all saved accounts
ag-quota accounts list

# Switch default account
ag-quota accounts switch user@gmail.com

# Remove an account session
ag-quota accounts remove old@user.com
```

### 3. Watch Mode & Notifications

Stay updated without manual refreshes.

```bash
# Watch with default 5-minute interval
ag-quota --watch

# Custom interval (e.g., every 2 minutes)
# Note: Use '=' for custom values (e.g., --watch=2)
ag-quota --watch=2
```

> [!TIP]
> **Telegram Setup**: For step-by-step instructions on setting up your notification bot, see the [Telegram Setup Guide](docs/telegram-setup.md).

### 4. Configuration & Testing

```bash
# Test your notification settings with dummy data
ag-quota config test-notify

# View current Telegram status
ag-quota config get-telegram
```

---

## 🛠️ Integration

### JSON Output
Perfect for custom scripts, status bars, or automation.

```bash
$ ag-quota --json | jq '.Models[] | select(.IsExhausted == true)'
```

---

## 📁 Technical Overview

- **Storage**: Auth tokens and config are stored in `~/.config/ag-quota/` (Linux/macOS) with `0600` permissions.
- **Auto-Refresh**: Tokens are automatically refreshed before expiration.
- **Retry Logic**: Built-in exponential backoff for API resilience.
- **Documentation**: 
  - [Technical Details](docs/technical.md)
  - [Build & CI/CD Flow](docs/build-flow.md)

---

## 🤝 Contributing

We welcome contributions! Please check the [Contributing Guide](CONTRIBUTING.md) to get started.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

**Built with ❤️ for the Anti-Gravity Community.**
