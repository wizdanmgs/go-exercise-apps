# Go Concurrent Web Scraper

A production-grade concurrent web scraper written in Go.

This project demonstrates clean architecture principles and implements:

- Worker pool concurrency
- Global and per-domain rate limiting
- Robots.txt compliance
- Exponential backoff retry with jitter
- Circuit breaker per domain
- Integration and unit tests

---

## Features

### Concurrent Crawling

- Configurable worker pool
- Context-aware graceful shutdown

### Rate Limiting

- Global RPS limit
- Per-domain RPS limit
- Burst configuration

### Robots.txt Compliance

- Caches robots.txt per domain
- Supports wildcard fallback
- Honors Crawl-Delay
- Proper port handling (supports `httptest` servers)

### Retry Strategy

- Exponential backoff
- Jitter to prevent thundering herd
- Retries on:
  - HTTP 429
  - HTTP 5xx
  - Network timeouts

### Circuit Breaker

- Per-domain breaker
- Opens after configurable failure threshold
- Auto reset after timeout

### Testing

- Table-driven unit tests
- Full integration tests using `httptest`
- Race detector compatible

---

## Project Structure

scraper/  
├── cmd/  
├── internal/  
│ ├── domain/  
│ └── usecase/  
│ ├── scraper_usecase.go  
│ ├── circuit_breaker.go  
│ └── scraper_integration_test.go  
├── go.mod  
└── README.md

---

## Installation

```bash
git clone <repository-url>
cd scraper
go mod tidy
```

---

## Running Tests

Run all tests:

```bash
go test ./...
```

Run with race detector:

```bash
go test -race ./...
```

Run verbose:

```bash
go test -v ./...
```

## Configuration Parameters

| Parameter        | Description                        |
| ---------------- | ---------------------------------- |
| workerCount      | Number of concurrent workers       |
| globalRPS        | Global requests per second         |
| globalBurst      | Global burst capacity              |
| domainRPS        | Per-domain requests per second     |
| domainBurst      | Per-domain burst capacity          |
| maxRetries       | Maximum retry attempts             |
| baseDelay        | Base delay for exponential backoff |
| breakerThreshold | Failures before circuit opens      |
| breakerTimeout   | Time before circuit resets         |
