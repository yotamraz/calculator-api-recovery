package main

import "time"

// Calculation represents a persisted calculation record in the database.
type Calculation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Operation string    `json:"operation"`
	A         float64   `json:"a"`
	B         float64   `json:"b"`
	Result    float64   `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}

// OperandsRequest represents the JSON request body for arithmetic endpoints.
type OperandsRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

// HealthResponse represents the JSON response body for the health check endpoint.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// ErrorResponse represents the JSON error response body, matching FastAPI's HTTPException format.
type ErrorResponse struct {
	Detail string `json:"detail"`
}
