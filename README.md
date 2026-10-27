# URL Shortener

A lightweight URL shortening service written in Go.

The service generates short aliases for long URLs, stores them in SQLite, and redirects users from the short URL to the original address.

## Features

- Create short URLs
- Generate random aliases
- Set custom aliases
- Redirect short URLs to original addresses
- SQLite database storage
- REST API
- Basic authentication for URL creation
- Request validation
- Structured logging
- Request ID middleware
- Database initialization on startup

## Tech Stack

- Go
- Chi Router
- SQLite
- Docker
- pgx-compatible project structure
- cleanenv
- validator
- slog

## Requirements

- Go 1.25+
- Docker Desktop
- SQLite
- Make is optional

## Project Structure

```text
my-url-shortener/
├── cmd/
│   └── url-shortener/
│       └── main.go
├── config/
│   └── config.yaml
├── internal/
│   ├── config/
│   ├── http-server/
│   ├── lib/
│   └── storage/
├── storage/
│   └── storage.db
├── main.go
├── go.mod
└── go.sum
```

## Configuration

The application requires the `CONFIG_PATH` environment variable.

Before running the application, update `config/config.yaml`:

```yaml
env: "local"

storage_path: "./storage/storage.db"

http_server:
  address: "localhost:8082"
  timeout: 4s
  idle_timeout: 30s
  user: "myuser"
  password: "mypass"
```

The `storage_path` value should point to a local SQLite database file.

## Run the Application

Run the application from the project root.

### Git Bash

```bash
CONFIG_PATH=config/config.yaml go run ./cmd/url-shortener
```

### PowerShell

```powershell
$env:CONFIG_PATH="config/config.yaml"
go run ./cmd/url-shortener
```

The server will be available at:

```text
http://localhost:8082
```

## API Endpoints

### Create a Short URL

```http
POST /url/
```

This endpoint requires Basic Authentication.

Request:

```json
{
  "url": "https://github.com/stfprod17",
  "alias": "stefan"
}
```

The `alias` field is optional. If it is not provided, the service generates a random six-character alias.

Example using cURL:

```bash
curl -X POST http://localhost:8082/url/ \
  -u myuser:mypass \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com/stfprod17","alias":"stefan"}'
```

Example response:

```json
{
  "status": "OK",
  "alias": "stefan"
}
```

### Redirect to the Original URL

```http
GET /{alias}
```

Example:

```text
http://localhost:8082/stefan
```

The service redirects the request to the original URL.

## Authentication

URL creation is protected with Basic Authentication.

Credentials are configured in:

```yaml
http_server:
  user: "myuser"
  password: "mypass"
```

The redirect endpoint does not require authentication.

## Error Handling

The service returns errors for:

- Invalid request bodies
- Invalid URLs
- Existing aliases
- Missing aliases
- Internal storage errors

## Storage

The application uses SQLite.

The database file is created at the path specified by `storage_path` in `config/config.yaml`.

The `url` table contains:

- Original URL
- Short alias
- Record ID

## Development

Download dependencies:

```bash
go mod download
```

Format the code:

```bash
gofmt -w .
```

Run the application:

```bash
CONFIG_PATH=config/config.yaml go run ./cmd/url-shortener
```

## Future Improvements

- Add automated tests
- Add click statistics
- Add expiration dates for links
- Add link deletion
- Add user accounts
- Add rate limiting
- Add Dockerfile
- Add production configuration
- Add API documentation with Swagger
- Add custom error response models

## License

This project is intended for educational and personal use.
