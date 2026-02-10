# Conventional Commits & Jira Integration Rule

All commit messages in this project MUST follow the Conventional Commits specification and include the relevant Jira Issue ID. This ensures compatibility with automatic changelog generation and traceability.

## Commit Message Format

```
<type>(<scope>): <description> (AGQ-XXX)
```

- **type**: The category of change (see list below).
- **scope** (optional): The specific part of the codebase being changed (e.g., `api`, `ui`, `auth`).
- **description**: A short, imperative-mood description of the change (e.g., "add", "fix", "change").
- **AGQ-XXX**: The mandatory Jira Issue ID associated with the change.

## Allowed Types

| Type | Purpose |
|------|---------|
| `feat` | A new feature for the user |
| `fix` | A bug fix for the user |
| `docs` | Documentation only changes |
| `style` | Changes that do not affect the meaning of the code (white-space, formatting, etc) |
| `refactor` | A code change that neither fixes a bug nor adds a feature |
| `perf` | A code change that improves performance |
| `test` | Adding missing tests or correcting existing tests |
| `chore` | Changes to the build process or auxiliary tools and libraries |
| `revert` | Reverts a previous commit |

## Examples

- `feat(ui): add compact mode support (AGQ-40)`
- `fix(auth): resolve race condition in token refresh (AGQ-30)`
- `docs: create CONTRIBUTING.md (AGQ-29)`
- `refactor(api): standardize error wrapping (AGQ-28)`

## Why?
- **Auto-Changelog**: `git-chglog` uses these types to categorize changes.
- **Jira Links**: The `(AGQ-XXX)` pattern allows `git-chglog` to generate direct links to Jira tickets.
- **Clarity**: Makes the project history easy to scan and understand.
