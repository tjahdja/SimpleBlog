# SimpleBlog API - REST API Assignment

A production-ready, containerized REST API built with Go (Golang), Gin, GORM, and PostgreSQL.

---

## Architecture Explanation

This project splits code into decoupled layers to ensure strict separation of concerns, high testability, and independence from external frameworks:

1. **Domain / Entity Layer (`internal/entity`)**: Contains models (`User`, `Post`, `Comment`).
2. **Repository Layer (`internal/repository`)**: Data mapping layer/Database Operations.
3. **Service / Usecase Layer (`internal/service`)**: Core application business logic (e.g., password hashing, token generation, validation rules).
4. **Handler / Controller Layer (`internal/handler`)**: Delivery mechanism. Parses JSON, manages Gin contexts, executes input structural validation (`binding` tags), and maps internal errors to clean HTTP status codes.

---

## 🛠️ Technology Choices Justification

* **Go (Golang)**: Chosen for its performance, minimal memory footprint, rapid speed.
* **Gin Gonic Framework**: High-performance HTTP router utilizing a custom Radix tree rendering engine. Used primarily for routing, middleware management, and JSON payload binding validations.
* **GORM & PostgreSQL**: Object-Relational Mapping library paired with a rock-solid relational database. GORM keeps database queries readable while safely managing connection pooling and soft deletes.
* **Docker & Docker Compose**: Ensures the entire backend stacks up identically across any machine with a single unified container runtime execution environment.

---

## 🚀 Setup and Running Instructions

### Prerequisites
* Docker and Docker Compose installed on your host machine.

### Getting Started

1. **Clone and Navigate into the Project Root:**
   ```bash
   cd SimpleBlog
   ```
2. **Build and Run the Application:**
   ```bash
   docker-compose up --build
   ```
3. **Access the Interactive API Documentation:**
   Open your browser and navigate to `http://localhost:8080/swagger/index.html`

## Test Coverage Report

To view the test coverage report, run the following command:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html coverage.out
or
go tool cover -html coverage
```
This will generate an HTML coverage report that you can open in your browser.

## Known Limitations and Future Improvements
1. Stateless Token Invalidation: JWT blocks are currently fully stateless. Introducing a Redis caching layer in future iterations would allow temporary blacklisting or immediate server-side revocation of tokens upon explicit user logout.
2. Refresh Token Rotation: Currently relies on a singular short-lived access token configuration layout. Implementing a rotation strategy with a secondary long-lived refresh token would maximize API safety boundaries.
