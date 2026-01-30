---
description: GitFlow workflow
---

# GitFlow Workflow

This project uses a simplified GitFlow branching strategy with two main branches and feature branches.

## Branch Structure

### Main Branches

1. **`master`** - Production-ready code
   - Always stable and deployable
   - Tagged with version numbers (e.g., `v0.1.0`, `v0.1.1`)
   - Only updated at end of sprint via merge from `dev`

2. **`dev`** - Integration branch
   - Latest development changes
   - Base for all feature branches
   - Merged into `master` at sprint end

### Supporting Branches

3. **`features/<task-name>`** - Feature development
   - Format: `features/AGQ-7-switch-accounts`
   - Created from: `dev`
   - Merged back to: `dev`
   - Deleted after merge

---

## Workflow Steps

### 1. Start New Feature

When starting a new task (e.g., AGQ-7):

```bash
# turbo
git checkout dev
# turbo
git pull origin dev
# turbo
git checkout -b features/AGQ-7-switch-accounts
```

**Rules:**
- Always create feature branch from latest `dev`
- Use task ID in branch name for traceability
- Branch name format: `features/<TASK-ID>-<short-description>`

---

### 2. Work on Feature

During development:

```bash
# Make changes, then commit
# turbo
git add -A
git commit -m "feat(AGQ-7): Add accounts list command"

# Push to remote regularly
git push origin features/AGQ-7-switch-accounts
```

**Commit Message Format:**
```
<type>(<task-id>): <description>

Types: feat, fix, docs, test, refactor, chore
Examples:
  feat(AGQ-7): Add accounts switch command
  fix(AGQ-8): Handle missing account error
  test(AGQ-7): Add unit tests for accounts display
  docs(AGQ-7): Update README with accounts commands
```

---

### 3. Complete Feature

When feature is done and tested:

```bash
# Ensure all tests pass
# turbo
make test

# Build to verify
# turbo
make build

# Checkout dev and update
# turbo
git checkout dev
# turbo
git pull origin dev

# Merge feature branch
# turbo
git merge features/AGQ-7-switch-accounts --no-edit

# Push to remote
git push origin dev

# Delete feature branch (optional)
git branch -d features/AGQ-7-switch-accounts
git push origin --delete features/AGQ-7-switch-accounts
```

**Checklist before merge:**
- [ ] All tests passing (`make test`)
- [ ] Code builds successfully (`make build`)
- [ ] Jira task moved to "In Review" or "Done"
- [ ] No merge conflicts with `dev`

---

### 4. End of Sprint Release

At sprint end, merge `dev` to `master` and create release:

```bash
# Checkout master and update
git checkout master
git pull origin master

# Merge dev into master
git merge dev --no-edit

# Update version in files
# Edit: Makefile (VERSION), go files if needed

# Update CHANGELOG.md
# Add new version section with changes

# Commit version bump
git add -A
git commit -m "chore: Bump version to 0.1.1"

# Create and push tag
git tag -a v0.1.1 -m "Release version 0.1.1"
git push origin master
git push origin v0.1.1

# Update dev from master
git checkout dev
git merge master --no-edit
git push origin dev
```

**Release Checklist:**
- [ ] All sprint tasks completed and merged to `dev`
- [ ] All tests passing on `dev`
- [ ] CHANGELOG.md updated with new version
- [ ] Version bumped in Makefile and relevant files
- [ ] Tag created with version number
- [ ] Both `master` and `dev` pushed to remote

---

## Branch Protection Rules (Recommended)

For GitHub repository settings:

### `master` branch:
- Require pull request reviews before merging
- Require status checks to pass (CI/CD if available)
- Require branches to be up to date before merging
- Do not allow force pushes

### `dev` branch:
- Require status checks to pass
- Allow direct pushes (for quick integration)
- Do not allow force pushes

---

## Quick Reference

| Action | Command |
|--------|---------|
| Start feature | `git checkout dev && git pull && git checkout -b features/TASK-ID-name` |
| Commit changes | `git commit -m "feat(TASK-ID): description"` |
| Finish feature | `git checkout dev && git merge features/TASK-ID-name` |
| Create release | `git checkout master && git merge dev && git tag v0.1.1` |

---

## Hotfix Workflow (Emergency)

For critical bugs in production:

```bash
# Create hotfix from master
git checkout master
git pull origin master
git checkout -b hotfix/critical-bug-fix

# Fix, test, commit
git commit -m "fix: Critical bug description"

# Merge to master
git checkout master
git merge hotfix/critical-bug-fix --no-edit
git tag -a v0.1.2 -m "Hotfix release 0.1.2"
git push origin master
git push origin v0.1.2

# Also merge to dev
git checkout dev
git merge hotfix/critical-bug-fix --no-edit
git push origin dev

# Delete hotfix branch
git branch -d hotfix/critical-bug-fix
```

---

## Current State

**Active Branches:**
- `master` - Latest stable release
- `dev` - Current development (to be created)
- `release/0.1.1` - Temporary release branch (to be merged and removed)

**Migration Steps:**
1. Create `dev` branch from `release/0.1.1`
2. Merge `release/0.1.1` to `master`
3. Delete `release/0.1.1` branch
4. Start using new GitFlow workflow