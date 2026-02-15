#!/usr/bin/env python3
"""A very simple calculator REST API server."""

from contextlib import asynccontextmanager
from typing import AsyncGenerator

import uvicorn
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from pydantic_settings import BaseSettings, SettingsConfigDict

from core import add, divide, multiply, subtract

# --- Settings ---


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""

    model_config = SettingsConfigDict(env_prefix="CALC_")

    database_url: str = "sqlite:///calculator.db"
    server_host: str = "0.0.0.0"
    server_port: int = 8000


settings = Settings()

# --- Models ---


class CalculationRequest(BaseModel):
    """Request body for calculation endpoints."""

    a: float
    b: float


class ResultResponse(BaseModel):
    """Response body for calculation endpoints."""

    result: float


class HealthResponse(BaseModel):
    """Response model for the health check endpoint."""

    status: str
    version: str


# --- App Setup ---


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    """Application lifespan context manager."""
    # Database table creation will be added in Milestone 2
    yield


app = FastAPI(title="Calculator API", version="0.1.0", lifespan=lifespan)


# --- Health Check ---


@app.get("/health", response_model=HealthResponse)
def health_check() -> HealthResponse:
    """Return the health status of the service."""
    return HealthResponse(status="ok", version=app.version)


# --- Calculator Endpoints ---


@app.post("/add", response_model=ResultResponse)
def api_add(req: CalculationRequest) -> ResultResponse:
    """Add two numbers."""
    return ResultResponse(result=add(req.a, req.b))


@app.post("/subtract", response_model=ResultResponse)
def api_subtract(req: CalculationRequest) -> ResultResponse:
    """Subtract b from a."""
    return ResultResponse(result=subtract(req.a, req.b))


@app.post("/multiply", response_model=ResultResponse)
def api_multiply(req: CalculationRequest) -> ResultResponse:
    """Multiply two numbers."""
    return ResultResponse(result=multiply(req.a, req.b))


@app.post("/divide", response_model=ResultResponse)
def api_divide(req: CalculationRequest) -> ResultResponse:
    """Divide a by b."""
    try:
        return ResultResponse(result=divide(req.a, req.b))
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e)) from e


def main() -> None:
    """Run the server."""
    uvicorn.run(app, host=settings.server_host, port=settings.server_port)


if __name__ == "__main__":
    main()
