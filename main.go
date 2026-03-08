package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CalculationRequest represents the JSON request body for arithmetic endpoints.
type CalculationRequest struct {
	A float64 `json:"a" binding:"required"`
	B float64 `json:"b" binding:"required"`
}

// ResultResponse represents the JSON response body for successful arithmetic operations.
type ResultResponse struct {
	Result float64 `json:"result"`
}

// ErrorResponse represents the JSON response body for error cases.
type ErrorResponse struct {
	Detail string `json:"detail"`
}

// HealthResponse represents the JSON response body for the health check endpoint.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

const serviceVersion = "0.1.0"

// getEnvOrDefault reads an environment variable or returns a fallback value.
func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// setupRouter creates and configures the Gin engine with all routes.
func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", handleHealth)
	r.POST("/add", handleAdd)
	r.POST("/subtract", handleSubtract)
	r.POST("/multiply", handleMultiply)
	r.POST("/divide", handleDivide)

	return r
}

// handleHealth returns the service health status.
func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Version: serviceVersion,
	})
}

// handleAdd adds two numbers.
func handleAdd(c *gin.Context) {
	var req CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Detail: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultResponse{Result: Add(req.A, req.B)})
}

// handleSubtract subtracts b from a.
func handleSubtract(c *gin.Context) {
	var req CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Detail: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultResponse{Result: Subtract(req.A, req.B)})
}

// handleMultiply multiplies two numbers.
func handleMultiply(c *gin.Context) {
	var req CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Detail: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultResponse{Result: Multiply(req.A, req.B)})
}

// handleDivide divides a by b.
func handleDivide(c *gin.Context) {
	var req CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Detail: err.Error()})
		return
	}

	result, err := Divide(req.A, req.B)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Detail: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultResponse{Result: result})
}

func main() {
	host := getEnvOrDefault("HOST", "0.0.0.0")
	port := getEnvOrDefault("PORT", "8080")

	r := setupRouter()
	r.Run(host + ":" + port)
}
