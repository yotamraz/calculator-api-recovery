package main

import "time"

// CalculationRequest is the request body for calculator endpoints (/add, /subtract, /multiply, /divide).
// Pointer types are used so that zero is a valid value (binding:"required" rejects zero-value non-pointers).
type CalculationRequest struct {
	A *float64 `json:"a" binding:"required"`
	B *float64 `json:"b" binding:"required"`
}

// ResultResponse is the response body for calculator endpoints.
type ResultResponse struct {
	Result float64 `json:"result"`
}

// HealthResponse is the response body for the health check endpoint.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// Calculation is the GORM database model for stored calculations.
type Calculation struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Operation string    `json:"operation"`
	A         float64   `json:"a"`
	B         float64   `json:"b"`
	Result    float64   `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}

// CalculationCreate is the request body for creating a new calculation via POST /calculations.
// Pointer types are used for numeric fields so that zero is a valid value.
type CalculationCreate struct {
	Operation string   `json:"operation" binding:"required"`
	A         *float64 `json:"a" binding:"required"`
	B         *float64 `json:"b" binding:"required"`
}

// CalculationResponse is the response body for calculation CRUD endpoints.
type CalculationResponse struct {
	ID        int       `json:"id"`
	Operation string    `json:"operation"`
	A         float64   `json:"a"`
	B         float64   `json:"b"`
	Result    float64   `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}
