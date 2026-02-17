# Calculator API

A simple calculator REST API built with FastAPI, backed by SQLite and SQLModel.

## Requirements

- Python 3.14+

## Setup

```bash
python3.14 -m venv .venv
source .venv/bin/activate
pip install -e ".[dev]"
```

## Running the Application

```bash
python server.py
```

Or after installing:

```bash
calculator-api
```

Or directly with Uvicorn:

```bash
uvicorn server:app --host 0.0.0.0 --port 8000
```

The server runs on `http://localhost:8000` by default.

## Configuration

Configuration is managed via environment variables with the `CALC_` prefix:

| Variable | Description | Default |
|----------|-------------|---------|
| `CALC_DATABASE_URL` | SQLite database URL | `sqlite:///calculator.db` |
| `CALC_SERVER_HOST` | Server bind host | `0.0.0.0` |
| `CALC_SERVER_PORT` | Server bind port | `8000` |

Example:

```bash
CALC_SERVER_PORT=9000 python server.py
```

## Endpoints

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Returns service status and version |

### Calculator Operations

All calculator endpoints accept POST with JSON body `{"a": <number>, "b": <number>}` and return `{"result": <number>}`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/add` | Add two numbers |
| POST | `/subtract` | Subtract b from a |
| POST | `/multiply` | Multiply two numbers |
| POST | `/divide` | Divide a by b |

Division by zero returns HTTP 400 with `{"detail": "Cannot divide by zero"}`.

### Calculation History (CRUD)

Store calculations in a SQLite database.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/calculations` | Create & store a calculation |
| GET | `/calculations` | List all stored calculations |
| GET | `/calculations/{id}` | Get a specific calculation |
| DELETE | `/calculations/{id}` | Delete a calculation |

## Examples

**Quick calculation:**

```bash
curl -X POST http://localhost:8000/add \
  -H "Content-Type: application/json" \
  -d '{"a": 5, "b": 3}'
```

Response: `{"result": 8.0}`

## Running Tests

```bash
pytest tests/ -v
```

## API Docs

FastAPI auto-generates interactive docs at:
- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc
