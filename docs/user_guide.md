# Hướng dẫn Sử dụng Anti-Gravity Quota

## Giới thiệu

**Anti-Gravity Quota** là công cụ dòng lệnh giúp bạn theo dõi quota sử dụng các AI models trong Antigravity IDE một cách nhanh chóng và tiện lợi.

---

## Cài đặt

### Cách 1: Tải binary có sẵn

```bash
# Linux
curl -L https://github.com/your-username/anti-gravity-quota/releases/latest/download/anti-gravity-quota-linux-amd64 -o anti-gravity-quota
chmod +x anti-gravity-quota
sudo mv anti-gravity-quota /usr/local/bin/
```

### Cách 2: Build từ source

```bash
git clone https://github.com/your-username/anti-gravity-quota.git
cd anti-gravity-quota
go build -o anti-gravity-quota ./cmd/anti-gravity-quota
sudo mv anti-gravity-quota /usr/local/bin/
```

### Xác nhận cài đặt

```bash
anti-gravity-quota --version
# Output: anti-gravity-quota version 1.0.0
```

---

## Sử dụng cơ bản

### Xem quota hiện tại

Chỉ cần gõ `anti-gravity-quota` trong terminal (yêu cầu Antigravity IDE đang chạy):

```bash
anti-gravity-quota
```

**Output:**
```
╭─────────────────────────────────────────────────────────────╮
│  Anti-Gravity Quota Monitor                                 │
│  Account: ngoanhtuan245@gmail.com                           │
├─────────────────────────────────────────────────────────────┤
│  MODEL           │  USED/LIMIT    │  REMAINING  │  RESETS   │
├──────────────────┼────────────────┼─────────────┼───────────┤
│  Gemini 3 Pro    │  1,500/5,000   │   70% ████░ │  4h 32m   │
│  Claude 4.5      │    200/1,000   │   80% ████░ │  4h 32m   │
│  GPT-OSS 120B    │      0/2,000   │  100% █████ │  4h 32m   │
╰─────────────────────────────────────────────────────────────╯
```

---

## Các lệnh

### Xem quota

| Lệnh | Mô tả |
|------|-------|
| `anti-gravity-quota` | Xem quota account hiện tại |
| `anti-gravity-quota --all` | Xem quota của tất cả accounts |
| `anti-gravity-quota --json` | Output dạng JSON |
| `anti-gravity-quota -c` hoặc `--compact` | Output compact (1 dòng) |
| `anti-gravity-quota --watch` | Auto-refresh theo thời gian thực |

### Quản lý accounts

| Lệnh | Mô tả |
|------|-------|
| `anti-gravity-quota accounts list` | Liệt kê tất cả accounts |
| `anti-gravity-quota accounts switch <email>` | Chuyển sang account khác |
| `anti-gravity-quota accounts add <email>` | Thêm account mới |
| `anti-gravity-quota accounts remove <email>` | Xóa account |

### Cấu hình

| Lệnh | Mô tả |
|------|-------|
| `anti-gravity-quota config show` | Hiện cấu hình hiện tại |
| `anti-gravity-quota config edit` | Mở file config trong editor |
| `anti-gravity-quota config reset` | Reset về cấu hình mặc định |

---

## Chi tiết các chế độ hiển thị

### Table mode (mặc định)

```bash
anti-gravity-quota
```
Hiển thị bảng đầy đủ với progress bar.

### JSON mode

```bash
anti-gravity-quota --json
```

```json
{
  "account": "ngoanhtuan245@gmail.com",
  "timestamp": "2026-01-26T05:30:00Z",
  "models": [
    {
      "name": "Gemini 3 Pro",
      "used": 1500,
      "limit": 5000,
      "remaining_percent": 70,
      "reset_at": "2026-01-26T10:00:00Z"
    }
  ]
}
```

### Compact mode

```bash
anti-gravity-quota -c
```

```
Gemini:70% | Claude:80% | GPT:100% | Reset:4h32m
```

Phù hợp để embed vào status bar hoặc tmux.

### Watch mode

```bash
anti-gravity-quota --watch
# hoặc với interval tùy chỉnh
anti-gravity-quota --watch --interval 30
```

Auto-refresh mỗi 30 giây (mặc định 60s).

---

## Cấu hình

File cấu hình: `~/.config/anti-gravity-quota/config.yaml`

```yaml
# Danh sách accounts
accounts:
  - email: ngoanhtuan245@gmail.com
    active: true
  - email: work@company.com
    active: false

# Hiển thị
display:
  format: table    # table | json | compact
  refresh: 60      # giây (cho watch mode)
  
# Cảnh báo (tính năng tương lai)
alerts:
  low_quota_threshold: 20  # %
  notify: false
```

---

## Ví dụ sử dụng

### 1. Kiểm tra nhanh trước khi code

```bash
anti-gravity-quota -c
# Output: Gemini:70% | Claude:80% | GPT:100% | Reset:4h32m
```

### 2. Theo dõi trong lúc làm việc

```bash
# Mở watch mode trong một terminal tab
anti-gravity-quota --watch
```

### 3. Script tự động kiểm tra

```bash
#!/bin/bash
# check-quota.sh

QUOTA=$(anti-gravity-quota --json | jq '.models[0].remaining_percent')

if [ "$QUOTA" -lt 20 ]; then
    notify-send "AGQ Warning" "Gemini quota low: ${QUOTA}%"
fi
```

### 4. Xem quota của work account

```bash
anti-gravity-quota accounts switch work@company.com
anti-gravity-quota
```

---

## Troubleshooting

### "Antigravity not running"

**Nguyên nhân:** Không tìm thấy Antigravity Language Server.

**Giải pháp:** Đảm bảo Antigravity IDE đang mở.

### "Authentication failed"

**Nguyên nhân:** Không thể đọc OAuth credentials.

**Giải pháp:** 
1. Kiểm tra file `~/.gemini/oauth_creds.json` tồn tại
2. Thử đăng nhập lại trong Antigravity IDE

### "Port not found"

**Nguyên nhân:** Không xác định được port của Language Server.

**Giải pháp:** Restart Antigravity IDE.

---

## Tips & Tricks

### Tích hợp với tmux status bar

```bash
# Thêm vào ~/.tmux.conf
set -g status-right '#(anti-gravity-quota -c 2>/dev/null || echo "AG offline")'
```

### Alias tiện lợi

```bash
# Thêm vào ~/.bashrc hoặc ~/.zshrc
alias q='anti-gravity-quota -c'
alias qq='anti-gravity-quota'
alias qw='anti-gravity-quota --watch'
```

### Cron job kiểm tra quota

```bash
# Crontab: Kiểm tra mỗi giờ và log
0 * * * * anti-gravity-quota --json >> ~/.local/share/anti-gravity-quota/quota_history.jsonl
```

---

## FAQ

**Q: AGQ có gửi dữ liệu ra internet không?**
A: Không. Tool chỉ communicate với Language Server trên localhost.

**Q: Tôi có thể dùng khi Antigravity không chạy không?**
A: Không. Tool cần Language Server để lấy quota data.

**Q: Quota refresh mỗi bao lâu?**  
A: Tùy vào subscription tier:
- Free: Weekly
- Pro: Mỗi 5 giờ  
- Ultra: Mỗi 5 giờ (với limit cao hơn)
