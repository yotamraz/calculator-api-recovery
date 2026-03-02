package main

import (
	"errors"
	"math"
)

// ErrDivideByZero is returned when division by zero is attempted.
var ErrDivideByZero = errors.New("Cannot divide by zero")

// Add returns the sum of a and b.
func Add(a, b float64) float64 {
	return a + b
}

// Subtract returns the difference of a and b (a - b).
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns the product of a and b.
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns the quotient of a and b (a / b).
// Returns ErrDivideByZero if b is zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	result := a / b
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, ErrDivideByZero
	}
	return result, nil
}
