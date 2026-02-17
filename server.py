#!/usr/bin/env python3
"""A very simple calculator REST API server."""

from datetime import UTC, datetime

from pydantic import BaseModel
from pydantic_settings import BaseSettings, SettingsConfigDict
from sqlmodel import Field, SQLModel, create_engine

from core import add, divide, multiply, subtract


# --- Settings ---


class Settings(BaseSettings):
    """Application settings with environment variable support."""

    model_config = SettingsConfigDict(env_prefix="CALC_")

    database_url: str = "sqlite:///calculator.db"
    server_host: str = "0.0.0.0"
    server_port: int = 8000


settings = Settings()


# --- Models ---


class CalculationBase(SQLModel):
    """Base model for calculations."""

    operation: str
    a: float
    b: float
    result: float


class Calculation(CalculationBase, table=True):
    """Database model for stored calculations."""

    id: int | None = Field(default=None, primary_key=True)
    created_at: datetime = Field(default_factory=lambda: datetime.now(UTC))


class CalculationCreate(BaseModel):
    """Request body for creating a calculation."""

    operation: str
    a: float
    b: float


class CalculationResponse(CalculationBase):
    """Response model for a calculation."""

    id: int
    created_at: datetime


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


# --- Database Setup ---

engine = create_engine(settings.database_url, echo=False)

# --- Operations ---

OPERATIONS: dict[str, callable] = {  # type: ignore[type-arg]
    "add": add,
    "sub": subtract,
    "mul": multiply,
    "div": divide,
}
