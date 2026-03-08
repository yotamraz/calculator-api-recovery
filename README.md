# Arithmetic Service

A lightweight Go service providing stateless arithmetic operations via a REST API. Built with [Gin](https://github.com/gin-gonic/gin) v1.10.

## Setup

### Prerequisites

- Go 1.23+

### Install Dependencies

```bash
go mod download
```

## Run

```bash
go run .
```

The server starts on `http://0.0.0.0:8080` by default.

### Configuration

Configuration is managed via environment variables:

| Variable | Default   | Description          |
|----------|-----------|----------------------|
| `HOST`   | `0.0.0.0` | Server bind address  |
| `PORT`   | `8080`    | Server listen port   |

Copy `.env.example` to `.env` and modify as needed:

```bash
cp .env.example .env
```

## API Endpoints

### Health Check

```
GET /health
```

Response:
```json
{"status": "ok", "version": "0.1.0"}
```

### Arithmetic Operations

All arithmetic endpoints accept `POST` with JSON body `{"a": <number>, "b": <number>}` and return `{"result": <number>}`.

| Endpoint     | Description        |
|--------------|--------------------|
| `POST /add`      | Add two numbers    |
| `POST /subtract` | Subtract b from a  |
| `POST /multiply` | Multiply two numbers |
| `POST /divide`   | Divide a by b      |

### Examples

**Addition:**

```bash
curl -X POST http://localhost:8080/add \
  -H "Content-Type: application/json" \
  -d '{"a": 5, "b": 3}'
```

Response: `{"result": 8}`

**Division:**

```bash
curl -X POST http://localhost:8080/divide \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 3}'
```

Response: `{"result": 3.3333333333333335}`

**Division by zero:**

```bash
curl -X POST http://localhost:8080/divide \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 0}'
```

Response (HTTP 400): `{"detail": "Cannot divide by zero"}`

## Testing

Run all tests:

```bash
go test -v ./...
```

Run only unit tests:

```bash
go test -v -run "Test(Add|Subtract|Multiply|Divide)$" ./...
```

Run only HTTP integration tests:

```bash
go test -v -run "TestHandler" ./...
```
