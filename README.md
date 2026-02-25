# Calculator API

A simple REST API server for arithmetic operations with persistent calculation history.

## Setup

```bash
# Create and activate a virtual environment
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate

# Install the package with dev dependencies
pip install -e ".[dev]"

# (Optional) Create a .env file from the example
cp .env.example .env
```

## Run

Start the server with auto-reload for development:

```bash
uvicorn calculator_api.main:app --reload
```

Or using the console script:

```bash
calculator-api
```

The server runs on `http://localhost:8000` by default.

## Configuration

Configuration is managed via environment variables (or a `.env` file):

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `sqlite:///calculator.db` | SQLite database connection URL |
| `HOST` | `0.0.0.0` | Server bind address |
| `PORT` | `8000` | Server port |

## Endpoints

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Service health status |

### Calculator Operations

All calculator endpoints accept POST with JSON body `{"a": <number>, "b": <number>}`.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/add` | Add two numbers |
| POST | `/subtract` | Subtract b from a |
| POST | `/multiply` | Multiply two numbers |
| POST | `/divide` | Divide a by b |

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

**Save a calculation to the database:**

```bash
curl -X POST http://localhost:8000/calculations \
  -H "Content-Type: application/json" \
  -d '{"operation": "mul", "a": 7, "b": 6}'
```

Response: `{"operation": "mul", "a": 7.0, "b": 6.0, "result": 42.0, "id": 1, "created_at": "..."}`

**List all saved calculations:**

```bash
curl http://localhost:8000/calculations
```

**Delete a calculation:**

```bash
curl -X DELETE http://localhost:8000/calculations/1
```

## Testing

```bash
pytest
```

## API Docs

FastAPI auto-generates docs at:
- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc

## Project Structure

```
calculator_api/
  __init__.py
  main.py          # FastAPI app creation, lifespan events, router inclusion
  models.py        # SQLModel models and Pydantic schemas
  routes/
    __init__.py
    calculator.py  # /add, /subtract, /multiply, /divide endpoints
    calculations.py# /calculations CRUD endpoints
    health.py      # /health endpoint
  core.py          # Pure arithmetic business logic
  database.py      # Engine and session dependency
  config.py        # pydantic-settings configuration

tests/
  __init__.py
  test_core.py
  test_calculator_routes.py
  test_calculations_routes.py
  test_health_route.py
```
