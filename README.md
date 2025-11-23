# KytaPay Webhook v2

Webhook service untuk menangani callback dari LinkQu dan PakaiLink.

## Setup

1. Copy `env.example` ke `.env` dan isi dengan konfigurasi yang sesuai
2. Install dependencies: `go mod download`
3. Run service: `go run main.go`

## Port

Default port: **8081** (dapat diubah via `WEBHOOK_PORT` environment variable)

## Endpoints

- `POST /payments/linkqu/qris` - Webhook callback untuk QRIS dari LinkQu
- `POST /payments/linkqu/ewallet` - Webhook callback untuk E-Wallet dari LinkQu
- `POST /payments/pakailink/va` - Webhook callback untuk Virtual Account dari PakaiLink
- `GET /health` - Health check endpoint

## Validasi

- **LinkQu**: Validasi menggunakan `client-id` dan `client-secret` dari header
- **PakaiLink**: Validasi menggunakan `X-SIGNATURE` dengan symmetric signature

## Response

Semua webhook endpoint selalu mengembalikan HTTP 200 OK dengan response:
```json
{
  "responseCode": "2002800",
  "responseMessage": "Successful"
}
```

