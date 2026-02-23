"""Database engine, session factory, and dependency."""

from collections.abc import Generator

from sqlmodel import Session, create_engine

from app.config import settings

engine = create_engine(settings.database_url, echo=False)


def get_session() -> Generator[Session, None, None]:
    """Yield a database session."""
    with Session(engine) as session:
        yield session
