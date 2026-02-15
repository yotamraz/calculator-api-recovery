"""Tests for API endpoints."""

from fastapi.testclient import TestClient
from server import app


client = TestClient(app)


def test_health():
    resp = client.get("/health")
    assert resp.status_code == 200
    data = resp.json()
    assert data["status"] == "ok"
    assert data["version"] == "0.1.0"


def test_add():
    resp = client.post("/add", json={"a": 10.5, "b": 5.3})
    assert resp.status_code == 200
    assert resp.json()["result"] == 15.8


def test_subtract():
    resp = client.post("/subtract", json={"a": 20.0, "b": 8.5})
    assert resp.status_code == 200
    assert resp.json()["result"] == 11.5


def test_multiply():
    resp = client.post("/multiply", json={"a": 6.0, "b": 7.0})
    assert resp.status_code == 200
    assert resp.json()["result"] == 42.0


def test_divide():
    resp = client.post("/divide", json={"a": 100.0, "b": 4.0})
    assert resp.status_code == 200
    assert resp.json()["result"] == 25.0


def test_divide_by_zero():
    resp = client.post("/divide", json={"a": 10.0, "b": 0.0})
    assert resp.status_code == 400
    assert resp.json()["detail"] == "Cannot divide by zero"
