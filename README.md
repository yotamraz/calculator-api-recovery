# Calculator API

A simple calculator REST API built with FastAPI, SQLModel, and SQLite.

## Features

- **Arithmetic endpoints** — `POST /add`, `/subtract`, `/multiply`, `/divide`
- **Health check** — `GET /health`
- **Persisted calculations** — CRUD via `/calculations` (coming soon)
- **Configuration** — Environment-based settings via `pydantic-settings`

## Quick Start

### 1. Create a virtual environment

```bash
python -m venv .venv
source .venv/bin/activate  # Linux/macOS
```

### 2. Install dependencies

```bash
pip install -e ".[dev]"
```

### 3. Configure (optional)

Copy the example environment file and adjust as needed:

```bash
cp .env.example .env
```

Available settings (with defaults):

| Variable       | Default                    | Description          |
|---------------|----------------------------|----------------------|
| `DATABASE_URL` | `sqlite:///calculator.db`  | SQLite database URL  |
| `HOST`         | `0.0.0.0`                 | Server bind address  |
| `PORT`         | `8000`                    | Server bind port     |

### 4. Run the server

```bash
uvicorn calculator_api.main:app --reload
```

Or use the installed script:

```bash
calculator-api
```

### 5. Run tests

```bash
pytest
```

## Project Structure

```
calculator_api/
  __init__.py
  main.py          # FastAPI app creation, lifespan, router inclusion
  config.py        # pydantic-settings Settings class
  core.py          # Pure arithmetic business logic
  models.py        # SQLModel ORM models and Pydantic DTOs
  database.py      # Engine, session dependency, table init
  routes/
    __init__.py
    health.py      # GET /health
    calculator.py  # POST /add, /subtract, /multiply, /divide
tests/
pyproject.toml
.env.example
```
