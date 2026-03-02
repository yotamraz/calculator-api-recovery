package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newPostRequest creates a POST request with JSON body for testing.
func newPostRequest(t *testing.T, path string, body any) *http.Request {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// assertStatus checks the HTTP status code matches expected.
func assertStatus(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if rr.Code != expected {
		t.Errorf("status = %d, want %d; body: %s", rr.Code, expected, rr.Body.String())
	}
}

// assertContentType checks the Content-Type header is application/json.
func assertContentType(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()
	ct := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
}

func TestHandleAddSuccess(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"integers", 5, 3, 8},
		{"decimals", 1.5, 2.5, 4},
		{"negative", -5, 3, -2},
		{"zeros", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newPostRequest(t, "/add", CalculationRequest{A: tt.a, B: tt.b})
			rr := httptest.NewRecorder()
			HandleAdd(rr, req)

			assertStatus(t, rr, http.StatusOK)
			assertContentType(t, rr)

			var resp ResultResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if resp.Result != tt.expected {
				t.Errorf("result = %v, want %v", resp.Result, tt.expected)
			}
		})
	}
}

func TestHandleSubtractSuccess(t *testing.T) {
	req := newPostRequest(t, "/subtract", CalculationRequest{A: 10, B: 3})
	rr := httptest.NewRecorder()
	HandleSubtract(rr, req)

	assertStatus(t, rr, http.StatusOK)
	assertContentType(t, rr)

	var resp ResultResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Result != 7 {
		t.Errorf("result = %v, want 7", resp.Result)
	}
}

func TestHandleMultiplySuccess(t *testing.T) {
	req := newPostRequest(t, "/multiply", CalculationRequest{A: 7, B: 6})
	rr := httptest.NewRecorder()
	HandleMultiply(rr, req)

	assertStatus(t, rr, http.StatusOK)
	assertContentType(t, rr)

	var resp ResultResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Result != 42 {
		t.Errorf("result = %v, want 42", resp.Result)
	}
}

func TestHandleDivideSuccess(t *testing.T) {
	req := newPostRequest(t, "/divide", CalculationRequest{A: 10, B: 2})
	rr := httptest.NewRecorder()
	HandleDivide(rr, req)

	assertStatus(t, rr, http.StatusOK)
	assertContentType(t, rr)

	var resp ResultResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Result != 5 {
		t.Errorf("result = %v, want 5", resp.Result)
	}
}

func TestHandleDivideByZero(t *testing.T) {
	req := newPostRequest(t, "/divide", CalculationRequest{A: 10, B: 0})
	rr := httptest.NewRecorder()
	HandleDivide(rr, req)

	assertStatus(t, rr, http.StatusBadRequest)
	assertContentType(t, rr)

	var resp ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if resp.Detail != "Cannot divide by zero" {
		t.Errorf("detail = %q, want %q", resp.Detail, "Cannot divide by zero")
	}
}

func TestHandleHealthSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	HandleHealth(rr, req)

	assertStatus(t, rr, http.StatusOK)
	assertContentType(t, rr)

	var resp HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("status = %q, want %q", resp.Status, "ok")
	}
	if resp.Version != "0.1.0" {
		t.Errorf("version = %q, want %q", resp.Version, "0.1.0")
	}
}

func TestHandleInvalidJSON(t *testing.T) {
	handlers := []struct {
		name    string
		handler http.HandlerFunc
		path    string
	}{
		{"add", HandleAdd, "/add"},
		{"subtract", HandleSubtract, "/subtract"},
		{"multiply", HandleMultiply, "/multiply"},
		{"divide", HandleDivide, "/divide"},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, h.path, strings.NewReader("not json"))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			h.handler(rr, req)

			assertStatus(t, rr, http.StatusBadRequest)
			assertContentType(t, rr)

			var resp ErrorResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode error response: %v", err)
			}
			if resp.Detail != "Invalid JSON request body" {
				t.Errorf("detail = %q, want %q", resp.Detail, "Invalid JSON request body")
			}
		})
	}
}

func TestHandleWrongMethod(t *testing.T) {
	// Arithmetic endpoints should reject GET
	arithmeticHandlers := []struct {
		name    string
		handler http.HandlerFunc
		path    string
	}{
		{"add", HandleAdd, "/add"},
		{"subtract", HandleSubtract, "/subtract"},
		{"multiply", HandleMultiply, "/multiply"},
		{"divide", HandleDivide, "/divide"},
	}

	for _, h := range arithmeticHandlers {
		t.Run(h.name+"_GET", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, h.path, nil)
			rr := httptest.NewRecorder()
			h.handler(rr, req)

			assertStatus(t, rr, http.StatusMethodNotAllowed)
		})
	}

	// Health endpoint should reject POST
	t.Run("health_POST", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/health", nil)
		rr := httptest.NewRecorder()
		HandleHealth(rr, req)

		assertStatus(t, rr, http.StatusMethodNotAllowed)
	})
}

func TestHandleEmptyBody(t *testing.T) {
	handlers := []struct {
		name    string
		handler http.HandlerFunc
		path    string
	}{
		{"add", HandleAdd, "/add"},
		{"subtract", HandleSubtract, "/subtract"},
		{"multiply", HandleMultiply, "/multiply"},
		{"divide", HandleDivide, "/divide"},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, h.path, strings.NewReader(""))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			h.handler(rr, req)

			assertStatus(t, rr, http.StatusBadRequest)
		})
	}
}

func TestHandleResponseFormat(t *testing.T) {
	// Verify the exact JSON format matches what Python FastAPI produces
	// Python returns {"result": 8.0} for 5 + 3
	req := newPostRequest(t, "/add", CalculationRequest{A: 5, B: 3})
	rr := httptest.NewRecorder()
	HandleAdd(rr, req)

	assertStatus(t, rr, http.StatusOK)

	// Decode the raw JSON to check structure
	var raw map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Should have exactly one key: "result"
	if len(raw) != 1 {
		t.Errorf("response has %d keys, want 1; got: %v", len(raw), raw)
	}
	if _, ok := raw["result"]; !ok {
		t.Error("response missing 'result' key")
	}
}

func TestHandleDivideByZeroResponseFormat(t *testing.T) {
	// Verify error response matches FastAPI format: {"detail": "Cannot divide by zero"}
	req := newPostRequest(t, "/divide", CalculationRequest{A: 1, B: 0})
	rr := httptest.NewRecorder()
	HandleDivide(rr, req)

	assertStatus(t, rr, http.StatusBadRequest)

	var raw map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Should have exactly one key: "detail" (matching FastAPI HTTPException format)
	if len(raw) != 1 {
		t.Errorf("error response has %d keys, want 1; got: %v", len(raw), raw)
	}
	if _, ok := raw["detail"]; !ok {
		t.Error("error response missing 'detail' key (expected FastAPI-compatible format)")
	}
}

func TestHandleHealthResponseFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	HandleHealth(rr, req)

	assertStatus(t, rr, http.StatusOK)

	var raw map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Should have exactly two keys: "status" and "version"
	if len(raw) != 2 {
		t.Errorf("health response has %d keys, want 2; got: %v", len(raw), raw)
	}
	if raw["status"] != "ok" {
		t.Errorf("health status = %v, want 'ok'", raw["status"])
	}
	if raw["version"] != "0.1.0" {
		t.Errorf("health version = %v, want '0.1.0'", raw["version"])
	}
}

func TestHandleDecimalResults(t *testing.T) {
	// Test that decimal results are returned correctly
	req := newPostRequest(t, "/divide", CalculationRequest{A: 7, B: 2})
	rr := httptest.NewRecorder()
	HandleDivide(rr, req)

	assertStatus(t, rr, http.StatusOK)

	var resp ResultResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Result != 3.5 {
		t.Errorf("result = %v, want 3.5", resp.Result)
	}
}
