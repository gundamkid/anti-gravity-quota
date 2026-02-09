# Changelog

All notable changes to this project will be documented in this file.

## [0.1.4] - 2026-02-09

### 🚀 Features
- **Watch Mode**: Real-time quota monitoring with `ag-quota --watch`.
  - Supports custom refresh intervals (e.g., `--watch=10`).
  - Clears screen and updates display automatically.
  - Handles `Ctrl+C` gracefully.
- **Improved UX**: Redesigned OAuth callback page with "NieR:Automata" theme ("Glory to Mankind") for a premium authentication experience.

### 🛠 Improvements
- **Concurrency**: Refactored `quota --all` to use `errgroup` for better error propagation and context cancellation.
- **Testing**: Improved test stability by mocking time in `TestFormatResetTime`, eliminating flaky tests.
- **CI/CD**: Optimized GitHub Actions pipeline by combining Lint and Test jobs into a single `quality-check` job, reducing execution time.

### 📝 Documentation
- **Retrospective**: Added comprehensive Sprint 2 retrospective covering process, engineering, and action items.

## [0.1.3] - 2026-01-30

### 🚀 Features
- **Account Tiers**: Display account subscription tier (Free 📦, Pro 💎, Ultra 🚀) in `quota`, `status`, and `accounts list` commands.
- **Improved UI**: Refined the `accounts list` table with a clean, rounded style and better indentation for consistent readability.

### 📄 Documentation & Workflow
- **GitFlow with PRs**: Updated development workflow to require Pull Requests for the `dev` branch.
- **Developer Skills**: Established strict verification rules and coding standards.

## [0.1.2] - 2026-01-29

### 🚀 CI/CD & Automation
- **GitHub Actions**: Integrated full CI/CD pipeline for automated testing and multi-platform distribution.
- **Automated Releases**: Configured automated GitHub Releases upon version tagging.
- **Multi-platform Build**: Updated build system to automatically support Linux (amd64/arm64), macOS (amd64/arm64), and Windows (amd64).

### 🛠 Quality & Maintenance
- **Linting**: Added comprehensive `golangci-lint` configuration and resolved multiple quality issues.
- **Go 1.25 Support**: Verified compatibility and optimized build flow for Go 1.25.1.
- **Dynamic Versioning**: Updated `build.sh` to dynamically extract version from `Makefile`.

### 📄 Documentation
- **Build Flow**: Added detailed `docs/build-flow.md` covering local development and CI/CD processes.

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
