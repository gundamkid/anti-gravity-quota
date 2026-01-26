# Anti-Gravity Quota CLI - Task Tracker

## 🎯 Current Sprint

### Phase 1: Project Setup & Core Structure
- [ ] Initialize Go module
- [ ] Set up project directory structure
- [ ] Add dependencies (cobra, color, tablewriter, oauth2)
- [ ] Create main.go entry point

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

_None yet_

---

## 📝 Notes

- Focus on Cloud Mode only (không cần Local Mode)
- Single account support cho MVP
- Prioritize simple, working implementation
