"""FastAPI application factory and ASGI entrypoint."""

from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager

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
    application = FastAPI(
        title=settings.app_name,
        version="0.1.0",
        lifespan=lifespan,
    )

    application.include_router(health.router)
    application.include_router(calculator.router)

    return application


app = create_app()


def run() -> None:
    """Run the server using Uvicorn."""
    uvicorn.run(app, host=settings.host, port=settings.port)


if __name__ == "__main__":
    run()
