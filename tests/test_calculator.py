"""Endpoint tests for calculator arithmetic routes."""

from fastapi.testclient import TestClient


class TestAddEndpoint:
    """Tests for the /add endpoint."""

    def test_add_two_numbers(self, client: TestClient) -> None:
        response = client.post("/add", json={"a": 5, "b": 3})
        assert response.status_code == 200
        assert response.json() == {"result": 8.0}

    def test_add_negative_numbers(self, client: TestClient) -> None:
        response = client.post("/add", json={"a": -1, "b": -2})
        assert response.status_code == 200
        assert response.json() == {"result": -3.0}

    def test_add_floats(self, client: TestClient) -> None:
        response = client.post("/add", json={"a": 1.5, "b": 2.5})
        assert response.status_code == 200
        assert response.json() == {"result": 4.0}


class TestSubtractEndpoint:
    """Tests for the /subtract endpoint."""

    def test_subtract_two_numbers(self, client: TestClient) -> None:
        response = client.post("/subtract", json={"a": 10, "b": 3})
        assert response.status_code == 200
        assert response.json() == {"result": 7.0}

    def test_subtract_negative_result(self, client: TestClient) -> None:
        response = client.post("/subtract", json={"a": 3, "b": 10})
        assert response.status_code == 200
        assert response.json() == {"result": -7.0}


class TestMultiplyEndpoint:
    """Tests for the /multiply endpoint."""

    def test_multiply_two_numbers(self, client: TestClient) -> None:
        response = client.post("/multiply", json={"a": 4, "b": 5})
        assert response.status_code == 200
        assert response.json() == {"result": 20.0}

    def test_multiply_by_zero(self, client: TestClient) -> None:
        response = client.post("/multiply", json={"a": 4, "b": 0})
        assert response.status_code == 200
        assert response.json() == {"result": 0.0}


class TestDivideEndpoint:
    """Tests for the /divide endpoint."""

    def test_divide_two_numbers(self, client: TestClient) -> None:
        response = client.post("/divide", json={"a": 10, "b": 2})
        assert response.status_code == 200
        assert response.json() == {"result": 5.0}

    def test_divide_with_remainder(self, client: TestClient) -> None:
        response = client.post("/divide", json={"a": 7, "b": 2})
        assert response.status_code == 200
        assert response.json() == {"result": 3.5}

    def test_divide_by_zero_returns_400(self, client: TestClient) -> None:
        response = client.post("/divide", json={"a": 5, "b": 0})
        assert response.status_code == 400
        assert response.json() == {"detail": "Cannot divide by zero"}
