# URL Shortener (In-Memory)

This project demonstrates:

- Clean layered architecture
- In-memory storage using `map`
- Interface-based dependency injection
- Functional options pattern
- Table-driven unit tests
- HTTP integration tests
- Graceful shutdown with context
- Chi router

---

## Features

- Generate short codes
- Redirect to original URL
- Expiry (TTL) support
- Collision-safe code generation
- Graceful server shutdown
- Full test coverage (unit + integration)

---

## Project Structure

url-shortener/  
├── cmd/  
│ └── server/  
│ └── main.go  
├── internal/  
│ ├── handler/  
│ │ ├── http.go  
│ │ ├── http_integration_test.go  
│ │ └── shutdown_integration_test.go  
│ ├── model/  
│ │ └── url.go  
│ ├── service/  
│ │ ├── shortener.go  
│ │ └── shortener_test.go  
│ └── store/  
│ ├── store.go  
│ └── memory.go  
├── go.mod  
└── README.md

---

## Architecture

HTTP Handler  
↓  
Service Layer  
↓  
URLStore (interface)  
↓  
MemoryStore (implementation)

### Key Principles

- Accept interfaces, return structs
- Constructor-based dependency injection
- Functional options for extensibility
- No global state
- Thread-safe map using `sync.RWMutex`

---

## Installation

```bash
git clone <your-repo-url>
cd urlshortener
go mod tidy
```

---

## Run the Server

```bash
go run ./cmd/server
```

Server runs on:

http://localhost:8080

---

## API Usage

### Create Short URL

**POST** `/shorten`

Request body:

```json
{
  "url": "https://example.com",
  "ttl": 3600
}
```

Response:

```json
{
  "code": "Ab3kLm9P"
}
```

---

### Redirect

**GET** `/{code}`

Example:

GET /Ab3kLm9P

Returns:

302 Found  
Location: https://example.com

---

## Expiry Behavior

- TTL is defined in seconds
- Expired URLs return `404`
- Expired entries are lazily deleted on access

---

## Testing

### Run All Tests

```bash
go test ./...
```

### Run With Race Detector

```bash
go test -race ./...
```

### Run With Coverage

```bash
go test -cover ./...
```

---

## Test Types

### Unit Tests

- Table-driven tests
- Service layer behavior
- Collision handling
- Expiry validation

### Integration Tests

- Real HTTP server (`httptest`)
- Full request lifecycle
- Redirect verification
- Expiry behavior
- Graceful shutdown testing

---

## Graceful Shutdown

The server:

- Listens for `SIGINT` and `SIGTERM`
- Stops accepting new connections
- Waits for active requests to finish
- Uses `context.WithTimeout` for shutdown control

---
