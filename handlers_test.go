package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// setupTestRouter creates a Gin engine configured for testing with an
// in-memory SQLite database.
func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, err := InitDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	registerRoutes(r)
	return r
}

// performRequest is a helper that sends an HTTP request to the test router
// and returns the recorded response.
func performRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// --- Health Check Tests ---

func TestHealthEndpoint(t *testing.T) {
	r := setupTestRouter(t)

	w := performRequest(r, http.MethodGet, "/health", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("health status = %q, want %q", resp.Status, "ok")
	}
	if resp.Version != "0.1.0" {
		t.Errorf("health version = %q, want %q", resp.Version, "0.1.0")
	}
}

// --- Calculator Endpoint Tests ---

func TestAddEndpoint(t *testing.T) {
	r := setupTestRouter(t)

	tests := []struct {
		name   string
		a, b   float64
		want   float64
		status int
	}{
		{"positive numbers", 2, 3, 5, http.StatusOK},
		{"negative numbers", -2, -3, -5, http.StatusOK},
		{"zeros", 0, 0, 0, http.StatusOK},
		{"decimals", 1.5, 2.5, 4, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]float64{"a": tt.a, "b": tt.b}
			w := performRequest(r, http.MethodPost, "/add", body)

			if w.Code != tt.status {
				t.Fatalf("POST /add status = %d, want %d, body: %s", w.Code, tt.status, w.Body.String())
			}

			var resp ResultResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}
			if resp.Result != tt.want {
				t.Errorf("result = %v, want %v", resp.Result, tt.want)
			}
		})
	}
}

func TestSubtractEndpoint(t *testing.T) {
	r := setupTestRouter(t)

	body := map[string]float64{"a": 10, "b": 3}
	w := performRequest(r, http.MethodPost, "/subtract", body)

	if w.Code != http.StatusOK {
		t.Fatalf("POST /subtract status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp ResultResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Result != 7 {
		t.Errorf("result = %v, want %v", resp.Result, 7.0)
	}
}

func TestMultiplyEndpoint(t *testing.T) {
	r := setupTestRouter(t)

	body := map[string]float64{"a": 4, "b": 5}
	w := performRequest(r, http.MethodPost, "/multiply", body)

	if w.Code != http.StatusOK {
		t.Fatalf("POST /multiply status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp ResultResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Result != 20 {
		t.Errorf("result = %v, want %v", resp.Result, 20.0)
	}
}

func TestDivideEndpoint(t *testing.T) {
	r := setupTestRouter(t)

	body := map[string]float64{"a": 10, "b": 4}
	w := performRequest(r, http.MethodPost, "/divide", body)

	if w.Code != http.StatusOK {
		t.Fatalf("POST /divide status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp ResultResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Result != 2.5 {
		t.Errorf("result = %v, want %v", resp.Result, 2.5)
	}
}

func TestDivideByZero(t *testing.T) {
	r := setupTestRouter(t)

	body := map[string]float64{"a": 10, "b": 0}
	w := performRequest(r, http.MethodPost, "/divide", body)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("POST /divide (by zero) status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}
	if resp["detail"] != "Cannot divide by zero" {
		t.Errorf("error detail = %q, want %q", resp["detail"], "Cannot divide by zero")
	}
}

// --- Validation Error Tests ---

func TestCalculatorMissingFields(t *testing.T) {
	r := setupTestRouter(t)

	// Send empty body
	w := performRequest(r, http.MethodPost, "/add", map[string]interface{}{})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("POST /add (empty body) status = %d, want %d, body: %s", w.Code, http.StatusBadRequest, w.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}
	if resp["detail"] == "" {
		t.Error("expected non-empty error detail for missing fields")
	}
}

func TestCalculatorMalformedJSON(t *testing.T) {
	r := setupTestRouter(t)

	// Send malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("POST /add (malformed JSON) status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCalculatorWithZeroValues(t *testing.T) {
	r := setupTestRouter(t)

	// Zero is a valid value and should not trigger validation error
	body := map[string]float64{"a": 0, "b": 5}
	w := performRequest(r, http.MethodPost, "/add", body)

	if w.Code != http.StatusOK {
		t.Fatalf("POST /add (a=0) status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp ResultResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Result != 5 {
		t.Errorf("result = %v, want %v", resp.Result, 5.0)
	}
}
