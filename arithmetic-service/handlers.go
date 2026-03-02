package main

import (
	"encoding/json"
	"net/http"
)

// Version is the service version, matching the source monolith.
const Version = "0.1.0"

// CalculationRequest represents the JSON request body for arithmetic endpoints.
type CalculationRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

// ResultResponse represents the JSON response body for arithmetic endpoints.
type ResultResponse struct {
	Result float64 `json:"result"`
}

// ErrorResponse represents the JSON error response body, matching FastAPI's HTTPException format.
type ErrorResponse struct {
	Detail string `json:"detail"`
}

// HealthResponse represents the JSON response body for the health check endpoint.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error response matching FastAPI's {"detail": "..."} format.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Detail: message})
}

// decodeRequest decodes a JSON request body into a CalculationRequest.
// Returns false and writes an error response if decoding fails.
func decodeRequest(w http.ResponseWriter, r *http.Request) (CalculationRequest, bool) {
	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON request body")
		return req, false
	}
	return req, true
}

// HandleAdd handles POST /add requests.
func HandleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	req, ok := decodeRequest(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, ResultResponse{Result: Add(req.A, req.B)})
}

// HandleSubtract handles POST /subtract requests.
func HandleSubtract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	req, ok := decodeRequest(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, ResultResponse{Result: Subtract(req.A, req.B)})
}

// HandleMultiply handles POST /multiply requests.
func HandleMultiply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	req, ok := decodeRequest(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, ResultResponse{Result: Multiply(req.A, req.B)})
}

// HandleDivide handles POST /divide requests.
func HandleDivide(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	req, ok := decodeRequest(w, r)
	if !ok {
		return
	}
	result, err := Divide(req.A, req.B)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ResultResponse{Result: result})
}

// HandleHealth handles GET /health requests.
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, HealthResponse{
		Status:  "ok",
		Version: Version,
	})
}
