package main

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 3, 4, 7},
		{"negative numbers", -3, -4, -7},
		{"mixed signs", -3, 4, 1},
		{"zeros", 0, 0, 0},
		{"with zero", 5, 0, 5},
		{"decimals", 1.5, 2.5, 4.0},
		{"large numbers", 1e10, 2e10, 3e10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 10, 3, 7},
		{"negative numbers", -3, -4, 1},
		{"mixed signs", -3, 4, -7},
		{"zeros", 0, 0, 0},
		{"from zero", 0, 5, -5},
		{"subtract zero", 5, 0, 5},
		{"decimals", 5.5, 2.5, 3.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 7, 6, 42},
		{"negative numbers", -3, -4, 12},
		{"mixed signs", -3, 4, -12},
		{"with zero", 5, 0, 0},
		{"both zeros", 0, 0, 0},
		{"with one", 5, 1, 5},
		{"decimals", 2.5, 4.0, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Multiply(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 10, 2, 5},
		{"negative numbers", -10, -2, 5},
		{"mixed signs", -10, 2, -5},
		{"zero numerator", 0, 5, 0},
		{"decimals", 7.5, 2.5, 3.0},
		{"non-integer result", 10, 3, 10.0 / 3.0},
		{"with one", 42, 1, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if err != nil {
				t.Errorf("Divide(%v, %v) returned unexpected error: %v", tt.a, tt.b, err)
				return
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := Divide(1, 0)
	if err == nil {
		t.Error("Divide(1, 0) expected error, got nil")
	}
	if err != ErrDivisionByZero {
		t.Errorf("Divide(1, 0) error = %v, want %v", err, ErrDivisionByZero)
	}
}

func TestDivideByZeroWithNegative(t *testing.T) {
	_, err := Divide(-5, 0)
	if err == nil {
		t.Error("Divide(-5, 0) expected error, got nil")
	}
	if err != ErrDivisionByZero {
		t.Errorf("Divide(-5, 0) error = %v, want %v", err, ErrDivisionByZero)
	}
}

func TestDivideByZeroWithZero(t *testing.T) {
	_, err := Divide(0, 0)
	if err == nil {
		t.Error("Divide(0, 0) expected error, got nil")
	}
	if err != ErrDivisionByZero {
		t.Errorf("Divide(0, 0) error = %v, want %v", err, ErrDivisionByZero)
	}
}
