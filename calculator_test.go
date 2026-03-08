package main

import (
	"errors"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 5, 3, 8},
		{"negative numbers", -5, -3, -8},
		{"mixed signs", -5, 3, -2},
		{"zeros", 0, 0, 0},
		{"decimal numbers", 1.5, 2.5, 4.0},
		{"large numbers", 1e15, 2e15, 3e15},
		{"small decimal", 0.1, 0.2, 0.30000000000000004},
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
		{"positive result", 10, 3, 7},
		{"negative result", 3, 10, -7},
		{"same numbers", 5, 5, 0},
		{"zeros", 0, 0, 0},
		{"decimal numbers", 5.5, 2.3, 3.2},
		{"negative numbers", -5, -3, -2},
		{"large numbers", 1e15, 5e14, 5e14},
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
		{"positive numbers", 5, 3, 15},
		{"by zero", 5, 0, 0},
		{"by one", 5, 1, 5},
		{"negative numbers", -5, -3, 15},
		{"mixed signs", -5, 3, -15},
		{"decimal numbers", 2.5, 4, 10},
		{"large numbers", 1e7, 1e7, 1e14},
		{"small numbers", 0.001, 0.001, 0.000001},
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
		name    string
		a, b    float64
		want    float64
		wantErr bool
	}{
		{"positive numbers", 10, 2, 5, false},
		{"with remainder", 10, 3, 10.0 / 3.0, false},
		{"negative numbers", -10, -2, 5, false},
		{"mixed signs", -10, 2, -5, false},
		{"decimal numbers", 7.5, 2.5, 3, false},
		{"divide by one", 42, 1, 42, false},
		{"zero numerator", 0, 5, 0, false},
		{"large numbers", 1e15, 1e5, 1e10, false},
		{"division by zero", 10, 0, 0, true},
		{"zero by zero", 0, 0, 0, true},
		{"negative by zero", -5, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Divide(%v, %v) expected error, got nil", tt.a, tt.b)
				}
				if !errors.Is(err, ErrDivisionByZero) {
					t.Errorf("Divide(%v, %v) error = %v, want ErrDivisionByZero", tt.a, tt.b, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Divide(%v, %v) unexpected error: %v", tt.a, tt.b, err)
				return
			}

			if math.Abs(got-tt.want) > 1e-10 {
				t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
