"""Unit tests for calculator_api.core arithmetic functions."""

import pytest

from calculator_api.core import add, divide, multiply, subtract


class TestAdd:
    """Tests for the add function."""

    def test_positive_numbers(self):
        assert add(2, 3) == 5

    def test_negative_numbers(self):
        assert add(-2, -3) == -5

    def test_mixed_signs(self):
        assert add(-2, 3) == 1

    def test_with_zero(self):
        assert add(0, 5) == 5
        assert add(5, 0) == 5

    def test_both_zero(self):
        assert add(0, 0) == 0

    def test_floats(self):
        assert add(1.5, 2.5) == 4.0


class TestSubtract:
    """Tests for the subtract function."""

    def test_positive_numbers(self):
        assert subtract(5, 3) == 2

    def test_negative_result(self):
        assert subtract(3, 5) == -2

    def test_negative_numbers(self):
        assert subtract(-2, -3) == 1

    def test_with_zero(self):
        assert subtract(5, 0) == 5
        assert subtract(0, 5) == -5

    def test_same_numbers(self):
        assert subtract(5, 5) == 0

    def test_floats(self):
        assert subtract(5.5, 2.5) == 3.0


class TestMultiply:
    """Tests for the multiply function."""

    def test_positive_numbers(self):
        assert multiply(3, 4) == 12

    def test_negative_numbers(self):
        assert multiply(-3, -4) == 12

    def test_mixed_signs(self):
        assert multiply(-3, 4) == -12

    def test_with_zero(self):
        assert multiply(5, 0) == 0
        assert multiply(0, 5) == 0

    def test_with_one(self):
        assert multiply(5, 1) == 5

    def test_floats(self):
        assert multiply(2.5, 4.0) == 10.0


class TestDivide:
    """Tests for the divide function."""

    def test_positive_numbers(self):
        assert divide(10, 2) == 5.0

    def test_negative_numbers(self):
        assert divide(-10, -2) == 5.0

    def test_mixed_signs(self):
        assert divide(-10, 2) == -5.0

    def test_fractional_result(self):
        assert divide(7, 2) == 3.5

    def test_with_one(self):
        assert divide(5, 1) == 5.0

    def test_divide_zero_by_number(self):
        assert divide(0, 5) == 0.0

    def test_divide_by_zero_raises_value_error(self):
        with pytest.raises(ValueError, match="Cannot divide by zero"):
            divide(10, 0)

    def test_divide_zero_by_zero_raises_value_error(self):
        with pytest.raises(ValueError, match="Cannot divide by zero"):
            divide(0, 0)

    def test_floats(self):
        assert divide(7.5, 2.5) == 3.0
