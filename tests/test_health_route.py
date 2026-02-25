"""Integration tests for the /health endpoint."""


def test_health_returns_200(client):
    """GET /health should return 200 OK."""
    response = client.get("/health")
    assert response.status_code == 200


def test_health_response_body(client):
    """GET /health should return status and version."""
    response = client.get("/health")
    data = response.json()
    assert data["status"] == "ok"
    assert data["version"] == "0.1.0"


def test_health_response_keys(client):
    """GET /health response should contain exactly the expected keys."""
    response = client.get("/health")
    data = response.json()
    assert set(data.keys()) == {"status", "version"}
