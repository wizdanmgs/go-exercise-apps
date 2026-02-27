# Simple Notes Service

A production-ready REST API for managing notes, built with Go using Clean Architecture principles.

This project demonstrates:

- net/http fundamentals
- chi router
- Clean Architecture layering
- DTO separation
- Custom domain errors
- Structured logging (slog)
- Unit tests (mock-based TDD)
- Integration tests (real HTTP stack)
- Graceful shutdown

---

## Architecture

The project follows Clean Architecture:

.
├── cmd/
│ └── main.go
└── internal/
├── delivery/
│ └── http/ -> HTTP handlers, DTOs, routing
├── domain/ -> Business entities & domain errors  
 ├── usecase/ -> Application business logic  
 └── respository/
└── memory/ -> Repository implementation (in-memory)

Dependency flow:

Delivery (HTTP / chi)  
↓  
Usecase (Business Logic)  
↓  
Domain (Interfaces & Entities)  
↑  
Repository (Memory Repository)

Domain does not depend on HTTP, router, or infrastructure.

---

## Features

- Create note
- Get all notes
- Get note by ID
- Update note
- Delete note
- In-memory storage
- JSON responses
- Custom error handling
- Structured logging
- Graceful shutdown

---

## Tech Stack

- Go
- github.com/go-chi/chi/v5 (HTTP router)
- log/slog (structured logging)
- net/http
- httptest (integration testing)

---

## API Endpoints

Base URL:

http://localhost:8080

### Create Note

POST `/notes/`

Request:

```json
{
  "id": "1",
  "title": "My Note"
}
```

Response: 201 Created

```json
{
  "id": "1",
  "title": "My Note"
}
```

---

### Get All Notes

GET `/notes/`

Response: 200 OK

```json
[
  {
    "id": "1",
    "title": "My Note"
  }
]
```

---

### Get Note By ID

GET `/notes/{id}/`

Response: 200 OK

```json
{
  "id": "1",
  "title": "My Note"
}
```

If not found:

```json
{
  "error": "not found"
}
```

---

### Update Note

PUT `/notes/{id}/`

Request:

```json
{
  "title": "Updated Title"
}
```

Response: 200 OK

```json
{
  "id": "1",
  "title": "Updated Title"
}
```

---

### Delete Note

DELETE `/notes/{id}/`

Response: 200 OK

```json
{
  "message": "deleted"
}
```

---

## Error Handling

Custom domain errors:

- ErrInvalidInput
- ErrNotFound

Errors are mapped to proper HTTP status codes:
| Error | HTTP Status |
|--|--|
| ErrInvalidInput | 400 |
| ErrNotFound | 404 |
| Other errors | 500 |

---

## Running the Application

### 1. Install dependencies

```sh
go mod tidy
```

### 2. Run the server

```sh
go run cmd/main.go
```

Server runs on:

:8080

---

## Graceful Shutdown

The server supports graceful shutdown:

- Listens for SIGINT / SIGTERM
- Stops accepting new requests
- Waits up to 5 seconds for in-flight requests
- Exits cleanly

---

## Testing

### Run all tests

```sh
go test ./...
```

### Run with race detector

```sh
go test -race ./...
```

### Run with coverage

```sh
go test -cover ./...
```

---

## Testing Strategy

### Unit Tests

- Use table-driven tests
- Mock repository layer
- Validate business logic independently

### Integration Tests

- Use httptest.NewServer
- Test full stack:
  - HTTP
  - Handler
  - Usecase
  - Repository
- Validate JSON contract

---

## Design Decisions

### DTO Layer

HTTP layer uses DTOs:

- CreateNoteRequest
- UpdateNoteRequest
- NoteResponse

Domain models are never exposed directly via JSON.

This keeps transport and business logic separated.

---

### Repository Pattern

Repository is defined as an interface in domain:

```go
type  NoteRepository  interface {
 Create(Note) error
 GetAll() ([]Note, error)
 GetByID(string) (Note, error)
 Update(string, Note) error
 Delete(string) error
}
```

In-memory implementation lives in infrastructure layer.

Swapping to PostgreSQL or another storage only requires replacing repository implementation.
