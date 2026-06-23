# Project overview
Building a robust, production-ready REST API for a simple blog

# Technologies & Preferences
- Language: Go (1.26.2)
- HTTP Router: Gin
- Database: PostgreSQL via GORM
- Architecture: Clean Architecture

# Folder Structure
Follow this strict Clean Architecture folder structure:
├── cmd/
│   └── api/
│       └── main.go         # Entry point (initializes DB, router, and wires layers)
├── internal/
│   ├── entity/             # 1. Domain Models & Structs (Pure data definitions)
│   ├── handler/            # 2. HTTP Layer (Gin context, JSON binding, request validation)
│   ├── service/            # 3. Business Logic Layer (Use-cases, validations, rules)
│   └── repository/         # 4. Data Access Layer (Database operations/queries)
├── pkg/                    # Shared utilities (middleware, logger, custom errors)
├── go.mod
└── instructions.md

## Code Style & Go Best Practices
- **No Over-Engineering:** Keep logic clean.
- **Explicit Error Handling:** Check and handle ALL errors explicitly (`if err != nil`). Never let handlers swallow errors or panic under normal operations.
- **Strict Layer Separation:** 
  - Handlers deal *only* with HTTP status codes and JSON binding.
  - Services deal *only* with business logic and receive/return pure Go types or interfaces.
  - Repositories deal *only* with database interactions.
- **Response Consistency:** Format all JSON API responses utilizing a consistent structured wrapper (e.g., success status, data payload, or error message).
- **Structured Logging:** Use the native `log/slog` package from the standard library for all application logging. 
  - Initialize a global JSON handler logger in `main.go` so logs are machine-readable.
  - Do not use raw `fmt.Println` or standard `log` package.
  - Log errors at `slog.Error` level with contextual attributes (e.g., `slog.String("err", err.Error())`).

## API Core Requirements
1. **Global Middleware:** 
   - Request/Response logger
   - Panic recovery
2. **Public Endpoints:**
   - POST /users/register -> Register an account
   - POST /users/login -> Authenticate and return a JWT
   - GET /posts -> List all blog posts
   - GET /posts/:id -> View a single blog post
3. **Protected Endpoints (Requires JWT Auth Middleware):**
   - POST /posts -> Create a new post (Set `userID` from context as AuthorID)
   - PUT /posts/:id -> Update a post (Service layer checks if `userID` == `AuthorID`)
   - DELETE /posts/:id -> Delete a post (Service layer checks if `userID` == `AuthorID`)


## Future Deliverables (To be implemented AFTER core functionality is verified)
- [Phase 1.5] Native Go table-driven unit tests for critical Service Layer logic.
- [Phase 2] OpenAPI documentation annotations on handler functions.
- [Phase 3] Multi-stage Dockerfile and docker-compose orchestration setup.
