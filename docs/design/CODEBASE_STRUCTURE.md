backend/
│── cmd/                     # Main application entry points
│   ├── main.go              # Main entry point for the monolith
│   ├── migrate.go           # Entry point for DB migrations
│
├── configs/                 # Configuration files
│   ├── config.yaml          # Application configurations
│   ├── config.go            # Loads & parses config
│
├── internal/                # Core business logic (monolithic modules, future microservices)
│   ├── auth/                # Auth Module (Example)
│   │   ├── handler.go       # HTTP handlers (controllers)
│   │   ├── service.go       # Business logic
│   │   ├── repository.go    # Database access layer
│   │   ├── model.go         # Data structures
│   │   ├── routes.go        # Register routes
│   │   ├── proto/           # (Future gRPC) .proto definitions for microservices
│   │   ├── events.go        # (Future Kafka) Event publishing/subscription
│   │
│   ├── problem/             # Problem Module
│   ├── challenge/           # Challenge Module
│   ├── rating/              # Rating Module
│   ├── history/             # User Progress & History Module
│
├── pkg/                     # Shared utilities (Reusable for all services)
│   ├── db/                  # Database setup & connection
│   │   ├── postgres.go      # PostgreSQL connection
│   │   ├── redis.go         # Redis connection
│   │
│   ├── jwt/                 # JWT authentication utilities
│   ├── logger/              # Structured logging (Zap)
│   ├── httpserver/          # HTTP server setup
│   ├── cache/               # Caching utilities (Redis)
│   ├── events/              # Event Bus (Kafka, Pub/Sub, future async event handling)
│
├── api/                     # API Definitions (REST + Future gRPC)
│   ├── openapi.yaml         # OpenAPI spec for REST APIs
│   ├── proto/               # Proto files for gRPC (future-proofing)
│
├── migrations/              # Database migration files (SQL)
│   ├── 001_create_users.sql
│   ├── 002_create_problems.sql
│
├── tests/                   # Integration & unit tests
│   ├── auth_test.go
│   ├── problem_test.go
│
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency lock file
