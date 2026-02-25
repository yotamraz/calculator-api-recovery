# Calculator API

A simple REST API server for arithmetic operations with persistent calculation history.

## Project Structure

```
calculator_api/
  __init__.py
  main.py          # FastAPI app creation, lifespan events, include routers
  models.py        # SQLModel and Pydantic models
  routes/
    __init__.py
    calculator.py  # /add, /subtract, /multiply, /divide endpoints
    health.py      # /health endpoint
  core.py          # Pure arithmetic business logic
  database.py      # Engine and get_session dependency
  config.py        # pydantic-settings Settings class
```

## Setup

1. **Create and activate a virtual environment:**

```bash
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
```

2. **Install dependencies:**

```bash
pip install -e .
```

For development (includes pytest and httpx):

```bash
pip install -e ".[dev]"
```

3. **Configuration (optional):**

Copy the example environment file and adjust values as needed:

```bash
cp .env.example .env
```

Available settings:

| Variable       | Default                    | Description            |
|----------------|----------------------------|------------------------|
| `DATABASE_URL` | `sqlite:///calculator.db`  | Database connection URL |
| `HOST`         | `0.0.0.0`                 | Server bind address     |
| `PORT`         | `8000`                    | Server port             |
| `APP_NAME`     | `Calculator API`           | Application title       |

## Run

```bash
uvicorn calculator_api.main:app --reload
```

Or using the installed script:

```bash
calculator-api
```

The server runs on `http://localhost:8000`.

## Endpoints

### Health Check

| Method | Endpoint  | Description              |
|--------|-----------|--------------------------|
| GET    | `/health` | Service health status    |

### Calculator Operations

All calculator endpoints accept POST with JSON body `{"a": <number>, "b": <number>}`.

| Endpoint    | Description          |
|-------------|----------------------|
| `/add`      | Add two numbers      |
| `/subtract` | Subtract b from a    |
| `/multiply` | Multiply two numbers |
| `/divide`   | Divide a by b        |

### Calculation History (CRUD)

Store calculations in a SQLite database.

| Method | Endpoint              | Description                  |
|--------|-----------------------|------------------------------|
| POST   | `/calculations`       | Create & store a calculation |
| GET    | `/calculations`       | List all stored calculations |
| GET    | `/calculations/{id}`  | Get a specific calculation   |
| DELETE | `/calculations/{id}`  | Delete a calculation         |

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

## Running Tests

```bash
pytest
```

## API Docs

FastAPI auto-generates docs at:
- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc
