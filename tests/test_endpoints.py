"""Endpoint tests for the Calculator API.

Tests arithmetic endpoints (/add, /subtract, /multiply, /divide),
the health endpoint (/health), using FastAPI's TestClient with
in-memory SQLite for database isolation.
"""

import pytest
from fastapi.testclient import TestClient
from sqlmodel import Session, SQLModel, create_engine

from server import app, get_session


@pytest.fixture(name="client")
def client_fixture():
    """Create a test client with an isolated in-memory SQLite database."""
    test_engine = create_engine("sqlite://", echo=False)
    SQLModel.metadata.create_all(test_engine)

    def get_test_session():
        with Session(test_engine) as session:
            yield session

    app.dependency_overrides[get_session] = get_test_session
    with TestClient(app) as client:
        yield client
    app.dependency_overrides.clear()


# --- Health Endpoint Tests ---


class TestHealthEndpoint:
    def test_health_returns_200(self, client: TestClient):
        response = client.get("/health")
        assert response.status_code == 200

    def test_health_contains_status_ok(self, client: TestClient):
        response = client.get("/health")
        data = response.json()
        assert data["status"] == "ok"

    def test_health_contains_version(self, client: TestClient):
        response = client.get("/health")
        data = response.json()
        assert "version" in data
        assert isinstance(data["version"], str)
        assert len(data["version"]) > 0


# --- Arithmetic Endpoint Tests ---


class TestAddEndpoint:
    def test_add_positive_numbers(self, client: TestClient):
        response = client.post("/add", json={"a": 2, "b": 3})
        assert response.status_code == 200
        assert response.json()["result"] == 5.0

    def test_add_negative_numbers(self, client: TestClient):
        response = client.post("/add", json={"a": -2, "b": -3})
        assert response.status_code == 200
        assert response.json()["result"] == -5.0

    def test_add_mixed_signs(self, client: TestClient):
        response = client.post("/add", json={"a": -10, "b": 3})
        assert response.status_code == 200
        assert response.json()["result"] == -7.0

    def test_add_zeros(self, client: TestClient):
        response = client.post("/add", json={"a": 0, "b": 0})
        assert response.status_code == 200
        assert response.json()["result"] == 0.0

    def test_add_floats(self, client: TestClient):
        response = client.post("/add", json={"a": 1.5, "b": 2.5})
        assert response.status_code == 200
        assert response.json()["result"] == 4.0


class TestSubtractEndpoint:
    def test_subtract_positive_numbers(self, client: TestClient):
        response = client.post("/subtract", json={"a": 10, "b": 3})
        assert response.status_code == 200
        assert response.json()["result"] == 7.0

    def test_subtract_negative_numbers(self, client: TestClient):
        response = client.post("/subtract", json={"a": -5, "b": -3})
        assert response.status_code == 200
        assert response.json()["result"] == -2.0

    def test_subtract_mixed_signs(self, client: TestClient):
        response = client.post("/subtract", json={"a": -2, "b": 3})
        assert response.status_code == 200
        assert response.json()["result"] == -5.0

    def test_subtract_zeros(self, client: TestClient):
        response = client.post("/subtract", json={"a": 0, "b": 0})
        assert response.status_code == 200
        assert response.json()["result"] == 0.0

    def test_subtract_floats(self, client: TestClient):
        response = client.post("/subtract", json={"a": 5.5, "b": 2.5})
        assert response.status_code == 200
        assert response.json()["result"] == 3.0


class TestMultiplyEndpoint:
    def test_multiply_positive_numbers(self, client: TestClient):
        response = client.post("/multiply", json={"a": 4, "b": 5})
        assert response.status_code == 200
        assert response.json()["result"] == 20.0

    def test_multiply_negative_numbers(self, client: TestClient):
        response = client.post("/multiply", json={"a": -2, "b": -3})
        assert response.status_code == 200
        assert response.json()["result"] == 6.0

    def test_multiply_mixed_signs(self, client: TestClient):
        response = client.post("/multiply", json={"a": -2, "b": 3})
        assert response.status_code == 200
        assert response.json()["result"] == -6.0

    def test_multiply_by_zero(self, client: TestClient):
        response = client.post("/multiply", json={"a": 5, "b": 0})
        assert response.status_code == 200
        assert response.json()["result"] == 0.0

    def test_multiply_floats(self, client: TestClient):
        response = client.post("/multiply", json={"a": 1.5, "b": 2.0})
        assert response.status_code == 200
        assert response.json()["result"] == 3.0


class TestDivideEndpoint:
    def test_divide_positive_numbers(self, client: TestClient):
        response = client.post("/divide", json={"a": 10, "b": 2})
        assert response.status_code == 200
        assert response.json()["result"] == 5.0

    def test_divide_negative_numbers(self, client: TestClient):
        response = client.post("/divide", json={"a": -6, "b": -3})
        assert response.status_code == 200
        assert response.json()["result"] == 2.0

    def test_divide_mixed_signs(self, client: TestClient):
        response = client.post("/divide", json={"a": -6, "b": 3})
        assert response.status_code == 200
        assert response.json()["result"] == -2.0

    def test_divide_floats(self, client: TestClient):
        response = client.post("/divide", json={"a": 7.5, "b": 2.5})
        assert response.status_code == 200
        assert response.json()["result"] == 3.0

    def test_divide_by_zero_returns_400(self, client: TestClient):
        response = client.post("/divide", json={"a": 5, "b": 0})
        assert response.status_code == 400

    def test_divide_by_zero_error_message(self, client: TestClient):
        response = client.post("/divide", json={"a": 5, "b": 0})
        data = response.json()
        assert data["detail"] == "Cannot divide by zero"
