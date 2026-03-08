package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func performRequest(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// --- Health Check Tests ---

func TestHandlerHealth(t *testing.T) {
	router := setupRouter()
	w := performRequest(router, "GET", "/health", "")

	if w.Code != http.StatusOK {
		t.Errorf("GET /health status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse health response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("health status = %q, want %q", resp.Status, "ok")
	}
	if resp.Version != "0.1.0" {
		t.Errorf("health version = %q, want %q", resp.Version, "0.1.0")
	}
}

// --- Add Endpoint Tests ---

func TestHandlerAdd(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		body       string
		wantCode   int
		wantResult float64
	}{
		{"positive numbers", `{"a": 5, "b": 3}`, http.StatusOK, 8},
		{"negative numbers", `{"a": -5, "b": -3}`, http.StatusOK, -8},
		{"decimal numbers", `{"a": 1.5, "b": 2.5}`, http.StatusOK, 4},
		{"zeros", `{"a": 0, "b": 0}`, http.StatusOK, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequest(router, "POST", "/add", tt.body)

			if w.Code != tt.wantCode {
				t.Errorf("POST /add status = %d, want %d", w.Code, tt.wantCode)
			}

			var resp ResultResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if resp.Result != tt.wantResult {
				t.Errorf("POST /add result = %v, want %v", resp.Result, tt.wantResult)
			}
		})
	}
}

// --- Subtract Endpoint Tests ---

func TestHandlerSubtract(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		body       string
		wantCode   int
		wantResult float64
	}{
		{"positive result", `{"a": 10, "b": 3}`, http.StatusOK, 7},
		{"negative result", `{"a": 3, "b": 10}`, http.StatusOK, -7},
		{"same numbers", `{"a": 5, "b": 5}`, http.StatusOK, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequest(router, "POST", "/subtract", tt.body)

			if w.Code != tt.wantCode {
				t.Errorf("POST /subtract status = %d, want %d", w.Code, tt.wantCode)
			}

			var resp ResultResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if resp.Result != tt.wantResult {
				t.Errorf("POST /subtract result = %v, want %v", resp.Result, tt.wantResult)
			}
		})
	}
}

// --- Multiply Endpoint Tests ---

func TestHandlerMultiply(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		body       string
		wantCode   int
		wantResult float64
	}{
		{"positive numbers", `{"a": 5, "b": 3}`, http.StatusOK, 15},
		{"by zero", `{"a": 5, "b": 0}`, http.StatusOK, 0},
		{"decimal numbers", `{"a": 2.5, "b": 4}`, http.StatusOK, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequest(router, "POST", "/multiply", tt.body)

			if w.Code != tt.wantCode {
				t.Errorf("POST /multiply status = %d, want %d", w.Code, tt.wantCode)
			}

			var resp ResultResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if resp.Result != tt.wantResult {
				t.Errorf("POST /multiply result = %v, want %v", resp.Result, tt.wantResult)
			}
		})
	}
}

// --- Divide Endpoint Tests ---

func TestHandlerDivide(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		body       string
		wantCode   int
		wantResult float64
	}{
		{"even division", `{"a": 10, "b": 2}`, http.StatusOK, 5},
		{"decimal result", `{"a": 7, "b": 2}`, http.StatusOK, 3.5},
		{"negative numbers", `{"a": -10, "b": -2}`, http.StatusOK, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequest(router, "POST", "/divide", tt.body)

			if w.Code != tt.wantCode {
				t.Errorf("POST /divide status = %d, want %d", w.Code, tt.wantCode)
			}

			var resp ResultResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if resp.Result != tt.wantResult {
				t.Errorf("POST /divide result = %v, want %v", resp.Result, tt.wantResult)
			}
		})
	}
}

// --- Division by Zero Error Test ---

func TestHandlerDivideByZero(t *testing.T) {
	router := setupRouter()
	w := performRequest(router, "POST", "/divide", `{"a": 10, "b": 0}`)

	if w.Code != http.StatusBadRequest {
		t.Errorf("POST /divide (by zero) status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}

	if resp.Detail != "Cannot divide by zero" {
		t.Errorf("POST /divide (by zero) detail = %q, want %q", resp.Detail, "Cannot divide by zero")
	}
}

// --- Invalid Request Body Tests ---

func TestHandlerInvalidBody(t *testing.T) {
	router := setupRouter()

	endpoints := []string{"/add", "/subtract", "/multiply", "/divide"}

	for _, endpoint := range endpoints {
		t.Run("empty body "+endpoint, func(t *testing.T) {
			w := performRequest(router, "POST", endpoint, "")

			if w.Code != http.StatusBadRequest {
				t.Errorf("POST %s (empty body) status = %d, want %d", endpoint, w.Code, http.StatusBadRequest)
			}
		})

		t.Run("invalid json "+endpoint, func(t *testing.T) {
			w := performRequest(router, "POST", endpoint, "not json")

			if w.Code != http.StatusBadRequest {
				t.Errorf("POST %s (invalid json) status = %d, want %d", endpoint, w.Code, http.StatusBadRequest)
			}
		})

		t.Run("missing fields "+endpoint, func(t *testing.T) {
			w := performRequest(router, "POST", endpoint, `{"a": 5}`)

			if w.Code != http.StatusBadRequest {
				t.Errorf("POST %s (missing fields) status = %d, want %d", endpoint, w.Code, http.StatusBadRequest)
			}
		})
	}
}

// --- Error Response Shape Tests ---

func TestErrorResponseShape(t *testing.T) {
	router := setupRouter()

	// Division by zero should return {"detail": "Cannot divide by zero"}
	w := performRequest(router, "POST", "/divide", `{"a": 1, "b": 0}`)

	var raw map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}

	// Verify only "detail" key is present
	if len(raw) != 1 {
		t.Errorf("Error response has %d fields, want 1", len(raw))
	}

	if _, ok := raw["detail"]; !ok {
		t.Error("Error response missing 'detail' field")
	}

	if detail, ok := raw["detail"].(string); !ok || detail != "Cannot divide by zero" {
		t.Errorf("Error response detail = %v, want %q", raw["detail"], "Cannot divide by zero")
	}
}

// --- Health Response Shape Tests ---

func TestHealthResponseShape(t *testing.T) {
	router := setupRouter()
	w := performRequest(router, "GET", "/health", "")

	var raw map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("Failed to parse health response: %v", err)
	}

	// Verify exact fields: "status" and "version"
	if len(raw) != 2 {
		t.Errorf("Health response has %d fields, want 2", len(raw))
	}

	if _, ok := raw["status"]; !ok {
		t.Error("Health response missing 'status' field")
	}

	if _, ok := raw["version"]; !ok {
		t.Error("Health response missing 'version' field")
	}
}

// --- Result Response Shape Tests ---

func TestResultResponseShape(t *testing.T) {
	router := setupRouter()
	w := performRequest(router, "POST", "/add", `{"a": 1, "b": 2}`)

	var raw map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("Failed to parse result response: %v", err)
	}

	// Verify only "result" key is present
	if len(raw) != 1 {
		t.Errorf("Result response has %d fields, want 1", len(raw))
	}

	if _, ok := raw["result"]; !ok {
		t.Error("Result response missing 'result' field")
	}
}
