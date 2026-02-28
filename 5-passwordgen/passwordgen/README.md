# Password Generator CLI

A secure and idiomatic CLI password generator written in Go.

This tool generates cryptographically secure random passwords using Go's `crypto/rand` package and provides a clean command-line interface powered by Cobra.

---

## Features

- Cryptographically secure randomness (`crypto/rand`)
- Custom password length
- Optional symbol inclusion
- Short and long CLI flags
- Table-driven unit tests
- Clean and idiomatic Go implementation

---

## Installation

### Clone the repository

```bash
git clone <this-repository>
cd passwordgen
```

### Build the binary

```bash
go build -o passwordgen
```

### Run directly without building

```bash
go run main.go
```

---

## Usage

### Default (16 characters, symbols enabled)

```bash
./passwordgen
```

### Custom length

```bash
./passwordgen --length 24
```

Short flag:

```bash
./passwordgen -l 24
```

### Disable symbols

```bash
./passwordgen --symbols=false
```

Short flag:

```bash
./passwordgen -s=false
```

### Combined example

```bash
./passwordgen -l 32 -s=false
```

---

## Flags

| Flag        | Short | Default | Description                      |
| ----------- | ----- | ------- | -------------------------------- |
| `--length`  | `-l`  | 16      | Length of the generated password |
| `--symbols` | `-s`  | true    | Include symbols in the password  |

---

## Example Output

```
dG7@pL9!xQ2#vZ8m
```

---

## Project Structure

```
passwordgen/
├── cmd/
│   ├── root.go
│   └── root_test.go
├── main.go
└── go.mod
```

---

## Testing

Run all tests:

```bash
go test ./...
```

Run with coverage:

```bash
go test -cover ./...
```

---

## Security Notes

- Uses `crypto/rand` for secure randomness.
- Does not rely on `math/rand`.
- Suitable for generating passwords for production use.
- Randomness is derived from the system's secure entropy source.
