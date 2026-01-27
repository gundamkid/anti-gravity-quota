# Anti-Gravity Quota CLI - Technical Documentation

## 📡 API Specification

### Base URLs

| Environment | URL |
|-------------|-----|
| Production | `https://cloudcode-pa.googleapis.com` |
| Sandbox | `https://daily-cloudcode-pa.sandbox.googleapis.com` |

### Authentication

**OAuth2 Configuration:**
```
Authorization URL: https://accounts.google.com/o/oauth2/v2/auth
Token URL: https://oauth2.googleapis.com/token
Client ID: 764086051850-6qr4p6gpi6hn506pt8ejuq83di341hur.apps.googleusercontent.com
Redirect URI: http://localhost:8085/callback
Scopes: openid email profile https://www.googleapis.com/auth/cloud-platform
```

**Request Headers:**
```http
Authorization: Bearer <access_token>
Content-Type: application/json
User-Agent: antigravity
```

---

## 📋 API Endpoints

### 1. Load Code Assist

Lấy thông tin project và trạng thái code assist.

**Endpoint:** `POST /v1internal:loadCodeAssist`

**Request Body:**
```json
{
  "metadata": {
    "ideType": "ANTIGRAVITY",
    "platform": "PLATFORM_UNSPECIFIED",
    "pluginType": "GEMINI"
  }
}
```

**Response:**
```json
{
  "codeAssistEnabled": true,
  "planInfo": {
    "monthlyPromptCredits": 1000,
    "planType": "FREE"
  },
  "availablePromptCredits": 850,
  "cloudaicompanionProject": "projects/123456789",
  "currentTier": {
    "id": "free-tier",
    "name": "Free",
    "description": "Free tier"
  },
  "allowedTiers": [
    {"id": "free-tier", "isDefault": true}
  ]
}
```

---

### 2. Fetch Available Models

Lấy danh sách models với thông tin quota.

**Endpoint:** `POST /v1internal:fetchAvailableModels`

**Request Body:**
```json
{
  "project": "projects/123456789"
}
```

**Response:**
```json
{
  "models": {
    "claude-sonnet-4-5": {
      "displayName": "Claude 4 Sonnet",
      "model": "claude-sonnet-4-5",
      "label": "Claude 4 Sonnet",
      "quotaInfo": {
        "remainingFraction": 0.85,
        "resetTime": "2026-01-26T12:00:00Z",
        "isExhausted": false
      },
      "maxTokens": 64000,
      "recommended": true,
      "supportsImages": true,
      "supportsThinking": false,
      "modelProvider": "claude"
    },
    "gemini-3-flash": {
      "displayName": "Gemini 3 Flash",
      "model": "gemini-3-flash",
      "label": "Gemini 3 Flash",
      "quotaInfo": {
        "remainingFraction": 1.0,
        "resetTime": "2026-01-26T14:00:00Z",
        "isExhausted": false
      },
      "modelProvider": "google"
    }
  },
  "defaultAgentModelId": "claude-sonnet-4-5"
}
```

---

## 🔐 OAuth2 Flow

### PKCE Flow (Proof Key for Code Exchange)

```
┌─────────┐                              ┌─────────────┐
│   CLI   │                              │   Google    │
└────┬────┘                              └──────┬──────┘
     │                                          │
     │ 1. Generate code_verifier (random)       │
     │ 2. Generate code_challenge = SHA256(verifier)
     │                                          │
     │ 3. Open browser ────────────────────────>│
     │    ?client_id=...                        │
     │    &redirect_uri=localhost:8085          │
     │    &code_challenge=...                   │
     │    &code_challenge_method=S256           │
     │    &scope=openid email profile...        │
     │    &state=random                         │
     │                                          │
     │ 4. User logs in & consents               │
     │                                          │
     │<──────────── 5. Redirect to callback ────│
     │    ?code=authorization_code              │
     │    &state=random                         │
     │                                          │
     │ 6. Exchange code for tokens ────────────>│
     │    POST /token                           │
     │    code=...                              │
     │    code_verifier=...                     │
     │                                          │
     │<──────────── 7. Return tokens ───────────│
     │    access_token                          │
     │    refresh_token                         │
     │    expires_in                            │
     │                                          │
```

### Token Refresh

Khi access_token hết hạn (thường sau 1 giờ):

```http
POST https://oauth2.googleapis.com/token
Content-Type: application/x-www-form-urlencoded

grant_type=refresh_token
refresh_token=<refresh_token>
client_id=<client_id>
```

---

## 📁 Data Structures

### Account Tokens (`~/.config/ag-quota/accounts/{email}.json`)

```json
{
  "access_token": "ya29.xxx...",
  "refresh_token": "1//xxx...",
  "token_type": "Bearer",
  "expiry": "2026-01-26T08:00:00Z",
  "email": "user@gmail.com"
}
```

### Application Config (`~/.config/ag-quota/config.json`)

```json
{
  "default_account": "user@gmail.com"
}
```

### Model Quota Info (Internal)

```go
type ModelQuota struct {
    ModelID           string    // "claude-sonnet-4-5"
    DisplayName       string    // "Claude 4 Sonnet"
    Label             string    // "Claude 4 Sonnet"
    Provider          string    // "claude" | "google"
    RemainingFraction float64   // 0.0 - 1.0
    ResetTime         time.Time // When quota resets
    IsExhausted       bool      // true if quota = 0
}
```

---

## ⚠️ Error Handling

### HTTP Status Codes

| Code | Meaning | Action |
|------|---------|--------|
| 401 | Unauthorized | Token expired, refresh or re-login |
| 403 | Forbidden | Invalid token, re-login required |
| 429 | Rate Limited | Retry with exponential backoff |
| 5xx | Server Error | Retry after delay |

### Retry Strategy

```go
// Exponential backoff với jitter
func getBackoffDelay(attempt int) time.Duration {
    base := 500 * math.Pow(2, float64(attempt-1))
    jitter := rand.Float64() * 100
    delay := math.Min(base+jitter, 4000)
    return time.Duration(delay) * time.Millisecond
}

// Max 3 attempts per request
const MaxRetryAttempts = 3
```

---

## 🔧 Configuration

### Config Directory

| OS | Path |
|----|------|
| Linux | `~/.config/ag-quota/` |
| macOS | `~/Library/Application Support/ag-quota/` |
| Windows | `%APPDATA%\ag-quota\` |

### Files

| File | Description |
|------|-------------|
| `config.json` | Global application configuration (default account, etc.) |
| `accounts/` | Directory containing OAuth tokens per email (chmod 600) |
| `token.json` | (Deprecated) Old token storage |

---

## 🛠️ Account Management

CLI hỗ trợ quản lý nhiều tài khoản Google cùng lúc thông qua các lệnh sub-command của `accounts`:

| Command | Description |
|---------|-------------|
| `ag-quota accounts list` | Liệt kê tất cả tài khoản đã lưu và trạng thái token |
| `ag-quota accounts default <email>` | Thiết lập tài khoản mặc định cho các lệnh tiếp theo |
| `ag-quota accounts switch <email>` | Alias của lệnh `default` giúp chuyển nhanh giữa các account |

---

## 🧪 Testing

### Automated Testing

Dự án sử dụng bộ test chuẩn của Go để đảm bảo tính ổn định của các module cốt lõi.

#### 1. Chạy toàn bộ tests
Sử dụng Makefile để chạy toàn bộ tests một cách nhanh chóng:
```bash
make test
```
Hoặc dùng lệnh Go trực tiếp:
```bash
go test -v ./...
```

#### 2. Chiến lược Testing
- **Unit Tests**: Kiểm tra logic của các hàm xử lý dữ liệu (`internal/models`, `internal/ui`).
- **Mocking**: Sử dụng `net/http/httptest` để giả lập API của Google Cloud Code (`internal/api`).
- **Browser-less Auth**: Kiểm tra logic trao đổi token và PKCE mà không cần mở trình duyệt thật (`internal/auth`).

#### 3. Test Coverage
Để kiểm tra độ bao phủ của code:
```bash
go test -cover ./...
```

### Manual API Test

```bash
# Get access token first, then:
curl -X POST https://cloudcode-pa.googleapis.com/v1internal:loadCodeAssist \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"metadata":{"ideType":"ANTIGRAVITY","platform":"PLATFORM_UNSPECIFIED","pluginType":"GEMINI"}}'
```

### Expected Models

Based on antigravity-usage source:
- `claude-sonnet-4-5` → Claude family
- `gemini-3-flash` → Gemini flash quota group
- `gemini-3-pro-low` → Gemini pro quota group
