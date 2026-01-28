# Changelog

All notable changes to this project will be documented in this file.

## [0.1.1] - 2026-01-28

### 🚀 Features

- **Multiple Account Support**: Manage multiple Google accounts simultaneously.
  - `ag-quota accounts list`: View all saved accounts.
  - `ag-quota accounts switch <email>`: Quickly switch between accounts.
  - `ag-quota accounts remove <email>`: Delete specific accounts from storage.
- **Batched Quota Check**: Added `ag-quota quota --all` to check quotas for all saved accounts in parallel.
- **Targeted Quota**: Added `--account` flag to check quota for a specific email without switching context.

### 🐛 Bug Fixes

- **Login Flow**: Fixed issue where users were not correctly recognized as logged in immediately after authentication.
- **Race Condition**: Fixed critical bug in `quota --all` where concurrent checks caused token conflicts (AGQ-12).

### 🛠 Improvements

- **Performance**: Optimized multi-account quota fetching with concurrent goroutines.
- **Documentation**: Updated README with comprehensive account management guides.

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
