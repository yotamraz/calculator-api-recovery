"""FastAPI application creation and ASGI entrypoint."""

from contextlib import asynccontextmanager
from typing import AsyncGenerator

import uvicorn
from fastapi import FastAPI

from .config import settings
from .database import init_db
from .routes import calculator, health


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    """Create database tables on startup."""
    init_db()
    yield


def create_app() -> FastAPI:
    """Create and configure the FastAPI application."""
    app = FastAPI(title=settings.app_name, version="0.1.0", lifespan=lifespan)

    app.include_router(health.router)
    app.include_router(calculator.router)

    return app


app = create_app()


def main() -> None:
    """Run the server."""
    uvicorn.run(app, host=settings.host, port=settings.port)


if __name__ == "__main__":
    main()
