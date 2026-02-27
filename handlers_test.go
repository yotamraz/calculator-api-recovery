package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// setupTestRouter creates an in-memory SQLite database and returns a configured
// Chi router and the GORM DB instance for use in tests.
func setupTestRouter(t *testing.T) (*chi.Mux, *gorm.DB) {
	t.Helper()

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	r := chi.NewRouter()
	RegisterRoutes(r, db)
	return r, db
}

// doRequest is a test helper that sends an HTTP request to the router and returns
// the response recorder.
func doRequest(r *chi.Mux, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

// --- Health Endpoint Tests ---

func TestHealthEndpoint(t *testing.T) {
	r, _ := setupTestRouter(t)

	rr := doRequest(r, http.MethodGet, "/health", nil)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /health status = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("Health status = %q, want %q", resp.Status, "ok")
	}
	if resp.Version != "0.1.0" {
		t.Errorf("Health version = %q, want %q", resp.Version, "0.1.0")
	}
}

// --- Arithmetic Endpoint Tests ---

func TestArithmeticAdd(t *testing.T) {
	r, _ := setupTestRouter(t)

	body := OperandsRequest{A: 3, B: 4}
	rr := doRequest(r, http.MethodPost, "/add", body)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /add status = %d, want %d", rr.Code, http.StatusCreated)
	}

	var calc Calculation
	if err := json.Unmarshal(rr.Body.Bytes(), &calc); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if calc.Result != 7 {
		t.Errorf("Add result = %v, want %v", calc.Result, 7.0)
	}
	if calc.Operation != "add" {
		t.Errorf("Operation = %q, want %q", calc.Operation, "add")
	}
	if calc.A != 3 {
		t.Errorf("A = %v, want %v", calc.A, 3.0)
	}
	if calc.B != 4 {
		t.Errorf("B = %v, want %v", calc.B, 4.0)
	}
	if calc.ID == 0 {
		t.Error("Expected non-zero ID for persisted calculation")
	}
	if calc.CreatedAt.IsZero() {
		t.Error("Expected non-zero CreatedAt for persisted calculation")
	}
}

func TestArithmeticSubtract(t *testing.T) {
	r, _ := setupTestRouter(t)

	body := OperandsRequest{A: 10, B: 3}
	rr := doRequest(r, http.MethodPost, "/subtract", body)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /subtract status = %d, want %d", rr.Code, http.StatusCreated)
	}

	var calc Calculation
	if err := json.Unmarshal(rr.Body.Bytes(), &calc); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if calc.Result != 7 {
		t.Errorf("Subtract result = %v, want %v", calc.Result, 7.0)
	}
	if calc.Operation != "subtract" {
		t.Errorf("Operation = %q, want %q", calc.Operation, "subtract")
	}
}

func TestArithmeticMultiply(t *testing.T) {
	r, _ := setupTestRouter(t)

	body := OperandsRequest{A: 7, B: 6}
	rr := doRequest(r, http.MethodPost, "/multiply", body)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /multiply status = %d, want %d", rr.Code, http.StatusCreated)
	}

	var calc Calculation
	if err := json.Unmarshal(rr.Body.Bytes(), &calc); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if calc.Result != 42 {
		t.Errorf("Multiply result = %v, want %v", calc.Result, 42.0)
	}
	if calc.Operation != "multiply" {
		t.Errorf("Operation = %q, want %q", calc.Operation, "multiply")
	}
}

func TestArithmeticDivide(t *testing.T) {
	r, _ := setupTestRouter(t)

	body := OperandsRequest{A: 10, B: 2}
	rr := doRequest(r, http.MethodPost, "/divide", body)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /divide status = %d, want %d", rr.Code, http.StatusCreated)
	}

	var calc Calculation
	if err := json.Unmarshal(rr.Body.Bytes(), &calc); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if calc.Result != 5 {
		t.Errorf("Divide result = %v, want %v", calc.Result, 5.0)
	}
	if calc.Operation != "divide" {
		t.Errorf("Operation = %q, want %q", calc.Operation, "divide")
	}
}

func TestArithmeticDivideByZero(t *testing.T) {
	r, _ := setupTestRouter(t)

	body := OperandsRequest{A: 1, B: 0}
	rr := doRequest(r, http.MethodPost, "/divide", body)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("POST /divide (by zero) status = %d, want %d", rr.Code, http.StatusBadRequest)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errResp.Detail == "" {
		t.Error("Expected non-empty error detail for division by zero")
	}
}

func TestArithmeticWithNegativeNumbers(t *testing.T) {
	r, _ := setupTestRouter(t)

	body := OperandsRequest{A: -5, B: 3}
	rr := doRequest(r, http.MethodPost, "/add", body)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /add (negative) status = %d, want %d", rr.Code, http.StatusCreated)
	}

	var calc Calculation
	if err := json.Unmarshal(rr.Body.Bytes(), &calc); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if calc.Result != -2 {
		t.Errorf("Add result = %v, want %v", calc.Result, -2.0)
	}
}

func TestArithmeticWithDecimals(t *testing.T) {
	r, _ := setupTestRouter(t)

	body := OperandsRequest{A: 1.5, B: 2.5}
	rr := doRequest(r, http.MethodPost, "/multiply", body)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /multiply (decimals) status = %d, want %d", rr.Code, http.StatusCreated)
	}

	var calc Calculation
	if err := json.Unmarshal(rr.Body.Bytes(), &calc); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if calc.Result != 3.75 {
		t.Errorf("Multiply result = %v, want %v", calc.Result, 3.75)
	}
}

func TestArithmeticInvalidJSON(t *testing.T) {
	r, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("POST /add (invalid JSON) status = %d, want %d", rr.Code, http.StatusBadRequest)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errResp.Detail == "" {
		t.Error("Expected non-empty error detail for invalid JSON")
	}
}

// --- Calculations CRUD Endpoint Tests ---

func TestListCalculationsEmpty(t *testing.T) {
	r, _ := setupTestRouter(t)

	rr := doRequest(r, http.MethodGet, "/calculations", nil)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /calculations status = %d, want %d", rr.Code, http.StatusOK)
	}

	var calcs []Calculation
	if err := json.Unmarshal(rr.Body.Bytes(), &calcs); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(calcs) != 0 {
		t.Errorf("Expected empty list, got %d calculations", len(calcs))
	}
}

func TestCreateThenListCalculations(t *testing.T) {
	r, _ := setupTestRouter(t)

	// Create several calculations via arithmetic endpoints.
	doRequest(r, http.MethodPost, "/add", OperandsRequest{A: 1, B: 2})
	doRequest(r, http.MethodPost, "/subtract", OperandsRequest{A: 10, B: 5})
	doRequest(r, http.MethodPost, "/multiply", OperandsRequest{A: 3, B: 4})

	// List calculations.
	rr := doRequest(r, http.MethodGet, "/calculations", nil)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /calculations status = %d, want %d", rr.Code, http.StatusOK)
	}

	var calcs []Calculation
	if err := json.Unmarshal(rr.Body.Bytes(), &calcs); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(calcs) != 3 {
		t.Errorf("Expected 3 calculations, got %d", len(calcs))
	}

	// Verify ordering: latest first (created_at desc).
	if len(calcs) >= 2 {
		for i := 0; i < len(calcs)-1; i++ {
			if calcs[i].CreatedAt.Before(calcs[i+1].CreatedAt) {
				t.Errorf("Calculations not ordered by created_at desc: index %d (%v) before index %d (%v)",
					i, calcs[i].CreatedAt, i+1, calcs[i+1].CreatedAt)
			}
		}
	}
}

func TestCreateThenGetCalculation(t *testing.T) {
	r, _ := setupTestRouter(t)

	// Create a calculation.
	createRR := doRequest(r, http.MethodPost, "/add", OperandsRequest{A: 5, B: 3})
	if createRR.Code != http.StatusCreated {
		t.Fatalf("POST /add status = %d, want %d", createRR.Code, http.StatusCreated)
	}

	var created Calculation
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("Failed to decode creation response: %v", err)
	}

	// Get by ID.
	rr := doRequest(r, http.MethodGet, fmt.Sprintf("/calculations/%d", created.ID), nil)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /calculations/%d status = %d, want %d", created.ID, rr.Code, http.StatusOK)
	}

	var fetched Calculation
	if err := json.Unmarshal(rr.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if fetched.ID != created.ID {
		t.Errorf("Fetched ID = %d, want %d", fetched.ID, created.ID)
	}
	if fetched.Operation != "add" {
		t.Errorf("Fetched operation = %q, want %q", fetched.Operation, "add")
	}
	if fetched.Result != 8 {
		t.Errorf("Fetched result = %v, want %v", fetched.Result, 8.0)
	}
}

func TestGetCalculationNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	rr := doRequest(r, http.MethodGet, "/calculations/999999", nil)

	if rr.Code != http.StatusNotFound {
		t.Errorf("GET /calculations/999999 status = %d, want %d", rr.Code, http.StatusNotFound)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errResp.Detail == "" {
		t.Error("Expected non-empty error detail for not found")
	}
}

func TestDeleteCalculation(t *testing.T) {
	r, _ := setupTestRouter(t)

	// Create a calculation.
	createRR := doRequest(r, http.MethodPost, "/add", OperandsRequest{A: 1, B: 1})
	if createRR.Code != http.StatusCreated {
		t.Fatalf("POST /add status = %d, want %d", createRR.Code, http.StatusCreated)
	}

	var created Calculation
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("Failed to decode creation response: %v", err)
	}

	// Delete the calculation.
	deleteRR := doRequest(r, http.MethodDelete, fmt.Sprintf("/calculations/%d", created.ID), nil)
	if deleteRR.Code != http.StatusNoContent {
		t.Errorf("DELETE /calculations/%d status = %d, want %d", created.ID, deleteRR.Code, http.StatusNoContent)
	}

	// Verify it no longer exists.
	getRR := doRequest(r, http.MethodGet, fmt.Sprintf("/calculations/%d", created.ID), nil)
	if getRR.Code != http.StatusNotFound {
		t.Errorf("GET /calculations/%d after delete status = %d, want %d", created.ID, getRR.Code, http.StatusNotFound)
	}
}

func TestDeleteCalculationNotFound(t *testing.T) {
	r, _ := setupTestRouter(t)

	rr := doRequest(r, http.MethodDelete, "/calculations/999999", nil)

	if rr.Code != http.StatusNotFound {
		t.Errorf("DELETE /calculations/999999 status = %d, want %d", rr.Code, http.StatusNotFound)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errResp.Detail == "" {
		t.Error("Expected non-empty error detail for delete not found")
	}
}
