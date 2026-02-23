"""Application factory and entry point for the Calculator API."""

from fastapi import FastAPI

from .config import get_settings
from .routes import calculator, health


def create_app() -> FastAPI:
    """Create and configure the FastAPI application."""
    settings = get_settings()
    application = FastAPI(
        title="Calculator API",
        version=settings.app_version,
    )

    application.include_router(calculator.router)
    application.include_router(health.router)

    return application


app = create_app()


def main() -> None:
    """Run the server."""
    import uvicorn

    settings = get_settings()
    uvicorn.run(
        "calculator_api.main:app",
        host=settings.host,
        port=settings.port,
        reload=True,
    )


if __name__ == "__main__":
    main()
