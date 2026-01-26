# Anti-Gravity Quota CLI - Task Tracker

## 🎯 Current Sprint

### Phase 2: Authentication Module
- [ ] Implement OAuth2 login flow
  - [ ] Generate PKCE code verifier/challenge
  - [ ] Start local HTTP callback server
  - [ ] Open browser for Google auth
  - [ ] Handle callback and exchange code
- [ ] Implement token storage
  - [ ] Save tokens to config file
  - [ ] Load tokens from config
  - [ ] Auto-refresh expired tokens
- [ ] Test login flow end-to-end

### Phase 3: Cloud Code API Client
- [ ] Implement HTTP client with auth headers
- [ ] Implement `loadCodeAssist` endpoint
- [ ] Implement `fetchAvailableModels` endpoint
- [ ] Handle API errors and rate limiting
- [ ] Parse quota response into models

### Phase 4: CLI Commands
- [ ] Implement `ag-quota` (default quota command)
- [ ] Implement `ag-quota login`
- [ ] Implement `ag-quota status`
- [ ] Implement `ag-quota logout`
- [ ] Add colored output and table formatting

### Phase 5: Polish & Release
- [ ] Add JSON output option (`--json`)
- [ ] Add error handling and user-friendly messages
- [ ] Write usage documentation
- [ ] Test on Linux
- [ ] Build release binaries

---

## ✅ Completed Tasks

### Phase 1: Project Setup & Core Structure ✅
- [x] Initialize Go module
- [x] Set up project directory structure
- [x] Add dependencies (cobra, color, tablewriter, oauth2)
- [x] Create main.go entry point

**Completed:** 2026-01-26

**Summary:**
- Created Go module at `github.com/gundamkid/anti-gravity-quota`
- Set up organized directory structure with `cmd/` and `internal/` packages
- Added all required dependencies: cobra (CLI), color (output), tablewriter (formatting), oauth2 (auth)
- Built working CLI with command structure: quota, login, status, logout
- Verified CLI builds and runs successfully

---

## 📝 Notes

- Focus on Cloud Mode only (không cần Local Mode)
- Single account support cho MVP
- Prioritize simple, working implementation
