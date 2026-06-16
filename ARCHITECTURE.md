# Architecture & Design Decisions

This document explains the architectural choices made in this Go backend project and the reasoning behind each decision.

## Table of Contents

1. [Layered Architecture](#layered-architecture)
2. [Why GoFiber](#why-gofiber)
3. [Why SQLC over an ORM](#why-sqlc-over-an-orm)
4. [Why PostgreSQL](#why-postgresql)
5. [Why go-playground/validator](#why-go-playgroundvalidator)
6. [Why Uber Zap](#why-uber-zap)
7. [Middleware Design](#middleware-design)
8. [Age Calculation Strategy](#age-calculation-strategy)
9. [Error Handling Strategy](#error-handling-strategy)
10. [Docker Strategy](#docker-strategy)

---

## Layered Architecture

```
┌──────────────────────────────────────────┐
│              HTTP Layer                  │
│         (GoFiber + Middleware)           │
├──────────────────────────────────────────┤
│              Handler Layer               │
│     (Parse request → call service)       │
├──────────────────────────────────────────┤
│              Service Layer               │
│  (Business logic, validation, age calc)  │
├──────────────────────────────────────────┤
│            Repository Layer              │
│      (Data access via SQLC queries)      │
├──────────────────────────────────────────┤
│              Database                    │
│           (PostgreSQL)                   │
└──────────────────────────────────────────┘
```

### Why This Pattern?

**Separation of Concerns** — Each layer has a single responsibility:

| Layer | Responsibility | Knows About |
|-------|---------------|-------------|
| **Handler** | HTTP concerns (parse params, set status codes, format JSON) | Service |
| **Service** | Business logic (validate inputs, compute age, orchestrate) | Repository, Models |
| **Repository** | Data access (run queries, map DB types) | SQLC, Database |
| **Models** | Data structures and validation rules | Nothing |

**Benefits:**
- **Testability**: Each layer can be tested independently with mock interfaces
- **Maintainability**: Changes to DB schema don't ripple into HTTP handlers
- **Readability**: New developers can understand the flow by reading top-down (handler → service → repository)

**Alternative considered:** *Flat structure* (handlers calling DB directly). Rejected because it mixes concerns, making testing hard and creating tight coupling between HTTP and database code.

---

## Why GoFiber

**Decision:** Use [GoFiber v3](https://gofiber.io/) as the HTTP framework.

**Reasoning:**

1. **Performance** — Built on `fasthttp` instead of Go's `net/http`. In benchmarks, GoFiber handles significantly more requests per second with lower memory allocation. For a CRUD API, this headroom matters as the user base scales.

2. **Developer Experience** — Express.js-like API design makes it intuitive:
   ```go
   app.Get("/users/:id", handler.GetUserByID)
   ```
   This reduces boilerplate compared to the standard library's `http.ServeMux`.

3. **Built-in Features** — Rate limiting, CORS, helmet, compression are available as middleware without third-party dependencies.

4. **Middleware Ecosystem** — Rich middleware support matches the requirements for request ID injection and request logging.

**Alternatives considered:**
- **`net/http` + gorilla/mux**: More idiomatic but more boilerplate. Gorilla is now archived.
- **Gin**: Similar DX to Fiber but uses `net/http` under the hood. Fiber's fasthttp foundation provides better raw throughput.
- **Echo**: Comparable to Fiber but smaller ecosystem.

---

## Why SQLC over an ORM

**Decision:** Use [SQLC](https://sqlc.dev/) for database access code generation.

**Reasoning:**

1. **Type Safety at Compile Time** — SQLC generates Go structs and methods from your SQL queries. If a query references a column that doesn't exist, you get a compilation error, not a runtime crash.

2. **No Runtime Reflection** — ORMs like GORM use reflection to map structs to queries at runtime. SQLC generates plain Go code with zero reflection, resulting in faster query execution and lower memory overhead.

3. **SQL Is the Source of Truth** — Developers write actual SQL queries. This means:
   - Full access to PostgreSQL features (CTEs, window functions, lateral joins)
   - No "ORM abstraction leak" where generated SQL is suboptimal
   - DBAs can review and optimize queries directly

4. **Generated Code Is Auditable** — The generated `users.sql.go` is committed to version control. You can review exactly what code runs against your database.

5. **Interface Generation** — SQLC generates a `Querier` interface, which enables mocking in tests without any additional tooling.

**Alternatives considered:**
- **GORM**: Feature-rich ORM but generates SQL that's hard to predict, uses heavy reflection, and makes complex queries awkward.
- **sqlx**: Good middle ground (raw SQL with struct scanning) but requires manual error-prone struct mapping. SQLC automates this.

**Trade-off:** SQLC requires a compilation step (`sqlc generate`) when queries change. This is minimal friction compared to the safety benefits.

---

## Why PostgreSQL

**Decision:** Use PostgreSQL 17 as the database.

**Reasoning:**

1. **Native DATE Type** — PostgreSQL has a first-class `DATE` type that stores only year-month-day without time zone complications. Perfect for `dob` storage.

2. **SERIAL Type** — Auto-incrementing integers for primary keys without needing extra sequences or setup.

3. **pgx Driver** — The `jackc/pgx` driver is the most performant PostgreSQL driver for Go, with native support for connection pooling (`pgxpool`), prepared statements, and PostgreSQL-specific types.

4. **Production Proven** — PostgreSQL is the industry standard for relational data with excellent tooling, monitoring, and community support.

---

## Why go-playground/validator

**Decision:** Use [go-playground/validator v10](https://github.com/go-playground/validator) for input validation.

**Reasoning:**

1. **Struct Tag Validation** — Define validation rules directly on struct fields:
   ```go
   Name string `validate:"required,min=1,max=255"`
   ```
   This co-locates the rule with the data definition, improving readability.

2. **Custom Validators** — The `dateformat` tag was registered to validate `YYYY-MM-DD` format and reject future dates. This is clean and reusable.

3. **Detailed Error Reporting** — `ValidationErrors` provides per-field error information, enabling clear API error responses:
   ```json
   {"error": "validation failed", "details": {"name": "name is required"}}
   ```

4. **Performance** — Validator v10 uses reflection only once during registration, then validates via compiled function chains. It's one of the fastest validators in the Go ecosystem.

**Alternative considered:** Manual validation with `if/else` chains. This works for small projects but becomes unwieldy quickly and doesn't standardize error formats.

---

## Why Uber Zap

**Decision:** Use [Uber Zap](https://github.com/uber-go/zap) for structured logging.

**Reasoning:**

1. **Zero-Allocation Logging** — Zap avoids `fmt.Sprintf` and reflection. Log calls that are below the configured level (e.g., `Debug` in production) allocate zero memory.

2. **Structured Fields** — Logs are key-value pairs, not formatted strings:
   ```go
   log.Info("user created", zap.Int32("id", 1), zap.String("name", "Alice"))
   ```
   This makes logs machine-parseable for tools like ELK, Datadog, or Grafana Loki.

3. **Environment-Aware Configuration** —
   - **Development**: Colorized, human-readable console output
   - **Production**: JSON format for log aggregation pipelines

4. **Caller Information** — Automatically includes file and line number in log entries for debugging.

**Alternative considered:** `log/slog` (Go 1.21+ standard library). While improving, it lacks Zap's performance characteristics and the ecosystem of encoders/sinks.

---

## Middleware Design

### Request ID Middleware

**What it does:** Generates a UUID v4 for every request and attaches it as the `X-Request-ID` response header.

**Why:**
- **Traceability** — When debugging a production issue, the request ID links a user's report to the exact log entries.
- **Idempotent** — If the client sends an `X-Request-ID` header, the middleware reuses it instead of generating a new one. This supports distributed tracing.

### Request Logger Middleware

**What it does:** Logs method, path, status code, duration, and request ID for every request.

**Why:**
- **Observability** — Every request is logged with its duration, enabling latency monitoring without external APM tools.
- **Audit Trail** — Provides a chronological record of all API interactions.
- **Zero Config** — No need for external logging middleware libraries.

**Design Choice:** Both middleware run before route handlers, ensuring all requests (including 404s) are logged and get a request ID.

---

## Age Calculation Strategy

**Decision:** Calculate age dynamically in Go code using `time.Now()`, not in SQL.

**Reasoning:**

1. **No Stored Derived Data** — Age changes daily. Storing it in the database would require a daily cron job to update all rows, which is fragile and wasteful.

2. **Correctness** — The Go `CalculateAge` function uses direct month/day comparison to handle leap year edge cases (e.g., someone born on Feb 29).

3. **Testability** — A pure Go function is trivially unit-testable with table-driven tests. SQL-based age calculation would require a running database to test.

4. **Performance** — Computing age is a subtraction and two comparisons — O(1) with negligible CPU cost. There's no benefit to pushing this to the database.

**Implementation detail:** Age is only included in GET responses (`GetUserByID`, `ListUsers`), not in create/update responses, matching the API specification.

---

## Error Handling Strategy

The API uses a consistent error response format:

```json
{
  "error": "human-readable message",
  "details": { "field": "field-specific message" }
}
```

**Error mapping in the handler layer:**

| Error Type | HTTP Status | Example |
|---|---|---|
| Validation error | `400 Bad Request` | Invalid date format, missing name |
| Parse error | `400 Bad Request` | Non-integer ID parameter |
| Not found | `404 Not Found` | User ID doesn't exist |
| Internal error | `500 Internal Server Error` | Database connection failure |

**Design choice:** Sentinel errors (like `ErrUserNotFound`) in the repository layer bubble up through the service layer and are caught by `handleError()` in the handler layer. This keeps error translation centralized.

---

## Docker Strategy

**Multi-stage build:**

```
Stage 1: golang:1.25-alpine (builder)
  → Downloads modules, compiles binary

Stage 2: alpine:3.21 (runtime)
  → Copies only the binary + migrations
  → Final image is ~15-20 MB
```

**Why multi-stage:**
- The Go toolchain and source code are discarded after compilation
- The runtime image has a minimal attack surface
- Image pull times are fast for CI/CD pipelines

**Docker Compose orchestration:**
- PostgreSQL has a health check (`pg_isready`)
- The Go backend uses `depends_on` with `condition: service_healthy`
- This ensures the database is accepting connections before the app starts

---

## Summary of Trade-offs

| Decision | Benefit | Trade-off |
|---|---|---|
| Layered architecture | Testable, maintainable | More files and boilerplate |
| SQLC | Type-safe, fast, auditable | Requires codegen step |
| GoFiber | High performance, great DX | Not `net/http` compatible |
| Zap | Zero-alloc structured logging | Heavier API than `log/slog` |
| Dynamic age calc | No stale data, testable | Computed on every request |
| Multi-stage Docker | Tiny image, secure | Slightly longer build time |

Each decision prioritizes **correctness, performance, and maintainability** appropriate for a production-grade REST API.
