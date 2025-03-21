# Aynshteyn Backend

Secure email collection backend for aynshteyn.dev

## Features

- REST API for email collection with browser information
- SQLite database for simple, reliable storage
- Security features:
  - Rate limiting
  - CORS protection
  - HMAC verification
  - Security headers
  - Input validation
- Admin endpoint for viewing collected data

## Setup

1. Install Go (1.18+ recommended)
2. Install dependencies:
```bash
go mod tidy
```

## Configuration

The server can be configured using command-line flags:

- `-port`: API server port (default: 4000)
- `-env`: Environment [development|staging|production] (default: development)
- `-db-dsn`: SQLite data source name (default: file:subscribers.db)
- `-cors-trusted-origins`: Trusted CORS origins (space separated, default: https://aynshteyn.dev)

## Run in Development

```bash
go run ./cmd/api
```

## Run in Production

```bash
# Build the binary
go build -o aynshteyn-backend ./cmd/api

# Run the server (with proper flags)
./aynshteyn-backend -env=production -port=4000 -cors-trusted-origins="https://aynshteyn.dev"
```

## API Endpoints

### POST /api/v1/subscribe

Collects email and browser information.

Request body:
```json
{
  "email": "user@example.com",
  "clientToken": "optional-client-token",
  "clientNonce": "random-string",
  "browserData": "browser information string",
  "timestamp": 1679012345,
  "clientHMAC": "optional-hmac-signature"
}
```

Response:
```json
{
  "success": true,
  "message": "Thank you for subscribing",
  "data": {
    "email_hash": "hashed-email",
    "created_at": "2023-03-17T12:00:00Z",
    "is_verified": false,
    "challenge_id": "abc123"
  }
}
```

### GET /api/v1/admin/subscribers

Admin endpoint to view all subscribers.

**Requires Basic Authentication**

Response:
```json
{
  "subscribers": [
    {
      "email": "user@example.com",
      "email_hash": "hashed-email",
      "user_agent": "Mozilla/5.0...",
      "ip_address": "192.168.1.1",
      "created_at": "2023-03-17T12:00:00Z",
      "is_verified": false,
      "challenge_id": "abc123"
    }
  ],
  "count": 1
}
```

## Security Notes

For production use:
1. Change the secret key in `handler/subscribe.go`
2. Use environment variables for secrets instead of hardcoding
3. Update admin credentials
4. Consider implementing rate limiting with a distributed cache 