# FocusList

A Go/Gin REST API learning project for organizing to-do items into buckets. PostgreSQL stores users, bucket membership, and items; Redis stores temporary registration and password-reset data. The source includes JWT middleware and bucket-sharing routes.

## Run locally

Prerequisites: the Go toolchain specified in [go.mod](go.mod) (Go 1.23.0 with toolchain 1.23.11), PostgreSQL, and Redis. The included [Docker Compose file](docker-compose.yml) starts Redis only; provision PostgreSQL separately and create a local database named `focuslist`.

```bash
git clone https://github.com/paranietharan/FocusList.git
cd FocusList
go mod download
docker compose up -d redis
```

Configure `.env` in the repository root with your own local database credentials:

```dotenv
DB_USERNAME=your_local_user
DB_PASSWORD=your_local_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=focuslist
JWT_SECRET_KEY=replace-with-a-local-random-secret
REDIS_ADDR=localhost
REDIS_PORT=6379
REDIS_DB=0
REDIS_USERNAME=
REDIS_PASSWORD=
```

The repository already tracks a `.env` file. Replace its values locally for your setup and do not commit credentials. `REDIS_ADDR` is the hostname only; the client appends `REDIS_PORT`.

For a fresh database, run the migration and then start the API:

```bash
go run ./migrations
go run ./cmd/server
```

The migration uses plain `CREATE TABLE` statements, so it is not an idempotent migration runner. It also creates a demonstration super-admin account; inspect [migrations/main.go](migrations/main.go) before running it against anything other than a local test database. The API listens on port 8080.

## Route groups

The complete route definitions are in [cmd/server/main.go](cmd/server/main.go).

| Area | Routes |
| --- | --- |
| Accounts | `POST /register`, `/verify-email`, `/login` |
| Password reset | `POST /reset-password`, `/confirm-reset-password` |
| Buckets | `POST /buckets`, `GET /buckets`, and GET/PUT/DELETE `/buckets/:bucketID` |
| Bucket membership | PUT/GET `/buckets/:bucketID/users`, DELETE `/buckets/:bucketID/users/:userEmail` |
| Items | POST/GET `/buckets/:bucketID/items`, and GET/PUT/DELETE `/buckets/:bucketID/items/:itemID` |

Use the DTOs and handlers for exact request bodies. Protected routes use the authentication middleware. Verification and reset codes are currently printed to the server console rather than sent through an email provider.

## Source guide

- [internal/handler](internal/handler): Request/response handling.
- [internal/service](internal/service): Account, bucket, and item logic.
- [internal/repository](internal/repository): PostgreSQL queries.
- [internal/cache](internal/cache): Redis access.
- [internal/config/config.go](internal/config/config.go): Environment configuration.
- [migrations/001_create_table.sql](migrations/001_create_table.sql): Database schema.

## Development status

This is a learning implementation, not a verified production service. Authentication and bucket ownership rules need dedicated tests. Some diagnostic paths log configuration and verification data, so use dummy credentials/data locally.

```bash
go build ./...
go test ./...
```

No Go test files are currently included. `internal/handler/test_handler.go` is an HTTP handler, not an automated test suite.
