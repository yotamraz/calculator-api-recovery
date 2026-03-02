package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// Version is the service version, matching the source monolith.
const Version = "0.1.0"

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

// CalculationRequest represents the JSON request body for arithmetic endpoints.
type CalculationRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

// ValidationError represents a single validation error in FastAPI/Pydantic format.
type ValidationError struct {
	Type  string      `json:"type"`
	Loc   []string    `json:"loc"`
	Msg   string      `json:"msg"`
	Input interface{} `json:"input"`
}

// ValidationErrorResponse represents the FastAPI 422 validation error response.
type ValidationErrorResponse struct {
	Detail []ValidationError `json:"detail"`
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

// validateAndDecodeRequest reads the JSON body, validates that both "a" and "b" fields
// are present and numeric, and returns them. Returns false if validation fails
// (the appropriate error response has already been written).
func validateAndDecodeRequest(w http.ResponseWriter, r *http.Request) (float64, float64, bool) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return 0, 0, false
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON request body")
		return 0, 0, false
	}

	var validationErrors []ValidationError
	var a, b float64

	aVal, aExists := raw["a"]
	bVal, bExists := raw["b"]

	if !aExists {
		validationErrors = append(validationErrors, ValidationError{
			Type:  "missing",
			Loc:   []string{"body", "a"},
			Msg:   "Field required",
			Input: raw,
		})
	} else {
		switch v := aVal.(type) {
		case float64:
			a = v
		case string:
			validationErrors = append(validationErrors, ValidationError{
				Type:  "float_parsing",
				Loc:   []string{"body", "a"},
				Msg:   "Input should be a valid number, unable to parse string as a number",
				Input: v,
			})
		default:
			validationErrors = append(validationErrors, ValidationError{
				Type:  "float_parsing",
				Loc:   []string{"body", "a"},
				Msg:   "Input should be a valid number, unable to parse string as a number",
				Input: aVal,
			})
		}
	}

	if !bExists {
		validationErrors = append(validationErrors, ValidationError{
			Type:  "missing",
			Loc:   []string{"body", "b"},
			Msg:   "Field required",
			Input: raw,
		})
	} else {
		switch v := bVal.(type) {
		case float64:
			b = v
		case string:
			validationErrors = append(validationErrors, ValidationError{
				Type:  "float_parsing",
				Loc:   []string{"body", "b"},
				Msg:   "Input should be a valid number, unable to parse string as a number",
				Input: v,
			})
		default:
			validationErrors = append(validationErrors, ValidationError{
				Type:  "float_parsing",
				Loc:   []string{"body", "b"},
				Msg:   "Input should be a valid number, unable to parse string as a number",
				Input: bVal,
			})
		}
	}

	if len(validationErrors) > 0 {
		writeJSON(w, 422, ValidationErrorResponse{Detail: validationErrors})
		return 0, 0, false
	}

	return a, b, true
}

// HandleAdd handles POST /add requests.
func HandleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	a, b, ok := validateAndDecodeRequest(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, ResultResponse{Result: Add(a, b)})
}

// HandleSubtract handles POST /subtract requests.
func HandleSubtract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	a, b, ok := validateAndDecodeRequest(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, ResultResponse{Result: Subtract(a, b)})
}

// HandleMultiply handles POST /multiply requests.
func HandleMultiply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	a, b, ok := validateAndDecodeRequest(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, ResultResponse{Result: Multiply(a, b)})
}

// HandleDivide handles POST /divide requests.
func HandleDivide(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	a, b, ok := validateAndDecodeRequest(w, r)
	if !ok {
		return
	}
	result, err := Divide(a, b)
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
