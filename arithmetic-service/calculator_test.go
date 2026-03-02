package main

import (
	"errors"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 5, 3, 8},
		{"negative numbers", -5, -3, -8},
		{"mixed signs", -5, 3, -2},
		{"zeros", 0, 0, 0},
		{"decimals", 1.5, 2.5, 4.0},
		{"large numbers", 1e15, 1e15, 2e15},
		{"small decimals", 0.1, 0.2, 0.30000000000000004}, // IEEE 754
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 5, 3, 2},
		{"negative numbers", -5, -3, -2},
		{"result negative", 3, 5, -2},
		{"zeros", 0, 0, 0},
		{"decimals", 5.5, 2.5, 3.0},
		{"same number", 7, 7, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subtract(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 5, 3, 15},
		{"negative numbers", -5, -3, 15},
		{"mixed signs", -5, 3, -15},
		{"by zero", 5, 0, 0},
		{"by one", 5, 1, 5},
		{"decimals", 2.5, 4, 10},
		{"large numbers", 1e10, 1e10, 1e20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 10, 2, 5},
		{"negative numbers", -10, -2, 5},
		{"mixed signs", -10, 2, -5},
		{"result decimal", 7, 2, 3.5},
		{"by one", 5, 1, 5},
		{"zero dividend", 0, 5, 0},
		{"large numbers", 1e20, 1e10, 1e10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)
			if err != nil {
				t.Errorf("Divide(%v, %v) returned unexpected error: %v", tt.a, tt.b, err)
			}
			if result != tt.expected {
				t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := Divide(10, 0)
	if err == nil {
		t.Fatal("Divide(10, 0) expected error, got nil")
	}
	if !errors.Is(err, ErrDivideByZero) {
		t.Errorf("Divide(10, 0) error = %v, want ErrDivideByZero", err)
	}
	if err.Error() != "Cannot divide by zero" {
		t.Errorf("error message = %q, want %q", err.Error(), "Cannot divide by zero")
	}
}

func TestDivideByZeroNegative(t *testing.T) {
	_, err := Divide(-5, 0)
	if err == nil {
		t.Fatal("Divide(-5, 0) expected error, got nil")
	}
	if !errors.Is(err, ErrDivideByZero) {
		t.Errorf("expected ErrDivideByZero, got %v", err)
	}
}

func TestDivideZeroByZero(t *testing.T) {
	_, err := Divide(0, 0)
	if err == nil {
		t.Fatal("Divide(0, 0) expected error, got nil")
	}
	if !errors.Is(err, ErrDivideByZero) {
		t.Errorf("expected ErrDivideByZero, got %v", err)
	}
}

func TestAddCommutativity(t *testing.T) {
	a, b := 3.14, 2.71
	if Add(a, b) != Add(b, a) {
		t.Error("Add should be commutative")
	}
}

func TestMultiplyCommutativity(t *testing.T) {
	a, b := 3.14, 2.71
	if Multiply(a, b) != Multiply(b, a) {
		t.Error("Multiply should be commutative")
	}
}

func TestDivideSpecialValues(t *testing.T) {
	// Inf / Inf should return error
	_, err := Divide(math.Inf(1), math.Inf(1))
	if err == nil {
		t.Error("Divide(Inf, Inf) should return error")
	}
}
