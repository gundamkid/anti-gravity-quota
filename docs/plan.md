# Anti-Gravity Quota CLI - Implementation Plan

## 📋 Project Overview

Xây dựng CLI tool bằng Go để kiểm tra quota của các model Anti-Gravity, hỗ trợ đăng nhập qua Cloud Mode sử dụng Google Cloud Code API.

---

## 🔍 Reverse Engineering Summary

Dựa trên phân tích repo [antigravity-usage](https://github.com/skainguyen1412/antigravity-usage):

### Architecture
```
┌─────────────────────────────────────────────────────────┐
│                    CLI Application                       │
├─────────────────────────────────────────────────────────┤
│  Commands: quota | login | accounts | status            │
├─────────────────────────────────────────────────────────┤
│                    Core Services                         │
├──────────────────────┬──────────────────────────────────┤
│    Cloud Client      │       Local Client               │
│  (Google Cloud API)  │   (IDE Language Server)          │
├──────────────────────┴──────────────────────────────────┤
│                   Token Manager                          │
│              (OAuth2 + Token Storage)                    │
└─────────────────────────────────────────────────────────┘
```

### Key Findings

#### 1. Cloud Mode (Primary Focus)
- **Base URL**: `https://cloudcode-pa.googleapis.com`
- **Backup URL**: `https://daily-cloudcode-pa.sandbox.googleapis.com`
- **Authentication**: Google OAuth2

**Key API Endpoints:**
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1internal:loadCodeAssist` | POST | Load code assist status, get project ID |
| `/v1internal:fetchAvailableModels` | POST | Fetch available models with quota info |
| `/v1internal:onboardUser` | POST | Onboard new users |

**Request Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
User-Agent: antigravity
```

**Metadata Structure:**
```json
{
  "metadata": {
    "ideType": "ANTIGRAVITY",
    "platform": "PLATFORM_UNSPECIFIED",
    "pluginType": "GEMINI"
  }
}
```

#### 2. Model Quota Response Structure
```json
{
  "models": {
    "modelId": {
      "displayName": "Claude 4 Sonnet",
      "model": "claude-sonnet-4-5",
      "label": "Claude 4 Sonnet",
      "quotaInfo": {
        "remainingFraction": 0.85,
        "resetTime": "2026-01-26T12:00:00Z",
        "isExhausted": false
      },
      "modelProvider": "claude"
    }
  },
  "defaultAgentModelId": "claude-sonnet-4-5"
}
```

#### 3. OAuth2 Configuration
- **Scopes**: `openid email profile https://www.googleapis.com/auth/cloud-platform`
- **Token Storage**: File-based in `~/.antigravity/` or XDG config dir
- **Token Refresh**: Automatic refresh when expired

---

## 🏗️ Implementation Plan

### Phase 1: Project Setup & Core Structure

#### 1.1 Project Structure
```
anti-gravity-quota/
├── cmd/
│   └── ag-quota/
│       └── main.go              # CLI entry point
├── internal/
│   ├── auth/
│   │   ├── oauth.go             # OAuth2 flow
│   │   └── token.go             # Token management
│   ├── api/
│   │   └── cloudcode.go         # Cloud Code API client
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── models/
│   │   └── quota.go             # Data models
│   └── ui/
│       └── display.go           # Terminal UI/Output
├── go.mod
├── go.sum
└── README.md
```

#### 1.2 Dependencies
```go
// go.mod
module github.com/gundamkid/anti-gravity-quota

require (
    github.com/spf13/cobra v1.8.0       // CLI framework
    github.com/fatih/color v1.16.0      // Colored output
    github.com/jedib0t/go-pretty/v6 v6.5.4 // Table display
    golang.org/x/oauth2 v0.16.0         // OAuth2 client
)
```

---

### Phase 2: Authentication Module

#### 2.1 OAuth2 Implementation

```go
// internal/auth/oauth.go

const (
    GoogleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
    GoogleTokenURL = "https://oauth2.googleapis.com/token"
    
    // Cloud Code OAuth client (from antigravity extension)
    clientID     = "764086051850-6qr4p6gpi6hn506pt8ejuq83di341hur.apps.googleusercontent.com"
    scopes = "openid email profile https://www.googleapis.com/auth/cloud-platform"
)

func Login() error {
    // 1. Generate state & code verifier (PKCE)
    // 2. Start local HTTP server for callback
    // 3. Open browser to Google auth URL
    // 4. Handle callback, exchange code for tokens
    // 5. Save tokens to config file
}
```

#### 2.2 Token Storage

```go
// internal/auth/token.go

type TokenData struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
    Email        string    `json:"email"`
}

func SaveToken(token *TokenData) error
func LoadToken() (*TokenData, error)
func RefreshToken(token *TokenData) (*TokenData, error)
func GetValidToken() (string, error)  // Auto-refresh if expired
```

---

### Phase 3: Cloud Code API Client

#### 3.1 API Client

```go
// internal/api/cloudcode.go

const (
    BaseURL = "https://cloudcode-pa.googleapis.com"
)

type CloudCodeClient struct {
    httpClient *http.Client
    token      string
    projectID  string
}

// Load code assist and get project ID
func (c *CloudCodeClient) LoadCodeAssist() (*LoadCodeAssistResponse, error) {
    body := map[string]interface{}{
        "metadata": map[string]string{
            "ideType":    "ANTIGRAVITY",
            "platform":   "PLATFORM_UNSPECIFIED",
            "pluginType": "GEMINI",
        },
    }
    return c.post("/v1internal:loadCodeAssist", body)
}

// Fetch available models with quota
func (c *CloudCodeClient) FetchAvailableModels() (*ModelsResponse, error) {
    body := map[string]interface{}{}
    if c.projectID != "" {
        body["project"] = c.projectID
    }
    return c.post("/v1internal:fetchAvailableModels", body)
}
```

#### 3.2 Response Models

```go
// internal/models/quota.go

type ModelQuota struct {
    ModelID           string
    DisplayName       string
    Provider          string
    RemainingFraction float64
    ResetTime         time.Time
    IsExhausted       bool
}

type QuotaInfo struct {
    Email       string
    ProjectID   string
    Models      []ModelQuota
    FetchedAt   time.Time
}
```

---

### Phase 4: CLI Commands

#### 4.1 Command Structure

```go
// cmd/ag-quota/main.go

var rootCmd = &cobra.Command{
    Use:   "ag-quota",
    Short: "Anti-Gravity Quota CLI",
}

// ag-quota quota - Check quota for all models
var quotaCmd = &cobra.Command{
    Use:   "quota",
    Short: "Check quota for all models",
    Run:   runQuota,
}

// ag-quota login - Login with Google
var loginCmd = &cobra.Command{
    Use:   "login",
    Short: "Login with Google account",
    Run:   runLogin,
}

// ag-quota status - Check auth status
var statusCmd = &cobra.Command{
    Use:   "status",
    Short: "Check authentication status",
    Run:   runStatus,
}
```

#### 4.2 Output Format

```
╔══════════════════════════════════════════════════════════════╗
║              Anti-Gravity Quota Status                        ║
╠══════════════════════════════════════════════════════════════╣
║  Account: user@gmail.com                                      ║
║  Fetched: 2026-01-26 07:12:48 UTC                            ║
╠══════════════════════════════════════════════════════════════╣
║  MODEL                 │ QUOTA      │ RESET IN    │ STATUS   ║
╠═════════════════════════════════════════════════════════════╣
║  Claude 4 Sonnet       │ 85%        │ 4h 23m      │ ✓ OK     ║
║  Claude 4 Opus         │ 100%       │ 5h 0m       │ ✓ OK     ║
║  Gemini 3 Flash        │ 0%         │ 2h 15m      │ ✗ EMPTY  ║
║  Gemini 3 Pro          │ 50%        │ 3h 45m      │ ✓ OK     ║
╚══════════════════════════════════════════════════════════════╝
```

---

### Phase 5: Configuration & Storage

```go
// internal/config/config.go

// Config directory: ~/.config/ag-quota/ (Linux/Mac)
func GetConfigDir() string {
    if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
        return filepath.Join(xdg, "ag-quota")
    }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".config", "ag-quota")
}

// Files:
// - ~/.config/ag-quota/token.json    - OAuth tokens
// - ~/.config/ag-quota/config.json   - User preferences
```

---

## ✅ Implementation Checklist

### Core Features (MVP)
- [x] OAuth2 login flow with Google
- [x] Token storage and refresh
- [x] Cloud Code API client
- [x] Quota fetching for all models
- [x] Pretty table output in terminal

### Commands
- [x] `ag-quota` / `ag-quota quota` - Show all model quotas
- [x] `ag-quota login` - Start OAuth login flow
- [x] `ag-quota status` - Show auth status
- [x] `ag-quota logout` - Clear stored tokens

### Nice to Have (Future)
- [ ] Multi-account support
- [x] JSON output option (`--json`)
- [ ] Watch mode (`--watch`)
- [ ] Cache support for offline viewing

---

## 🔧 Build & Run

```bash
# Build
go build -o ag-quota ./cmd/ag-quota

# Or run directly
go run ./cmd/ag-quota

# Install globally
go install ./cmd/ag-quota
```

---

## 📝 Notes

1. **Client ID**: Using Google Cloud Code extension's public OAuth client ID
2. **No API Key Required**: OAuth2 tokens are sufficient for authentication
3. **Rate Limiting**: API may rate limit requests; implement backoff
4. **Token Security**: Store tokens with appropriate file permissions (0600)

---

## 🎯 Next Steps

1. Bạn review plan này và cho feedback
2. Sau khi approve, tôi sẽ bắt đầu implement từ Phase 1
3. Mỗi phase sẽ được test trước khi chuyển sang phase tiếp theo
