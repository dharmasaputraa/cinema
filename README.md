# Cinema API

A RESTful API for cinema management system built with Go and Clean Architecture principles.

## Description

Cinema API is a backend service for managing cinema operations, including user authentication, cinema details, screen management, and seat booking. The project follows Clean Architecture principles with a focus on maintainability, scalability, and testability.

## Features

- User Authentication: JWT-based authentication system
- Cinema Management: Manage cinema information
- Screen Management: Handle multiple screens per cinema
- Seat Management: Manage seat configurations and availability
- Graceful Shutdown: Proper handling of server shutdown
- Auto-migration: Automatic database migrations
- Request Tracing: Request ID middleware for tracking
- CORS Support: Configurable CORS origins
- Structured Logging: Zap logger with production/development modes

## Tech Stack & Dependencies

### Core Technologies

- Go 1.26.1: Programming language
- Gin: HTTP web framework
- GORM: ORM for database operations
- PostgreSQL: Primary database
- Redis: Caching layer
- Docker & Docker Compose: Containerization

### Key Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `gorm.io/gorm` - ORM
- `gorm.io/driver/postgres` - PostgreSQL driver
- `github.com/redis/go-redis/v9` - Redis client
- `go.uber.org/zap` - Structured logging
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/golang-migrate/migrate/v4` - Database migrations
- `github.com/google/wire` - Dependency injection
- `github.com/spf13/viper` - Configuration management
- `github.com/gin-contrib/cors` - CORS middleware

## Architecture

The project follows Clean Architecture with the following structure:

```
.
├── cmd/                    # Application entry points
│   └── main.go            # Main application
├── internal/              # Private application code
│   ├── auth/              # Authentication module
│   │   ├── delivery/      # HTTP handlers
│   │   ├── domain/        # Domain models
│   │   ├── repository/    # Data access
│   │   └── usecase/       # Business logic
│   ├── cinema/            # Cinema management module
│   │   ├── delivery/      # HTTP handlers
│   │   ├── domain/        # Domain models
│   │   ├── repository/    # Data access
│   │   └── usecase/       # Business logic
│   ├── user/              # User management module
│   │   ├── delivery/      # HTTP handlers
│   │   ├── domain/        # Domain models
│   │   ├── repository/    # Data access
│   │   └── usecase/       # Business logic
│   └── infrastructure/    # Infrastructure layer
│       ├── cache/         # Redis setup
│       ├── config/        # Configuration
│       └── db/            # Database setup & migrations
└── pkg/                   # Public/reusable packages
    ├── errors/           # Error handling
    ├── helper/           # Helper functions
    ├── middleware/       # HTTP middleware
    ├── pagination/       # Pagination utilities
    ├── response/         # Response formatting
    └── validator/        # Validation utilities
```

### Architecture Layers

1. Delivery Layer (`delivery/`): HTTP handlers that receive requests and send responses
2. Use Case Layer (`usecase/`): Business logic and application rules
3. Repository Layer (`repository/`): Data access and database operations
4. Domain Layer (`domain/`): Core business entities and interfaces
5. Infrastructure Layer (`infrastructure/`): External services (DB, Cache, Config)
6. Package Layer (`pkg/`): Reusable utilities across modules

### Middleware Stack

The application uses the following middleware in order:

1. RequestID: Adds unique request ID for tracing
2. CORS: Handles cross-origin requests
3. Logger: Logs HTTP requests
4. ErrorHandler: Centralized error handling

## Prerequisites

Before running this project, ensure you have:

- Go 1.26.1 or higher - [Download Go](https://golang.org/dl/)
- Docker & Docker Compose - [Download Docker](https://www.docker.com/get-started)
- PostgreSQL (if not using Docker)
- Redis (if not using Docker)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/dharmasaputraa/cinema.git
cd cinema/services/api
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Install Air (Hot Reload)

```bash
go install github.com/cosmtrek/air@latest
```

### 4. Setup Environment Variables

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cinema_db
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=your-secret-key-here
JWT_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_HOURS=168
```

## How to Run

### Development Mode (Recommended)

Using Docker Compose for dependencies and Air for hot reload:

```bash
make dev
```

This command will:

- Start PostgreSQL and Redis using Docker Compose
- Start the API server with hot reload

### Stop Development Environment

```bash
make dev-down
```

### Reset Development Environment (includes DB)

```bash
make dev-reset
```

### Production Mode

```bash
make prod
```

This command will:

- Build and start all services using Docker Compose
- Run in production mode

### Stop Production

```bash
make prod-down
```

### Manual Run (Without Docker)

1. Start PostgreSQL and Redis manually
2. Run the application:

```bash
go run cmd/main.go
```

### Build Binary

```bash
go build -o cinema-api cmd/main.go
./cinema-api
```

## API Documentation

### Base URL

```
http://localhost:8080
```

### Health Check

```bash
GET /health
```

Response:

```json
{
  "status": "ok",
  "request_id": "uuid-here"
}
```

### API Endpoints

All API endpoints are prefixed with `/api/v1`:

- Cinema Module
  - `GET /api/v1/cinemas` - Get all cinemas
  - `GET /api/v1/cinemas/:id` - Get cinema by ID
  - `POST /api/v1/cinemas` - Create new cinema
  - `PUT /api/v1/cinemas/:id` - Update cinema
  - `DELETE /api/v1/cinemas/:id` - Delete cinema

- Screen Module
  - `GET /api/v1/cinemas/:cinema_id/screens` - Get all screens for a cinema
  - `GET /api/v1/screens/:id` - Get screen by ID
  - `POST /api/v1/screens` - Create new screen
  - `PUT /api/v1/screens/:id` - Update screen
  - `DELETE /api/v1/screens/:id` - Delete screen

- Seat Module
  - `GET /api/v1/screens/:screen_id/seats` - Get all seats for a screen
  - `GET /api/v1/seats/:id` - Get seat by ID
  - `POST /api/v1/seats` - Create new seat
  - `PUT /api/v1/seats/:id` - Update seat
  - `DELETE /api/v1/seats/:id` - Delete seat

- Auth Module
  - `POST /api/v1/auth/register` - Register new user
  - `POST /api/v1/auth/login` - User login
  - `POST /api/v1/auth/refresh` - Refresh access token

## Environment Variables

| Variable                 | Description                                      | Default              |
| ------------------------ | ------------------------------------------------ | -------------------- |
| APP_ENV                  | Application environment (development/production) | development          |
| APP_PORT                 | Server port                                      | 8080                 |
| DB_HOST                  | Database host                                    | localhost            |
| DB_PORT                  | Database port                                    | 5432                 |
| DB_USER                  | Database username                                | postgres             |
| DB_PASSWORD              | Database password                                | postgres             |
| DB_NAME                  | Database name                                    | cinema_db            |
| DB_SSLMODE               | Database SSL mode                                | disable              |
| REDIS_HOST               | Redis host                                       | localhost            |
| REDIS_PORT               | Redis port                                       | 6379                 |
| REDIS_PASSWORD           | Redis password                                   | (empty)              |
| REDIS_DB                 | Redis database number                            | 0                    |
| JWT_SECRET               | JWT secret key                                   | your-secret-key-here |
| JWT_EXPIRY_HOURS         | JWT token expiry in hours                        | 24                   |
| JWT_REFRESH_EXPIRY_HOURS | JWT refresh token expiry in hours                | 168                  |

## Project Structure

```
cinema-api/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── auth/                   # Authentication module
│   │   ├── delivery/
│   │   │   └── http/
│   │   ├── domain/
│   │   ├── repository/
│   │   │   └── postgres/
│   │   └── usecase/
│   ├── cinema/                 # Cinema management module
│   │   ├── delivery/
│   │   │   └── http/
│   │   ├── domain/
│   │   ├── repository/
│   │   │   └── postgres/
│   │   └── usecase/
│   ├── user/                   # User management module
│   │   ├── delivery/
│   │   ├── domain/
│   │   ├── repository/
│   │   └── usecase/
│   └── infrastructure/
│       ├── cache/              # Redis cache
│       ├── config/             # Configuration
│       └── db/                 # Database & migrations
├── pkg/
│   ├── errors/                 # Error handling
│   ├── helper/                 # Helper functions
│   ├── middleware/             # HTTP middleware
│   │   ├── cors.go
│   │   ├── error_handler.go
│   │   ├── logger.go
│   │   └── request_id.go
│   ├── pagination/             # Pagination utilities
│   ├── response/               # Response formatting
│   └── validator/              # Validation utilities
├── .env                        # Environment variables (not in git)
├── .env.example                # Example environment variables
├── .gitignore                  # Git ignore rules
├── .air.toml                   # Air hot reload configuration
├── docker-compose.yml          # Production Docker Compose
├── docker-compose.dev.yml      # Development Docker Compose
├── Dockerfile                  # Docker image
├── go.mod                      # Go module definition
├── go.sum                      # Go dependencies checksum
├── Makefile                    # Build and run commands
└── README.md                   # This file
```

## Development Commands

### Using Makefile

```bash
# Development
make dev              # Start development environment with hot reload
make dev-down         # Stop development environment
make dev-reset        # Reset development environment (includes DB)

# Production
make prod             # Start production environment
make prod-down        # Stop production environment

# Utilities
make logs             # View Docker logs
make ps               # View running containers
```

### Manual Commands

```bash
# Install dependencies
go mod download
go mod tidy

# Run with hot reload (Air)
air

# Run without hot reload
go run cmd/main.go

# Build binary
go build -o cinema-api cmd/main.go

# Run tests
go test ./...

# View Docker logs
docker compose logs -f

# View running containers
docker compose ps
```

## Key Features Implementation

- Clean Architecture: Separation of concerns with clear boundaries
- Dependency Injection: Manual wiring of dependencies in `main.go`
- Middleware Stack: Request tracking, CORS, logging, error handling
- Graceful Shutdown: Proper cleanup on SIGINT/SIGTERM
- Auto-migration: Automatic database schema updates
- Environment-based Configuration: Different configs for dev/prod
- Structured Logging: Zap logger with production/development modes
- Error Handling: Centralized error handling with consistent responses
- Validation: Request validation using go-playground/validator
- Pagination: Built-in pagination support for list endpoints

## Notes

- The application uses `golang-migrate` for database migrations
- Auto-migration can be disabled by setting `AUTO_MIGRATE=false` in environment
- The server supports graceful shutdown with a 5-second timeout
- All HTTP responses include a `request_id` for tracing
- CORS origins are configurable via `CORS_ORIGINS` environment variable

## Debugging

### Check Logs

```bash
# Development logs
make logs

# Or view specific service logs
docker compose logs -f postgres
docker compose logs -f redis
```
