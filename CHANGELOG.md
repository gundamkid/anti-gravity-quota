# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2026-01-26

### 🚀 Features

- **Initial MVP Release**: Launched the first version of Anti-Gravity Quota CLI.
- **Google Login**: Implemented single-account Google OAuth2 login flow with PKCE security.
- **Quota Monitoring**: Ability to view real-time quota usage for all Anti-Gravity models (Claude 4, Gemini 3, etc.).
- **JSON Output**: Added `--json` flag to provide machine-readable output for scripts and automation.
- **Multi-OS Support**: Included `build.sh` script to easily build binaries for Linux, macOS, and Windows.

### 🛠 Improvements

- **Pretty Output**: Enhanced terminal output with colored tables and visual status indicators using `go-pretty`.
- **Automatic Refresh**: Implemented automatic token refreshing to handle expiration gracefully.
