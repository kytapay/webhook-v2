# KytaPay Webhook v2

Webhook service untuk menangani callback dari LinkQu dan PakaiLink.

## 📖 Dokumentasi

- **[HOW_TO_RUN.md](./HOW_TO_RUN.md)** - Link ke deployment guide (lihat `../api-v2/HOW_TO_RUN.md`)
- **[DATABASE_SETUP.md](./DATABASE_SETUP.md)** - Setup database untuk cPanel shared hosting
- **[deploy/](./deploy/)** - File konfigurasi deployment (systemd, nginx)

## 🚀 Quick Start (Development)

```bash
# Copy environment file
cp env.example .env

# Edit .env dengan konfigurasi yang sesuai
nano .env

# Install dependencies
go mod download

# Run service
go run main.go
```

## 🌐 Endpoints

### Payment Webhooks
- `POST /payments/linkqu/qris` - LinkQu QRIS webhook
- `POST /payments/linkqu/ewallet` - LinkQu E-Wallet webhook
- `POST /payments/pakailink/va` - PakaiLink VA webhook

### Payout Webhooks
- `POST /payouts/linkqu/bank` - LinkQu Bank payout webhook
- `POST /payouts/linkqu/ewallet` - LinkQu E-Wallet payout webhook
- `POST /payouts/pakailink/bank` - PakaiLink Bank payout webhook
- `POST /payouts/pakailink/ewallet` - PakaiLink E-Wallet payout webhook

### Health
- `GET /health` - Health check endpoint

## 🔐 Validasi

- **LinkQu**: Validasi menggunakan `client-id` dan `client-secret` dari header
- **PakaiLink**: Validasi menggunakan `X-SIGNATURE` dengan symmetric signature (HMAC SHA-512) atau asymmetric signature (RSA SHA-256)

## 📤 Response

Semua webhook endpoint selalu mengembalikan HTTP 200 OK dengan response:
```json
{
  "responseCode": "2002800",
  "responseMessage": "Successful"
}
```

## 📝 Port

Default port: **8081** (dapat diubah via `WEBHOOK_PORT` environment variable)

## 🔐 Security

**⚠️ PENTING: File-file berikut TIDAK BOLEH di-push ke Git:**
- `.env` files
- `*.pem` files (RSA keys)
- Binary files
- Log files

Pastikan semua file sensitif sudah ada di `.gitignore`!

