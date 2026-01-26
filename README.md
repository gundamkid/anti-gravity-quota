# Anti-Gravity Quota - CLI Tool

CLI tool nhẹ, nhanh để quản lý và hiển thị quota của các AI models trong Antigravity, hỗ trợ nhiều Google accounts.

## Tính năng chính

- 📊 **Theo dõi quota real-time**: Xem quota đã sử dụng, còn lại, và thời gian refresh
- 👥 **Multi-account**: Quản lý nhiều Google accounts
- 🚀 **Nhanh & Nhẹ**: Viết bằng Go, binary ~5MB, khởi động tức thì
- 🖥️ **Cross-platform**: Hỗ trợ Linux, macOS, Windows

## Cài đặt

```bash
# Từ source
git clone https://github.com/your-username/anti-gravity-quota.git
cd anti-gravity-quota
go build -o anti-gravity-quota ./cmd/anti-gravity-quota

# Di chuyển vào PATH
sudo mv anti-gravity-quota /usr/local/bin/
```

## Sử dụng nhanh

```bash
# Xem quota của account hiện tại
anti-gravity-quota

# Xem quota của tất cả accounts
anti-gravity-quota --all

# Xem dạng JSON
anti-gravity-quota --json

# Xem compact (1 dòng)
anti-gravity-quota -c

# Watch mode (auto refresh)
anti-gravity-quota --watch
```

## Yêu cầu

- Antigravity IDE đang chạy (để detect Language Server)
- Go 1.21+ (nếu build từ source)

## Tài liệu

- [Hướng dẫn sử dụng](docs/user_guide.md)
- [Tài liệu kỹ thuật](docs/technical.md)
- [Implementation Plan](docs/implementation/implementation_plan.md)

## License

MIT License
