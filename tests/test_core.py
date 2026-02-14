"""Unit tests for core calculator functions."""

import pytest

from core import add, divide, multiply, subtract


# --- add ---


class TestAdd:
    def test_positive_numbers(self):
        assert add(2, 3) == 5

    def test_negative_numbers(self):
        assert add(-2, -3) == -5

    def test_mixed_sign(self):
        assert add(-2, 3) == 1

    def test_zero(self):
        assert add(0, 0) == 0

    def test_zero_identity(self):
        assert add(5, 0) == 5
        assert add(0, 5) == 5

    def test_floating_point(self):
        assert add(1.5, 2.5) == 4.0

    def test_large_numbers(self):
        assert add(1e15, 1e15) == 2e15


# --- subtract ---


class TestSubtract:
    def test_positive_numbers(self):
        assert subtract(5, 3) == 2

    def test_negative_numbers(self):
        assert subtract(-5, -3) == -2

    def test_mixed_sign(self):
        assert subtract(-2, 3) == -5

    def test_zero(self):
        assert subtract(0, 0) == 0

    def test_zero_identity(self):
        assert subtract(5, 0) == 5

    def test_result_negative(self):
        assert subtract(3, 5) == -2

    def test_floating_point(self):
        assert subtract(5.5, 2.5) == 3.0

    def test_large_numbers(self):
        assert subtract(1e15, 1e14) == 9e14


# --- multiply ---


class TestMultiply:
    def test_positive_numbers(self):
        assert multiply(2, 3) == 6

    def test_negative_numbers(self):
        assert multiply(-2, -3) == 6

    def test_mixed_sign(self):
        assert multiply(-2, 3) == -6

    def test_zero(self):
        assert multiply(5, 0) == 0
        assert multiply(0, 5) == 0

    def test_identity(self):
        assert multiply(5, 1) == 5

    def test_floating_point(self):
        assert multiply(2.5, 4.0) == 10.0

    def test_large_numbers(self):
        assert multiply(1e7, 1e7) == 1e14


# --- divide ---


class TestDivide:
    def test_positive_numbers(self):
        assert divide(6, 3) == 2.0

    def test_negative_numbers(self):
        assert divide(-6, -3) == 2.0

    def test_mixed_sign(self):
        assert divide(-6, 3) == -2.0

    def test_zero_numerator(self):
        assert divide(0, 5) == 0.0

    def test_identity(self):
        assert divide(5, 1) == 5.0

    def test_floating_point(self):
        assert divide(7.5, 2.5) == 3.0

    def test_large_numbers(self):
        assert divide(1e15, 1e5) == 1e10

    def test_fractional_result(self):
        result = divide(1, 3)
        assert abs(result - 0.3333333333333333) < 1e-10

    def test_division_by_zero_raises_value_error(self):
        with pytest.raises(ValueError, match="Cannot divide by zero"):
            divide(5, 0)

    def test_division_by_zero_with_zero_numerator(self):
        with pytest.raises(ValueError, match="Cannot divide by zero"):
            divide(0, 0)
