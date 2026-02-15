"""Tests for core calculator functions."""

import pytest

from core import add, divide, multiply, subtract


class TestAdd:
    def test_positive_numbers(self):
        assert add(2, 3) == 5

    def test_negative_numbers(self):
        assert add(-1, -2) == -3

    def test_zero(self):
        assert add(0, 0) == 0

    def test_mixed_sign(self):
        assert add(-1, 3) == 2


class TestSubtract:
    def test_positive_numbers(self):
        assert subtract(5, 3) == 2

    def test_negative_result(self):
        assert subtract(3, 5) == -2

    def test_zero(self):
        assert subtract(0, 0) == 0

    def test_negative_numbers(self):
        assert subtract(-1, -2) == 1


class TestMultiply:
    def test_positive_numbers(self):
        assert multiply(4, 3) == 12

    def test_by_zero(self):
        assert multiply(5, 0) == 0

    def test_negative_numbers(self):
        assert multiply(-2, -3) == 6

    def test_mixed_sign(self):
        assert multiply(-2, 3) == -6


class TestDivide:
    def test_positive_numbers(self):
        assert divide(10, 2) == 5.0

    def test_float_result(self):
        assert divide(7, 2) == 3.5

    def test_negative_numbers(self):
        assert divide(-6, -3) == 2.0

    def test_mixed_sign(self):
        assert divide(-6, 3) == -2.0

    def test_divide_by_zero(self):
        with pytest.raises(ValueError, match="Cannot divide by zero"):
            divide(1, 0)
