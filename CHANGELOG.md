## v0.1.5 (2026-02-10)

### ✨ Features

- refine notification messages with deltas and reset times (008b172)
- integrate notifications into watch mode (f932b4f)
- implement state tracker and message formatter for notifications (7986032)

### 🐛 Bug Fixes

- install git-chglog and fix configuration for changelog generation (8bc504d)
- resolve lint errors and update migration tests (74bde25)

### 📚 Documentation

- update README and Telegram setup guide (AGQ-39) (4f63214)
- create CONTRIBUTING.md (AGQ-29) (2022d0f)

### 🔨 Chore

- bump version to v0.1.5 and update changelog configuration (1b54eb4)
- ignore .agent directory and remove from index (AGQ-29) (5abd374)
- extract conventional-commits rule to shared file (a90d140)
- Bump version to 0.1.4 in Makefile (9fe6923)

## v0.1.4 (2026-02-09)

### 🔨 Chore

- Bump version to 0.1.4 (2de5003)
- stop tracking .agent directory (ca2a3ed)

## v0.1.3 (2026-01-30)

### 🐛 Bug Fixes

- check error return values for token saving functions (29e24e4)

### 📚 Documentation

- add PR creation rules to workflow (e3688c4)
- update workflow to require PR for dev branch (550e7a6)

### 🔨 Chore

- bump version to 0.1.3 and update changelog (07c7d0e)
- pin golangci-lint version to v1.64.8 for CI stability (0cee44c)
- update .gitignore to ignore .agent directory (d289b7d)

## v0.1.2 (2026-01-29)

### 📚 Documentation

- update ci/cd flow description to reflect strict lint->test->build sequence (88186e5)
- link build-flow documentation in README and technical docs (2c14597)
- add build and release flow documentation (2620c68)

### 🔨 Chore

- enable ci trigger for PRs targeting dev branch (3500448)
- add dry build verification step to test job (c236e6a)
- refactor ci/cd pipeline to enforce lint -> test -> build/release flow (f648a94)
- correct version to 0.1.2 and refresh release (9f1356e)
- bump version to 0.1.3 and update changelog (2aeff94)

## v0.1.1 (2026-01-28)

### ✨ Features

- Release v0.1.1 with multi-account support and bug fixes, and add new agent skills and GitFlow workflow documentation. (0a1b944)
- Implement multi-account migration, parallelize quota fetching, and add an `accounts remove` command. (91e729b)
- Implement multi-account management with refactored token storage, new PKCE support, and comprehensive tests across authentication, API, config, models, and UI. (644ec23)

### 📚 Documentation

- Add automated testing instructions to README and detailed technical documentation on testing strategy and coverage. (58cd1cc)

### 🔨 Chore

- Bump version to 0.1.1 (666562e)

## v0.1.0 (2026-01-26)

### ✨ Features

- Update Google OAuth client ID and explicitly set redirect URL in the configuration. (a8c49d0)
- Add JSON output option and a multi-platform build system with `build.sh` and `Makefile`. (2fbc7c6)
- Implement `ag-quota` command to display quota information using a new UI module with colored and formatted output. (0f0e9fd)
- implement Cloud Code API client for fetching model quota information, including authentication, retry logic, and defining new data models. (b6351b1)
- Implement OAuth2 authentication with PKCE, token management, and CLI commands for login, status, and logout. (24234dd)
- Initialize Go module, add core dependencies, and configure project setup with a new gitignore. (fe91dfa)
- Initialize Go module, establish basic Cobra CLI with core commands, add dependencies, and configure gitignore. (b591dac)
- Implement CSRF token handling in the client, enhance language server process detection, and outline future cloud mode tasks. (87cc1e5)
- Implement CLI with quota display, configuration management, and mock servers, updating task progress. (eacbbe5)
- implement GetUserStatus API client and update task tracking documentation. (ed0a6e3)
- Initialize Go module and implement Antigravity Language Server process and port detection. (1b02490)
- add initial project documentation including technical design, implementation tasks, and user guide. (b8edf74)

### 📚 Documentation

- add CHANGELOG for v0.1.0 (542b314)
- Refine project plan. (d2c84c3)
- Remove port 8085 requirement and troubleshooting instructions from README. (0eef64f)

### 🔧 Improvements

- remove account and project information display and filter out models without a display name. (53d63eb)
- Migrate table rendering from `tablewriter` to `go-pretty`, updating UI examples and documentation. (f6c6442)

### 🔨 Chore

- Add .claude/ to .gitignore. (8902a96)
- remove compiled mock binary from git tracking (4b6ffef)

