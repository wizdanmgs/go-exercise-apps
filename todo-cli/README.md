# Todo CLI

A simple Todo CLI application written in Go using Clean Architecture principles and built with Cobra.

## Features

- Add a task
- List all tasks
- Mark task as done
- Delete task
- Auto-renumber IDs after delete
- Persist tasks to `tasks.json`
- Clean Architecture structure (domain, usecase, repository, delivery)

---

## Project Structure

todo-cli/
│
├── cmd/
│ ├── root.go
│ ├── add.go
│ ├── list.go
│ ├── delete.go
│ ├── done.go
│ └── main.go
│
├── internal/
│ ├── domain/
│ │ └── task.go
│ │
│ ├── repository/
│ │ ├── task_repository.go
│ │ └── json_repository.go
│ │
│ └── usecase/
│ └── task_usecase.go
│
└── tasks.json

---

## Architecture Overview

The project follows Clean Architecture principles:

- **Domain Layer**
  - Contains business entities.
  - No external dependencies.

- **Usecase Layer**
  - Contains business logic.
  - Depends only on interfaces.

- **Repository Layer**
  - Handles data persistence.
  - Implements repository interfaces.

- **Delivery Layer (CLI)**
  - Uses Cobra for command handling.
  - Calls usecase layer.

Dependency direction:

CLI → Usecase → Repository Interface → Repository Implementation

---

## Installation

Clone the repository:

```bash
git clone <your-repo-url>
cd todo-cli
```

Install dependencies:

```bash
go mod tidy
```

Build the binary:

```bash
go build -o todo
```

---

## Usage

### Add a Task

```bash
./todo add "Learn Go Clean Architecture"
```

### List Tasks

```bash
./todo list
```

Output example:

[ ] 1: Learn Go Clean Architecture
[✓] 2: Build CLI App

### Mark Task as Done

```bash
./todo done 1
```

### Delete Task

```bash
./todo delete 1
```

After deletion, IDs are automatically renumbered.

---

## Example tasks.json

```json
[
  {
    "id": 1,
    "name": "Learn Go Clean Architecture",
    "done": false
  }
]
```
