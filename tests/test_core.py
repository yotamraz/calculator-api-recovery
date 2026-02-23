"""Unit tests for core calculator functions."""

import pytest

from app.core import add, divide, multiply, subtract


class TestAdd:
    """Tests for the add function."""

    def test_add_positive_numbers(self) -> None:
        assert add(2, 3) == 5

    def test_add_negative_numbers(self) -> None:
        assert add(-1, -2) == -3

    def test_add_mixed_numbers(self) -> None:
        assert add(-1, 5) == 4

    def test_add_zeros(self) -> None:
        assert add(0, 0) == 0

    def test_add_floats(self) -> None:
        assert add(1.5, 2.5) == 4.0


class TestSubtract:
    """Tests for the subtract function."""

    def test_subtract_positive_numbers(self) -> None:
        assert subtract(5, 3) == 2

    def test_subtract_negative_result(self) -> None:
        assert subtract(3, 5) == -2

    def test_subtract_negative_numbers(self) -> None:
        assert subtract(-1, -2) == 1

    def test_subtract_zeros(self) -> None:
        assert subtract(0, 0) == 0

    def test_subtract_floats(self) -> None:
        assert subtract(5.5, 2.5) == 3.0


class TestMultiply:
    """Tests for the multiply function."""

    def test_multiply_positive_numbers(self) -> None:
        assert multiply(3, 4) == 12

    def test_multiply_by_zero(self) -> None:
        assert multiply(5, 0) == 0

    def test_multiply_negative_numbers(self) -> None:
        assert multiply(-2, -3) == 6

    def test_multiply_mixed_signs(self) -> None:
        assert multiply(-2, 3) == -6

    def test_multiply_floats(self) -> None:
        assert multiply(2.5, 4.0) == 10.0


class TestDivide:
    """Tests for the divide function."""

    def test_divide_positive_numbers(self) -> None:
        assert divide(10, 2) == 5

    def test_divide_with_remainder(self) -> None:
        assert divide(7, 2) == 3.5

    def test_divide_negative_numbers(self) -> None:
        assert divide(-6, -3) == 2

    def test_divide_mixed_signs(self) -> None:
        assert divide(-6, 3) == -2

    def test_divide_by_zero_raises_error(self) -> None:
        with pytest.raises(ValueError, match="Cannot divide by zero"):
            divide(5, 0)

    def test_divide_zero_by_number(self) -> None:
        assert divide(0, 5) == 0
