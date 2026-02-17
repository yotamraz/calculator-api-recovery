"""Tests for server module: Settings, models, engine, and OPERATIONS."""

from datetime import UTC, datetime

from core import add, divide, multiply, subtract
from server import (
    Calculation,
    CalculationBase,
    CalculationCreate,
    CalculationRequest,
    CalculationResponse,
    HealthResponse,
    OPERATIONS,
    ResultResponse,
    Settings,
    engine,
    settings,
)


class TestSettings:
    def test_default_database_url(self):
        assert settings.database_url == "sqlite:///calculator.db"

    def test_default_server_host(self):
        assert settings.server_host == "0.0.0.0"

    def test_default_server_port(self):
        assert settings.server_port == 8000

    def test_settings_env_prefix(self):
        s = Settings()
        assert s.model_config["env_prefix"] == "CALC_"


class TestModels:
    def test_calculation_base_fields(self):
        calc = CalculationBase(operation="add", a=1.0, b=2.0, result=3.0)
        assert calc.operation == "add"
        assert calc.a == 1.0
        assert calc.b == 2.0
        assert calc.result == 3.0

    def test_calculation_has_created_at(self):
        calc = Calculation(operation="add", a=1.0, b=2.0, result=3.0)
        assert calc.created_at is not None
        assert calc.created_at.tzinfo is not None

    def test_calculation_created_at_uses_utc(self):
        before = datetime.now(UTC)
        calc = Calculation(operation="add", a=1.0, b=2.0, result=3.0)
        after = datetime.now(UTC)
        assert before <= calc.created_at <= after

    def test_calculation_default_id_is_none(self):
        calc = Calculation(operation="add", a=1.0, b=2.0, result=3.0)
        assert calc.id is None

    def test_calculation_create_fields(self):
        req = CalculationCreate(operation="mul", a=3.0, b=4.0)
        assert req.operation == "mul"
        assert req.a == 3.0
        assert req.b == 4.0

    def test_calculation_request_fields(self):
        req = CalculationRequest(a=5.0, b=10.0)
        assert req.a == 5.0
        assert req.b == 10.0

    def test_result_response_fields(self):
        resp = ResultResponse(result=42.0)
        assert resp.result == 42.0

    def test_health_response_fields(self):
        resp = HealthResponse(status="ok", version="0.1.0")
        assert resp.status == "ok"
        assert resp.version == "0.1.0"

    def test_calculation_response_fields(self):
        resp = CalculationResponse(
            operation="add", a=1.0, b=2.0, result=3.0,
            id=1, created_at=datetime.now(UTC),
        )
        assert resp.id == 1
        assert resp.operation == "add"


class TestEngine:
    def test_engine_url_matches_settings(self):
        assert str(engine.url) == settings.database_url


class TestOperations:
    def test_operations_keys(self):
        assert set(OPERATIONS.keys()) == {"add", "sub", "mul", "div"}

    def test_operations_add(self):
        assert OPERATIONS["add"] is add

    def test_operations_sub(self):
        assert OPERATIONS["sub"] is subtract

    def test_operations_mul(self):
        assert OPERATIONS["mul"] is multiply

    def test_operations_div(self):
        assert OPERATIONS["div"] is divide
