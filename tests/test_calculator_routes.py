"""Integration tests for the arithmetic calculator endpoints."""


class TestAddEndpoint:
    """Tests for POST /add."""

    def test_add_positive_numbers(self, client):
        response = client.post("/add", json={"a": 3, "b": 4})
        assert response.status_code == 200
        assert response.json() == {"result": 7.0}

    def test_add_negative_numbers(self, client):
        response = client.post("/add", json={"a": -3, "b": -4})
        assert response.status_code == 200
        assert response.json() == {"result": -7.0}

    def test_add_with_zero(self, client):
        response = client.post("/add", json={"a": 5, "b": 0})
        assert response.status_code == 200
        assert response.json() == {"result": 5.0}

    def test_add_floats(self, client):
        response = client.post("/add", json={"a": 1.5, "b": 2.5})
        assert response.status_code == 200
        assert response.json() == {"result": 4.0}


class TestSubtractEndpoint:
    """Tests for POST /subtract."""

    def test_subtract_positive_numbers(self, client):
        response = client.post("/subtract", json={"a": 10, "b": 3})
        assert response.status_code == 200
        assert response.json() == {"result": 7.0}

    def test_subtract_negative_result(self, client):
        response = client.post("/subtract", json={"a": 3, "b": 10})
        assert response.status_code == 200
        assert response.json() == {"result": -7.0}

    def test_subtract_same_numbers(self, client):
        response = client.post("/subtract", json={"a": 5, "b": 5})
        assert response.status_code == 200
        assert response.json() == {"result": 0.0}


class TestMultiplyEndpoint:
    """Tests for POST /multiply."""

    def test_multiply_positive_numbers(self, client):
        response = client.post("/multiply", json={"a": 5, "b": 6})
        assert response.status_code == 200
        assert response.json() == {"result": 30.0}

    def test_multiply_with_zero(self, client):
        response = client.post("/multiply", json={"a": 5, "b": 0})
        assert response.status_code == 200
        assert response.json() == {"result": 0.0}

    def test_multiply_negative_numbers(self, client):
        response = client.post("/multiply", json={"a": -3, "b": -4})
        assert response.status_code == 200
        assert response.json() == {"result": 12.0}


class TestDivideEndpoint:
    """Tests for POST /divide."""

    def test_divide_positive_numbers(self, client):
        response = client.post("/divide", json={"a": 10, "b": 2})
        assert response.status_code == 200
        assert response.json() == {"result": 5.0}

    def test_divide_fractional_result(self, client):
        response = client.post("/divide", json={"a": 7, "b": 2})
        assert response.status_code == 200
        assert response.json() == {"result": 3.5}

    def test_divide_by_zero_returns_400(self, client):
        response = client.post("/divide", json={"a": 10, "b": 0})
        assert response.status_code == 400
        assert response.json() == {"detail": "Cannot divide by zero"}

    def test_divide_negative_numbers(self, client):
        response = client.post("/divide", json={"a": -10, "b": -2})
        assert response.status_code == 200
        assert response.json() == {"result": 5.0}


class TestInvalidRequests:
    """Tests for invalid request bodies."""

    def test_missing_field(self, client):
        response = client.post("/add", json={"a": 5})
        assert response.status_code == 422

    def test_invalid_type(self, client):
        response = client.post("/add", json={"a": "foo", "b": 5})
        assert response.status_code == 422

    def test_empty_body(self, client):
        response = client.post("/add", json={})
        assert response.status_code == 422
