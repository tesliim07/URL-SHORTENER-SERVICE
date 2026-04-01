# URL-SHORTENER-SERVICE
A REST API built in Go that shortens URLs, stores them in PostgreSQL, and caches lookups in Redis.

## Architecture
url-shortener-service/
├── cmd/
│   └── main.go                  # Entry point — wires everything together
├── config/
│   └── config.go                # Reads environment variables
├── internal/
│   ├── handler/
│   │   └── handler.go           # HTTP handlers — request/response layer
│   ├── service/
│   │   └── service.go           # Business logic — shorten and redirect
│   ├── repository/
│   │   └── repository.go        # PostgreSQL queries
│   └── cache/
│       └── cache.go             # Redis cache logic
├── migrations/
│   └── 001_initial.sql          # Database migration
├── docs/                        # Auto-generated Swagger docs
├── docker-compose.yml
├── .env                         # Environment variables (not committed)
└── go.mod


## API Endpoints
 - `POST/shorten` : Takes a long URL, returns a short URLGET/{code}
    Content-Type: application/json
    {
      "url": "https://www.google.com"
    }
    Response:
    {
      "short_url": "http://localhost:8080/abc123"
    }
 - `GET /{code}` : Redirects the client to the original URL with a 302 

## Prerequisites
Make sure you have the following installed before running the project:
 - Go 1.23+
 - VS Code
 - Docker Desktop

## Getting Started
1. Clone the repository
 - git clone https://github.com/yourname/url-shortener-service
 - cd url-shortener-service
2. Create your .env file
 - DB_HOST=localhost
 - DB_PORT=5432
 - DB_USER=your_db_user
 - DB_PASSWORD=your_db_password
 - DB_NAME=your_db_name

 - REDIS_HOST=localhost
 - REDIS_PORT=6379

 - PGADMIN_DEFAULT_EMAIL=admin@localhost.com
 - PGADMIN_DEFAULT_PASSWORD=admin
3. Start PostgreSQL and Redis
 - docker-compose up -d
4. Install dependencies
 - go mod tidy
5. Generate Swagger docs
 - swag init -g cmd/main.go
6. Run the app
 - go run cmd/main.go
 - The server starts on http://localhost:8080.

## Docker Commands

### Start services
docker-compose up -d

### Stop services
docker-compose down

### Stop and delete all data
docker-compose down -v

### View logs
docker-compose logs -f
