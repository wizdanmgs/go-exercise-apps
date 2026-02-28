# File Upload Server

A production-style file upload server written in Go.

This project demonstrates:

- Multipart form handling
- Image validation via decoding (`image.DecodeConfig`)
- Clean architecture structure
- Chi router usage
- Static file serving
- Table-driven unit tests
- Integration tests
- Fuzz testing

---

## Features

- Upload image via `multipart/form-data`
- Validate image by decoding (not just MIME sniffing)
- Restrict allowed formats (JPEG, PNG, GIF)
- Limit image dimensions
- Save uploaded files to local directory
- Serve uploaded files via HTTP
- Table-driven unit tests
- Integration tests using `httptest`
- Fuzz testing for image validator

---

## Project Structure

file-upload-server/  
├── cmd/  
│ └── server/  
│ └── main.go  
├── internal/  
│ ├── handler/  
│ ├── service/  
│ ├── validator/  
│ └── server/  
└── uploads/

### Architecture

- `handler` → HTTP layer
- `service` → business logic
- `validator` → image validation logic
- `server` → router setup
- `cmd/server` → application entry point

---

## Installation

Clone the repository:

```bash
git clone <repository-url>
cd file-upload-server
```

Install dependencies

```bash
go mod tidy
```

---

## Running the Server

```bash
go run ./cmd/server
```

Server runs on:

```http
http://localhost:8080
```

---

## Uploading an Image

Example using curl:

```bash
curl  -X POST http://localhost:8080/api/upload \
  -F  "image=@test.png"
```

Successful response:

```json
{
  "message": "upload successful",
  "filename": "example.png",
  "url": "/uploads/example.png"
}
```

---

## Accessing Uploaded Files

Uploaded files are served statically:

```http
http://localhost:8080/uploads/<filename>
```

---

## Image Validation Strategy

Images are validated using:

```go
image.DecodeConfig(file)
```

This ensures:

- The file is a real decodable image
- The format is allowed (jpeg, png, gif)
- The dimensions are valid
- The dimensions are within allowed limits

This approach is stronger than checking file extensions or MIME headers alone.

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

Generate coverage report:

```bash
go test -cover ./...
```

---

## Running Fuzz Tests

Run fuzzing for the validator:

```bash
go test ./internal/validator -fuzz=FuzzValidateImage
```

Stop fuzzing with:

```bash
CTRL + C
```

If a failure is found, Go will automatically save the failing input under:

```bash
testdata/fuzz/
```

---

## Security Considerations

- Images are decoded to verify authenticity
- Only allowed formats are accepted
- Maximum image dimensions enforced
- File system writes are isolated during tests
- Fuzz testing protects against malformed input
