"""Centralized configuration via Pydantic Settings."""

from functools import lru_cache

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""

    database_url: str = "sqlite:///calculator.db"
    host: str = "0.0.0.0"
    port: int = 8000
    app_version: str = "0.1.0"

    model_config = {
        "env_file": ".env",
    }


@lru_cache()
def get_settings() -> Settings:
    """Return cached application settings."""
    return Settings()
