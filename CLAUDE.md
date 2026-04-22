# TripSearch — CLAUDE.md

TripSearch is a travel listing platform for hotels and flights. Go 1.22+ backend (Gin), React 18 + TypeScript frontend (Vite), PostgreSQL primary database, Redis for caching/sessions/rate limiting, RabbitMQ for event-driven messaging, deployed on AWS (ECS Fargate, RDS, ElastiCache, S3, SES). This is a portfolio project built to demonstrate production-grade engineering depth.

---

## Commands

### Backend
```bash
cd backend
go run ./cmd/api/main.go          # start dev server
go test ./...                     # run all tests
go test ./internal/hotel/...      # run tests for a domain
golangci-lint run                 # lint
migrate -path db/migrations -database $DATABASE_URL up   # run migrations
sqlc generate                     # regenerate sqlc types after editing queries
```

### Frontend
```bash
cd frontend
npm run dev        # start dev server (Vite)
npm run build      # production build
npm run test       # Vitest
npm run lint       # ESLint
```

### Docker (local dev)
```bash
docker compose up -d              # start postgres, redis, rabbitmq
docker compose down               # stop all
```

---

## Architecture

### Layered Architecture — Handler → Service → Repository

Every domain in `backend/internal/` follows this exact four-file pattern:

| File | Responsibility |
|------|---------------|
| `handler.go` | Parse HTTP request, call service, write response. Never touches DB. |
| `service.go` | Business logic, validation, orchestration. Never knows about HTTP. |
| `repository.go` | Database queries only via sqlc. Returns domain types. Never contains business logic. |
| `types.go` | Structs, interfaces, sentinel errors for this domain. |

**The rule is absolute:** handlers never touch the database. Repositories never contain business logic. Services never import `net/http`.

### Request Lifecycle
```
HTTP Request
  → Gin engine
  → Global middleware: RequestID → Logger → Recovery → RateLimit
  → Route group middleware: JWT → RequireRole
  → handler.go     — parse + validate
  → service.go     — business logic
  → repository.go  — sqlc query
  → service.go     — transform result
  → handler.go     — c.JSON(200, response)
  → Logger middleware — write duration, status, request_id
  → HTTP Response
```

### Dependency Injection — Manual wiring in main.go only
All dependencies are wired manually in `backend/cmd/api/main.go`. No DI framework. No `init()` magic.

```go
// Example pattern — always follow this structure
hotelRepo    := hotel.NewRepository(db)
hotelService := hotel.NewService(hotelRepo, s3, broker)
hotelHandler := hotel.NewHandler(hotelService)
```

---

## Repo Structure

```
tripsearch/
├── frontend/
│   └── src/
│       ├── components/     # Shared UI components (presentational)
│       ├── pages/          # Route-level page components (containers)
│       ├── hooks/          # useListings, useAuth, useFavorites (TanStack Query)
│       ├── store/          # Zustand — auth state, filter state
│       └── lib/            # API client, utils
│
├── backend/
│   ├── cmd/api/main.go     # Entry point — bootstrap and wire all dependencies
│   ├── internal/
│   │   ├── hotel/          # handler.go, service.go, repository.go, types.go
│   │   ├── flight/         # handler.go, service.go, repository.go, types.go
│   │   ├── user/           # handler.go, service.go, repository.go, types.go
│   │   ├── search/         # handler.go, service.go, query_builder.go, sort_strategy.go
│   │   ├── auth/           # handler.go, jwt.go, middleware.go, oauth.go
│   │   └── alert/          # handler.go, service.go, worker.go
│   └── pkg/
│       ├── storage/        # S3 adapter behind StorageProvider interface
│       ├── email/          # SES adapter behind EmailProvider interface
│       ├── cache/          # Redis client singleton
│       ├── broker/         # RabbitMQ publisher behind MessageBroker interface
│       ├── middleware/     # RequestID, Logger, Recovery, RateLimit, JWT, RequireRole
│       └── logger/         # zerolog structured logger
│
├── notification-svc/       # Phase 5 — standalone Go gRPC microservice
├── infra/                  # Terraform — ECS, RDS, ElastiCache, S3, SNS/SQS
├── load-tests/             # k6 scripts
└── docs/adr/               # Architecture Decision Records
```

---

## Error Handling — Non-negotiable rules

Go errors are values. Every error must be returned and handled explicitly at the call site.

```go
// Always wrap errors with context
hotel, err := s.repo.FindByID(ctx, id)
if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrHotelNotFound   // domain sentinel error
    }
    return nil, fmt.Errorf("FindByID: %w", err)  // wrap with operation name
}
```

**Three rules — never break these:**
1. Never use `_` to discard an error unless failure truly cannot affect the program
2. Always wrap with `fmt.Errorf("operation: %w", err)` to build a traceable chain
3. Define domain errors as sentinel values in `types.go`: `var ErrHotelNotFound = errors.New("hotel not found")`

**Expected error chain in logs:**
```
"CreateHotel: InsertImages: s3.Upload: connection timeout"
```

---

## Testing Strategy

### Unit Tests — `internal/*/service.go`
- Test business logic in complete isolation — no DB, no Redis, no HTTP
- Inject mock implementations of every interface dependency
- Use `testify/assert` and `testify/mock`
- **Table-driven tests** for anything with multiple input/output combinations

```go
// Always structure unit tests like this
func TestHotelService_Create_InvalidPrice(t *testing.T) {
    mockRepo    := new(MockHotelRepository)
    mockStorage := new(MockStorageProvider)
    mockBroker  := new(MockMessageBroker)
    svc := hotel.NewService(mockRepo, mockStorage, mockBroker)

    _, err := svc.Create(ctx, hotel.CreateRequest{PricePerNight: -100})

    assert.ErrorIs(t, err, hotel.ErrInvalidPrice)
    mockRepo.AssertNotCalled(t, "Insert")
}
```

### Integration Tests — `internal/*/repository.go`
- Test against a real PostgreSQL instance using `testcontainers-go`
- Verify actual SQL, indexes, constraints, and cursor pagination
- Each test spins up its own container — no shared state

### Handler Tests — `internal/*/handler.go`
- Use `net/http/httptest` — no real server
- Real Gin router, mock service layer
- Test auth enforcement, status codes, response shapes

### Frontend Tests — Vitest + React Testing Library
- Component tests for `HotelCard`, `SearchFilters`
- Hook tests for `useAuth`, `useFavorites` (optimistic update + rollback)
- Integration tests for auth-protected route redirects

---

## Coding Conventions

### Go
- All interfaces defined in `types.go` of the owning domain
- Constructor functions always named `New<Type>` (e.g. `NewHotelService`)
- Return `(result, error)` from every function that can fail
- HTTP handlers always use `c.JSON` — never write directly to `c.Writer`
- Middleware added in `main.go`, never inside domain packages
- Use `context.Context` as the first argument to every service and repo method
- Log with `zerolog` — always include `request_id` and `user_id` fields

### Frontend
- Pages (containers) fetch data via custom hooks
- UI components (presenters) receive data as props — no direct API calls
- All TanStack Query logic lives inside `hooks/` — never inline in components
- Zustand stores in `store/` — one store per concern (auth, filters)
- Named exports only — no default exports for components

---

## Key Design Patterns

| Pattern | Location | Rule |
|---------|----------|------|
| Repository | `internal/*/repository.go` | Services never import sqlc directly |
| Strategy | `internal/search/sort_strategy.go` | New sort = new struct implementing the interface, no switch statements |
| Builder | `internal/search/query_builder.go` | Chain filters step by step, build SQL params at the end |
| Adapter | `pkg/storage/s3.go`, `pkg/email/ses.go` | AWS SDK never imported outside `pkg/` |
| Observer | `internal/alert/worker.go` + RabbitMQ | Publish `HotelCreated` event; consumers are independent |
| Singleton | `pkg/cache/redis.go`, `pkg/logger/logger.go` | One shared instance, initialised once in `main.go` |
| Facade | `internal/search/service.go` | One `Search()` method hides DB query + cache check + filter logic |

---

## Database

### Core Tables
`users`, `roles`, `hotels`, `hotel_images`, `flights`, `locations`, `favorites`, `inquiries`, `price_alerts`, `sessions`, `audit_logs`

### Rules
- All migrations in `backend/db/migrations/` using `golang-migrate` (up/down pairs)
- All SQL queries in `backend/db/queries/` — never write raw SQL outside these files
- Run `sqlc generate` after editing any `.sql` query file
- Use `tsvector` columns for full-text search on hotels and flights
- Cursor-based pagination only — never `OFFSET`

---

## Interfaces to know

```go
// pkg/storage/storage.go
type StorageProvider interface {
    Upload(ctx context.Context, key string, r io.Reader) (string, error)
    Delete(ctx context.Context, key string) error
    PresignedURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}

// pkg/email/email.go
type EmailProvider interface {
    Send(ctx context.Context, to, subject, body string) error
}

// pkg/broker/broker.go
type MessageBroker interface {
    Publish(ctx context.Context, event string, payload any) error
}
```

Services depend on these interfaces — never on the concrete AWS implementations.

---

## Phases

| Phase | Status | Focus |
|-------|--------|-------|
| 1 | Foundation | Monorepo, Docker, DB schema, CI skeleton |
| 2 | Core | Hotel/flight CRUD, auth, search, image upload |
| 3 | Production | AWS deploy, Datadog RUM + APM, full CI/CD |
| 4 | Polish | RabbitMQ migration, k6 load tests, ADRs |
| 5 (optional) | Microservices | notification-svc, gRPC, API Gateway, SNS/SQS |

---

## Rules — Never break these

- Handlers never import the `database/sql` or `pgx` packages
- Repositories never import `net/http` or `github.com/gin-gonic/gin`
- Services never construct HTTP responses or read from `*gin.Context`
- No global variables except the logger and Redis client singletons in `pkg/`
- No `init()` functions — all setup happens explicitly in `main.go`
- Every new domain gets its own directory under `internal/` with all four files
- All AWS SDK usage is behind an interface in `pkg/` — never called directly from `internal/`
