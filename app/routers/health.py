"""Health check endpoint."""

from fastapi import APIRouter, Request

from app.schemas import HealthResponse

router = APIRouter()


@router.get("/health", response_model=HealthResponse)
def health_check(request: Request) -> HealthResponse:
    """Return the health status of the service."""
    return HealthResponse(status="ok", version=request.app.version)
