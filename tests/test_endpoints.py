"""Endpoint tests for the Calculator API.

Tests arithmetic endpoints (/add, /subtract, /multiply, /divide),
the health endpoint (/health), and CRUD calculation endpoints
(/calculations), using FastAPI's TestClient with in-memory SQLite
for database isolation.
"""

from datetime import datetime

import pytest
from fastapi.testclient import TestClient
from sqlalchemy.pool import StaticPool
from sqlmodel import Session, SQLModel, create_engine

from server import app, get_session


@pytest.fixture(name="client")
def client_fixture():
    """Create a test client with an isolated in-memory SQLite database."""
    test_engine = create_engine(
        "sqlite://",
        echo=False,
        connect_args={"check_same_thread": False},
        poolclass=StaticPool,
    )
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


# --- CRUD Calculation Endpoint Tests ---


class TestCreateCalculation:
    def test_create_valid_add(self, client: TestClient):
        response = client.post(
            "/calculations", json={"operation": "add", "a": 2, "b": 3}
        )
        assert response.status_code == 201
        data = response.json()
        assert data["operation"] == "add"
        assert data["a"] == 2.0
        assert data["b"] == 3.0
        assert data["result"] == 5.0
        assert "id" in data
        assert "created_at" in data

    def test_create_valid_sub(self, client: TestClient):
        response = client.post(
            "/calculations", json={"operation": "sub", "a": 10, "b": 4}
        )
        assert response.status_code == 201
        data = response.json()
        assert data["operation"] == "sub"
        assert data["result"] == 6.0

    def test_create_valid_mul(self, client: TestClient):
        response = client.post(
            "/calculations", json={"operation": "mul", "a": 3, "b": 5}
        )
        assert response.status_code == 201
        assert response.json()["result"] == 15.0

    def test_create_valid_div(self, client: TestClient):
        response = client.post(
            "/calculations", json={"operation": "div", "a": 10, "b": 2}
        )
        assert response.status_code == 201
        assert response.json()["result"] == 5.0

    def test_create_has_valid_id(self, client: TestClient):
        response = client.post(
            "/calculations", json={"operation": "add", "a": 1, "b": 1}
        )
        data = response.json()
        assert isinstance(data["id"], int)
        assert data["id"] > 0

    def test_create_has_valid_created_at(self, client: TestClient):
        response = client.post(
            "/calculations", json={"operation": "add", "a": 1, "b": 1}
        )
        data = response.json()
        # Verify created_at is a valid ISO datetime string
        parsed = datetime.fromisoformat(data["created_at"])
        assert isinstance(parsed, datetime)

    def test_create_unknown_operation_returns_400(self, client: TestClient):
        response = client.post(
            "/calculations", json={"operation": "modulo", "a": 5, "b": 3}
        )
        assert response.status_code == 400

    def test_create_unknown_operation_error_message(self, client: TestClient):
        response = client.post(
            "/calculations", json={"operation": "modulo", "a": 5, "b": 3}
        )
        data = response.json()
        detail = data["detail"]
        assert "modulo" in detail
        # The error message should list valid operations
        assert "add" in detail
        assert "sub" in detail
        assert "mul" in detail
        assert "div" in detail

    def test_create_division_by_zero_returns_400(self, client: TestClient):
        response = client.post(
            "/calculations", json={"operation": "div", "a": 5, "b": 0}
        )
        assert response.status_code == 400

    def test_create_division_by_zero_error_message(self, client: TestClient):
        response = client.post(
            "/calculations", json={"operation": "div", "a": 5, "b": 0}
        )
        data = response.json()
        assert data["detail"] == "Cannot divide by zero"


class TestListCalculations:
    def test_list_empty(self, client: TestClient):
        response = client.get("/calculations")
        assert response.status_code == 200
        assert response.json() == []

    def test_list_includes_created_records(self, client: TestClient):
        client.post("/calculations", json={"operation": "add", "a": 1, "b": 2})
        client.post("/calculations", json={"operation": "sub", "a": 5, "b": 3})
        response = client.get("/calculations")
        assert response.status_code == 200
        data = response.json()
        assert len(data) == 2

    def test_list_ordered_by_created_at_descending(self, client: TestClient):
        client.post("/calculations", json={"operation": "add", "a": 1, "b": 1})
        client.post("/calculations", json={"operation": "sub", "a": 5, "b": 3})
        response = client.get("/calculations")
        data = response.json()
        # Most recently created should be first
        first_created = datetime.fromisoformat(data[0]["created_at"])
        second_created = datetime.fromisoformat(data[1]["created_at"])
        assert first_created >= second_created

    def test_list_contains_expected_fields(self, client: TestClient):
        client.post("/calculations", json={"operation": "mul", "a": 3, "b": 4})
        response = client.get("/calculations")
        data = response.json()
        record = data[0]
        assert "id" in record
        assert "created_at" in record
        assert "operation" in record
        assert "a" in record
        assert "b" in record
        assert "result" in record


class TestGetCalculation:
    def test_get_existing_calculation(self, client: TestClient):
        create_resp = client.post(
            "/calculations", json={"operation": "add", "a": 7, "b": 3}
        )
        calc_id = create_resp.json()["id"]
        response = client.get(f"/calculations/{calc_id}")
        assert response.status_code == 200
        data = response.json()
        assert data["id"] == calc_id
        assert data["operation"] == "add"
        assert data["a"] == 7.0
        assert data["b"] == 3.0
        assert data["result"] == 10.0

    def test_get_nonexistent_calculation_returns_404(self, client: TestClient):
        response = client.get("/calculations/99999")
        assert response.status_code == 404

    def test_get_nonexistent_calculation_error_message(self, client: TestClient):
        response = client.get("/calculations/99999")
        data = response.json()
        assert data["detail"] == "Calculation not found"


class TestDeleteCalculation:
    def test_delete_existing_calculation(self, client: TestClient):
        create_resp = client.post(
            "/calculations", json={"operation": "add", "a": 1, "b": 2}
        )
        calc_id = create_resp.json()["id"]
        response = client.delete(f"/calculations/{calc_id}")
        assert response.status_code == 204

    def test_delete_nonexistent_calculation_returns_404(self, client: TestClient):
        response = client.delete("/calculations/99999")
        assert response.status_code == 404

    def test_delete_nonexistent_calculation_error_message(self, client: TestClient):
        response = client.delete("/calculations/99999")
        data = response.json()
        assert data["detail"] == "Calculation not found"

    def test_deleted_record_not_retrievable(self, client: TestClient):
        create_resp = client.post(
            "/calculations", json={"operation": "sub", "a": 10, "b": 5}
        )
        calc_id = create_resp.json()["id"]
        # Delete the record
        delete_resp = client.delete(f"/calculations/{calc_id}")
        assert delete_resp.status_code == 204
        # Verify it's gone
        get_resp = client.get(f"/calculations/{calc_id}")
        assert get_resp.status_code == 404
        assert get_resp.json()["detail"] == "Calculation not found"
