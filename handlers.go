package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// registerRoutes wires all HTTP endpoints to the given Gin engine.
func registerRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", healthHandler)

	// Calculator endpoints
	r.POST("/add", addHandler)
	r.POST("/subtract", subtractHandler)
	r.POST("/multiply", multiplyHandler)
	r.POST("/divide", divideHandler)
}

// healthHandler returns the service health status and version.
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Version: "0.1.0",
	})
}

// addHandler adds two numbers and returns the result.
func addHandler(c *gin.Context) {
	var req CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultResponse{Result: Add(*req.A, *req.B)})
}

// subtractHandler subtracts b from a and returns the result.
func subtractHandler(c *gin.Context) {
	var req CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultResponse{Result: Subtract(*req.A, *req.B)})
}

// multiplyHandler multiplies two numbers and returns the result.
func multiplyHandler(c *gin.Context) {
	var req CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultResponse{Result: Multiply(*req.A, *req.B)})
}

// divideHandler divides a by b and returns the result. Returns 400 on division by zero.
func divideHandler(c *gin.Context) {
	var req CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	result, err := Divide(*req.A, *req.B)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultResponse{Result: result})
}
