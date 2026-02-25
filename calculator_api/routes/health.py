"""Health check endpoint."""

from pydantic import BaseModel
from fastapi import APIRouter

from ..config import settings

router = APIRouter()


class HealthResponse(BaseModel):
    """Response model for the health check endpoint."""

    status: str
    version: str


@router.get("/health", response_model=HealthResponse)
def health_check() -> HealthResponse:
    """Return the health status of the service."""
    return HealthResponse(status="ok", version=settings.app_version)
