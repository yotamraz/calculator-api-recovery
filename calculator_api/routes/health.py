"""Health check endpoint."""

from fastapi import APIRouter

from ..models import HealthResponse

router = APIRouter()


@router.get("/health", response_model=HealthResponse)
def health_check() -> HealthResponse:
    """Return the health status of the service."""
    return HealthResponse(status="ok", version="0.1.0")
