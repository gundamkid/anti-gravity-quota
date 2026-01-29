# 🏗️ Build & Release Flow

This document describes the build, continuous integration, and release process for the `ag-quota` CLI tool.

---

## 🛡️ Contribution & Branch Rules (Strict)

To ensure code quality in this Open Source project, we enforce strict **Branch Protection Rules** on GitHub.

### 1. Protected Branches
*   **`master`**: Production-ready code. No direct pushes allowed.
*   **`dev`**: Main development branch. No direct pushes allowed.

### 2. Contribution Workflow
1.  **Fork & Branch**: Create a feature branch from `dev` (e.g., `features/my-new-feature`).
2.  **Pull Request (PR)**: Submit a PR to merge your branch into `dev`.
3.  **Automated Checks**: The following CI/CD checks **MUST PASS** before merging:
    *   ✅ **Test**: All unit tests passing.
    *   ✅ **Lint**: Code style meets strict `golangci-lint` standards.
    *   ✅ **Build**: Code compiles successfully on all platforms.
4.  **Review**: At least 1 approval is required from maintainers.

### 3. Merging
*   Once checks pass and code is approved, squash and merge into `dev`.
*   Maintainers will periodically merge `dev` into `master` for releases.

---

## 💻 Local Development

For local development and testing, use the provided `Makefile`.

### 1. Requirements
- Go 1.25.1+
- `golangci-lint` (for local linting)
- `make`

### 2. Common Commands
| Command | Description |
|---------|-------------|
| `make build` | Builds the binary for the current platform in the root directory. |
| `make test` | Runs all unit tests. |
| `make fmt` | Formats the codebase using `go fmt`. |
| `make lint` | Runs `golangci-lint` (must be installed). |
| `make clean` | Removes build artifacts and the `dist/` directory. |

### 3. All-Platform Build
To build for all supported platforms (Linux, macOS, Windows) locally:
```bash
make release
```
This executes `./build.sh`, which generates binaries in the `dist/` folder.

---

## 🚀 CI/CD Pipeline (GitHub Actions)

The project uses GitHub Actions to automate testing and distribution. The workflow is defined in `.github/workflows/build.yml`.

### 1. Workflow Triggers
- **Push to `master`**: Triggers full test, lint, and build. Builds are uploaded as artifacts.
- **Pull Request to `master`**: Triggers test and lint to ensure quality before merging.
- **Git Tags (`v*`)**: Triggers the **Release** flow.
- **Manual Trigger**: Use the "Run workflow" button in the GitHub Actions tab to test any branch.

### 2. Jobs Flow (Sequential)
1.  **Lint**:
    *   Runs first to ensure code quality and style.
    *   Must pass strict `golangci-lint` checks.
2.  **Test**:
    *   Runs **only after Lint passes**.
    *   Executes unit tests (`go test -v -race ./...`).
    *   **Result**: This is the primary Status Check for PRs.
3.  **Build & Release** (Tag Only):
    *   Runs **only after Test passes** AND **only on Tags (`v*`)**.
    *   Compiles binaries for all platforms (Linux, macOS, Windows).
    *   Creates GitHub Release and uploads artifacts.

---

## 📦 Release Process

To publish a new version of `ag-quota`:

1. **Prepare the release**:
   - Update the `VERSION` variable in `Makefile`.
   - Update `CHANGELOG.md` with new changes.
2. **Merge to master**:
   - Ensure all features are merged into `master`.
3. **Create and push a tag**:
   ```bash
   git checkout master
   git pull origin master
   git tag v0.1.2
   git push origin v0.1.2
   ```
4. **Automated Release**:
   - GitHub Actions will detect the tag.
   - It will build all binaries and create a release on the GitHub project page.
   - Users can download the binaries directly from the "Releases" section.

---

## 🛠️ Build Matrix

| OS | Arch | Binary Name |
|----|------|-------------|
| Linux | amd64 | `ag-quota-linux-amd64` |
| Linux | arm64 | `ag-quota-linux-arm64` |
| Darwin (macOS) | amd64 | `ag-quota-darwin-amd64` |
| Darwin (macOS) | arm64 | `ag-quota-darwin-arm64` |
| Windows | amd64 | `ag-quota-windows-amd64.exe` |
