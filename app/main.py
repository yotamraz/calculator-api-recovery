"""FastAPI application creation, lifespan, and router registration."""

from collections.abc import AsyncIterator
from contextlib import asynccontextmanager
from typing import Callable

import uvicorn
from fastapi import FastAPI
from sqlmodel import SQLModel

from app.config import settings
from app.core import add, divide, multiply, subtract
from app.database import engine
from app.routers import calculator, health

OPERATIONS: dict[str, Callable[[float, float], float]] = {
    "add": add,
    "sub": subtract,
    "mul": multiply,
    "div": divide,
}


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncIterator[None]:
    """Create database tables on startup."""
    SQLModel.metadata.create_all(engine)
    yield


app = FastAPI(
    title=settings.app_title,
    version=settings.app_version,
    lifespan=lifespan,
)

app.include_router(calculator.router)
app.include_router(health.router)


def main() -> None:
    """Run the server."""
    uvicorn.run(app, host=settings.host, port=settings.port)


if __name__ == "__main__":
    main()
