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
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed signs", -2, 3, 1},
		{"zeros", 0, 0, 0},
		{"large numbers", 1e15, 2e15, 3e15},
		{"decimal numbers", 0.1, 0.2, 0.30000000000000004},
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
		{"positive numbers", 5, 3, 2},
		{"negative result", 3, 5, -2},
		{"negative numbers", -2, -3, 1},
		{"zeros", 0, 0, 0},
		{"subtract from zero", 0, 5, -5},
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
		{"positive numbers", 2, 3, 6},
		{"negative numbers", -2, -3, 6},
		{"mixed signs", -2, 3, -6},
		{"multiply by zero", 5, 0, 0},
		{"multiply by one", 5, 1, 5},
		{"large numbers", 1e10, 1e10, 1e20},
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
		name      string
		a, b      float64
		want      float64
		wantError bool
	}{
		{"positive numbers", 6, 3, 2, false},
		{"negative numbers", -6, -3, 2, false},
		{"mixed signs", -6, 3, -2, false},
		{"result is fraction", 1, 3, 1.0 / 3.0, false},
		{"divide zero", 0, 5, 0, false},
		{"division by zero", 1, 0, 0, true},
		{"zero divided by zero", 0, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if tt.wantError {
				if err == nil {
					t.Errorf("Divide(%v, %v) expected error, got nil", tt.a, tt.b)
				}
				if err != nil && err.Error() != "Cannot divide by zero" {
					t.Errorf("Divide(%v, %v) error = %q, want %q", tt.a, tt.b, err.Error(), "Cannot divide by zero")
				}
			} else {
				if err != nil {
					t.Errorf("Divide(%v, %v) unexpected error: %v", tt.a, tt.b, err)
				}
				if math.Abs(got-tt.want) > 1e-15 {
					t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
				}
			}
		})
	}
}
