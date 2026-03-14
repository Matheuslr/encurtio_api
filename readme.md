# 🔗 Encurtio

A fast, lightweight URL shortener API built with **Go**, **Gin** and **Apache Cassandra**.

URLs are shortened using a **SHA-256 + Base62** encoding strategy that produces compact ~8-character codes from any input URL.

---

## 📋 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Architecture](#-architecture)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [Configuration](#-configuration)
- [API Reference](#-api-reference)
- [Testing](#-testing)
- [Makefile Commands](#-makefile-commands)
- [Docker](#-docker)
- [How the Shortener Works](#-how-the-shortener-works)

---

## ✨ Features

- **URL Shortening** — `POST` a long URL, get back a short code
- **URL Redirect** — `GET /:code` redirects (302) to the original URL
- **Health Check** — `GET /api/v1/health` for monitoring
- **Base62 Short Codes** — deterministic, ~8-character codes via SHA-256 + Base62
- **Cassandra Storage** — highly available, write-optimized persistence
- **Panic Recovery Middleware** — graceful error handling on unexpected failures
- **Unit Tests** — 22 tests covering all layers (handler → service → shortener → config → middleware)
- **Docker Compose** — one-command setup for the entire stack
- **Hot Reload** — development with [Air](https://github.com/air-verse/air)

---

## 🛠 Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.25+ |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| Database | Apache Cassandra 5 |
| Encoding | SHA-256 → Base62 ([jxskiss/base62](https://github.com/jxskiss/base62)) |
| Config | Environment variables via [godotenv](https://github.com/joho/godotenv) |
| Testing | [testify](https://github.com/stretchr/testify) (assert + mock) |
| Hot Reload | [Air](https://github.com/air-verse/air) |
| Containerization | Docker + Docker Compose |

---

## 🏗 Architecture

```
Client
  │
  ▼
┌──────────────────────┐
│   Gin HTTP Router    │  ← middleware (logger, panic recovery)
├──────────────────────┤
│      Handlers        │  ← URL handler, Health handler
├──────────────────────┤
│      Services        │  ← Business logic (shorten, resolve, build URL)
├──────────────────────┤
│    Repositories      │  ← Data access layer (interface + Cassandra impl)
├──────────────────────┤
│  Apache Cassandra    │  ← Persistence
└──────────────────────┘
```

### Shortening Flow

```
POST /api/v1/url/shorten  { "url": "https://example.com/..." }
  → Handler parses & validates request body
    → Service calls shortener.Encode(url)
      → SHA-256(url) → take first 6 bytes → Base62 encode → "sRVrAaAE"
    → Service saves { short_code, original_url, created_at } to Cassandra
  → Returns 201 { "short_url": "http://localhost:8080/sRVrAaAE" }
```

### Redirect Flow

```
GET /sRVrAaAE
  → Handler extracts :code from path
    → Service queries Cassandra by short_code
  → Returns 302 redirect → Location: https://example.com/...
```

---

## 📁 Project Structure

```
encurtio/
├── cmd/
│   └── api/
│       └── main.go              # Application entrypoint
├── configs/
│   └── config.go                # Environment-based configuration
├── internal/
│   ├── database/
│   │   └── cassandra.go         # Cassandra session factory
│   ├── handler/
│   │   ├── healthcheck.go       # GET /api/v1/health
│   │   └── url_handler.go       # POST /shorten + GET /:code
│   ├── middleware/
│   │   └── error_middleware.go   # Panic recovery + error logging
│   ├── model/
│   │   └── url.go               # URL domain model
│   ├── repository/
│   │   ├── url_repository.go    # Repository interface
│   │   └── cassandra_url_repository.go  # Cassandra implementation
│   ├── service/
│   │   └── url_service.go       # Business logic
│   └── shortener/
│       └── base62.go            # SHA-256 + Base62 encoder
├── test/                        # All unit tests (mirrors internal/ tree)
│   ├── mocks/
│   │   └── url_repository_mock.go
│   ├── configs/
│   │   └── config_test.go
│   ├── handler/
│   │   ├── url_handler_test.go
│   │   └── healthcheck_test.go
│   ├── middleware/
│   │   └── error_middleware_test.go
│   ├── service/
│   │   └── url_service_test.go
│   └── shortener/
│       └── base62_test.go
├── cassandra/
│   └── init.cql                 # Keyspace + table bootstrap
├── migrations/
│   ├── 000001_create_keyspace.up.cql
│   ├── 000001_create_keyspace.down.cql
│   ├── 000002_create_urls_table.up.cql
│   └── 000002_create_urls_table.down.cql
├── docker-compose.yaml
├── dockerfile
├── Makefile
├── go.mod
├── go.sum
├── .env                         # Local env vars (not committed)
└── .air.toml                    # Hot-reload config
```

---

## 🚀 Getting Started

### Prerequisites

- **Go** 1.25+
- **Docker** & **Docker Compose** (for Cassandra)
- **Make** (optional, for convenience)

### 1. Clone the repository

```bash
git clone https://github.com/Matheuslr/encurtio_api.git
cd encurtio_api
```

### 2. Start Cassandra

```bash
make cassandra-up
# or
docker compose up -d cassandra
```

> Wait ~30 seconds for Cassandra to fully initialize. The `init.cql` file will automatically create the keyspace and table.

### 3. Configure environment

Create a `.env` file in the project root (or export the variables):

```env
APP_PORT=8080
APP_URL=http://localhost:8080
CASSANDRA_HOST=127.0.0.1
CASSANDRA_KEYSPACE=encurtio
```

### 4. Run the API

```bash
# directly
make run

# or with hot-reload
make air

# or with Docker (API + Cassandra)
make docker-up
```

---

## ⚙️ Configuration

All configuration is done via environment variables (loaded from `.env` if present):

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_PORT` | `8080` | HTTP server port |
| `APP_URL` | `http://localhost:8080` | Base URL used to build short URLs |
| `CASSANDRA_HOST` | `127.0.0.1` | Cassandra node address |
| `CASSANDRA_KEYSPACE` | `encurtio` | Cassandra keyspace name |

---

## 📡 API Reference

### Shorten URL

```http
POST /api/v1/url/shorten
Content-Type: application/json

{
  "url": "https://example.com/some/very/long/path?with=params"
}
```

**Response** — `201 Created`

```json
{
  "short_url": "http://localhost:8080/sRVrAaAE"
}
```

### Redirect

```http
GET /:code
```

**Response** — `302 Found`

```
Location: https://example.com/some/very/long/path?with=params
```

**Example with curl:**

```bash
# Shorten
curl -s -X POST http://localhost:8080/api/v1/url/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'

# Redirect (follow)
curl -L http://localhost:8080/sRVrAaAE

# Redirect (inspect headers)
curl -I http://localhost:8080/sRVrAaAE
```

### Health Check

```http
GET /api/v1/health
```

**Response** — `200 OK`

```json
{
  "message": "Health check at 2026-03-14T12:00:00Z"
}
```

---

## 🧪 Testing

Tests live in the `test/` directory, mirroring the `internal/` structure. They use **black-box testing** (`_test` package suffix) with **testify** for assertions and mocks.

| Test Suite | Tests | Covers |
|------------|-------|--------|
| `test/shortener/` | 5 | Base62 encoding (short, deterministic, unique, empty input) |
| `test/service/` | 6 | Shorten, GetOriginalURL, BuildShortURL (success + error) |
| `test/handler/` | 8 | HTTP handlers (Shorten, Redirect, Health, error cases) |
| `test/middleware/` | 3 | Panic recovery, normal flow, handler errors |
| `test/configs/` | 2 | Default config, env var overrides |

```bash
# Run all tests
make test

# Run with coverage report
make test-cover

# Open HTML coverage in browser
make test-cover-html
```

---

## 🔧 Makefile Commands

```bash
make help
```

| Command | Description |
|---------|-------------|
| `make build` | Build the binary to `./bin/encurtio` |
| `make run` | Run locally with `go run` |
| `make air` | Run with hot-reload (Air) |
| `make test` | Run all unit tests |
| `make test-cover` | Tests + coverage report |
| `make test-cover-html` | Open coverage in browser |
| `make lint` | Run `go vet` + `gofmt` |
| `make tidy` | `go mod tidy` |
| `make deps` | Download dependencies |
| `make docker-up` | Start all containers (Cassandra + API) |
| `make docker-down` | Stop all containers |
| `make docker-logs` | Tail container logs |
| `make docker-restart` | Restart all containers |
| `make cassandra-up` | Start only Cassandra |
| `make cassandra-cql` | Open `cqlsh` shell |
| `make clean` | Remove build artefacts |

---

## 🐳 Docker

### Full stack (API + Cassandra)

```bash
docker compose up -d --build
```

### Cassandra only (for local development)

```bash
docker compose up -d cassandra
```

### Database schema

The schema is automatically created via `cassandra/init.cql`:

```sql
CREATE KEYSPACE IF NOT EXISTS encurtio
WITH replication = { 'class': 'SimpleStrategy', 'replication_factor': 1 };

CREATE TABLE IF NOT EXISTS urls (
    short_code text PRIMARY KEY,
    original_url text,
    created_at timestamp
);
```

---

## 🧠 How the Shortener Works

The encoding pipeline in `internal/shortener/base62.go`:

```
Input URL
    │
    ▼
SHA-256 hash (32 bytes)
    │
    ▼
Truncate to first 6 bytes (48 bits ≈ 281 trillion unique codes)
    │
    ▼
Base62 encode → ~8-character alphanumeric string
```

- **Deterministic** — same URL always produces the same code
- **Fixed-length** — output is always ~8 characters regardless of input size
- **Collision-resistant** — 48 bits gives ~281 trillion unique codes

Example outputs:

| Input | Code | Length |
|-------|------|--------|
| `https://www.google.com/search?q=golang` | `zmLlHqRj` | 8 |
| `https://github.com/Matheuslr/encurtio_api` | `rkVPGgTO` | 8 |
| `https://example.com` | `sRVrAaAE` | 8 |

---

## 📝 License

This project is open source and available under the [MIT License](LICENSE).