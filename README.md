# Ainyx Solution — Go Backend API

A RESTful API built with **Go** to manage users with their `name` and `dob` (date of birth). The API dynamically calculates and returns a user's **age** when fetching user details.

## Tech Stack

| Technology | Purpose |
|---|---|
| [GoFiber v3](https://gofiber.io/) | HTTP framework (Express-like, built on fasthttp) |
| [PostgreSQL 17](https://www.postgresql.org/) | Relational database |
| [SQLC](https://sqlc.dev/) | Type-safe SQL → Go code generation |
| [pgx v5](https://github.com/jackc/pgx) | PostgreSQL driver with connection pooling |
| [go-playground/validator v10](https://github.com/go-playground/validator) | Struct-tag input validation |
| [Uber Zap](https://github.com/uber-go/zap) | Structured, zero-allocation logging |
| [Docker](https://www.docker.com/) | Containerization |

## Project Structure

```
go-backend/
├── cmd/server/main.go           # Application entry point
├── config/config.go             # Environment-based configuration
├── db/
│   ├── migrations/              # SQL migration files
│   └── sqlc/
│       ├── queries/users.sql    # SQLC query definitions
│       └── generated/           # SQLC-generated Go code (DO NOT EDIT)
├── internal/
│   ├── handler/                 # HTTP handlers (GoFiber)
│   ├── service/                 # Business logic layer
│   ├── repository/              # Data access layer
│   ├── models/                  # Request/response types, validation, age calc
│   ├── middleware/               # Request ID & request logger middleware
│   ├── routes/                  # Route registration
│   └── logger/                  # Uber Zap logger setup
├── sqlc.yaml                    # SQLC configuration
├── Dockerfile                   # Multi-stage Docker build
├── .env                         # Local environment variables
└── go.mod / go.sum              # Go module files
```

## Prerequisites

- **Go** 1.25+
- **PostgreSQL** 17+ (or Docker)
- **Docker & Docker Compose** (optional, for containerized setup)

## Getting Started

### Option 1: Run with Docker Compose (Recommended)

```bash
# From the project root (Ainyx_Solution_Assignment/)
docker compose up --build -d

# The API will be available at http://localhost:9090
# PostgreSQL runs on port 5432
```

### Option 2: Run Locally

1. **Start PostgreSQL** (via Docker or local install):
   ```bash
   docker compose up postgresql -d
   ```

2. **Configure environment** — edit `go-backend/.env`:
   ```env
   SERVER_PORT=9090
   DATABASE_URL=postgres://admin:secret_password@localhost:5432/ainyx?sslmode=disable
   APP_ENV=development
   ```

3. **Run the server**:
   ```bash
   cd go-backend
   go run ./cmd/server/
   ```

4. The API will be available at `http://localhost:9090`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/users` | Create a new user |
| `GET` | `/users` | List all users (paginated) |
| `GET` | `/users/:id` | Get a user by ID (includes calculated age) |
| `PUT` | `/users/:id` | Update a user |
| `DELETE` | `/users/:id` | Delete a user |
| `GET` | `/health` | Health check endpoint |

### Create User

```bash
curl -X POST http://localhost:9090/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "dob": "1990-05-10"}'
```

**Response** (201 Created):
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### Get User by ID

```bash
curl http://localhost:9090/users/1
```

**Response** (200 OK):
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}
```

### List All Users (Paginated)

```bash
curl "http://localhost:9090/users?page=1&limit=10"
```

**Response** (200 OK):
```json
{
  "data": [
    {
      "id": 1,
      "name": "Alice",
      "dob": "1990-05-10",
      "age": 35
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 1,
  "total_pages": 1
}
```

### Update User

```bash
curl -X PUT http://localhost:9090/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Updated", "dob": "1991-03-15"}'
```

**Response** (200 OK):
```json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

### Delete User

```bash
curl -X DELETE http://localhost:9090/users/1
```

**Response**: `204 No Content`

## Input Validation

- `name`: Required, 1–255 characters
- `dob`: Required, must be in `YYYY-MM-DD` format, must not be in the future

Invalid input returns `400 Bad Request` with field-level error messages:
```json
{
  "error": "validation failed",
  "details": {
    "name": "name is required",
    "dob": "dob must be a valid date in YYYY-MM-DD format and not in the future"
  }
}
```

## Running Tests

```bash
cd go-backend
go test ./internal/models/... -v
```

## Regenerating SQLC Code

If you modify SQL queries or migrations:

```bash
cd go-backend
.\sqlc_bin\sqlc.exe generate
```

## Middleware

- **X-Request-ID**: Every response includes a unique `X-Request-ID` header for traceability
- **Request Logger**: Logs method, path, status code, duration, and request ID for every request

## License

This project is for the Ainyx Solutions assignment evaluation.
