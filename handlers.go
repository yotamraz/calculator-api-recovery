package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// RegisterRoutes registers all HTTP routes on the given Chi router.
func RegisterRoutes(r chi.Router, db *gorm.DB) {
	r.Get("/health", HealthHandler())

	r.Post("/add", ArithmeticHandler(db, "add"))
	r.Post("/subtract", ArithmeticHandler(db, "subtract"))
	r.Post("/multiply", ArithmeticHandler(db, "multiply"))
	r.Post("/divide", ArithmeticHandler(db, "divide"))

	r.Route("/calculations", func(r chi.Router) {
		r.Get("/", ListCalculationsHandler(db))
		r.Get("/{id}", GetCalculationHandler(db))
		r.Delete("/{id}", DeleteCalculationHandler(db))
	})
}

// HealthHandler returns an http.HandlerFunc that responds with the service health status.
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, HealthResponse{
			Status:  "ok",
			Version: "0.1.0",
		})
	}
}

// ArithmeticHandler returns an http.HandlerFunc that performs the named arithmetic
// operation, persists the result, and returns the full Calculation record.
func ArithmeticHandler(db *gorm.DB, operation string) http.HandlerFunc {
	// Build a dispatch map for the core functions.
	type mathFunc func(a, b float64) (float64, error)

	ops := map[string]mathFunc{
		"add":      func(a, b float64) (float64, error) { return Add(a, b), nil },
		"subtract": func(a, b float64) (float64, error) { return Subtract(a, b), nil },
		"multiply": func(a, b float64) (float64, error) { return Multiply(a, b), nil },
		"divide":   Divide,
	}

	fn, ok := ops[operation]
	if !ok {
		// Programming error – will panic at startup if mis-configured.
		log.Fatalf("ArithmeticHandler: unknown operation %q", operation)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req OperandsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		result, err := fn(req.A, req.B)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		calc := Calculation{
			Operation: operation,
			A:         req.A,
			B:         req.B,
			Result:    result,
		}
		if err := db.Create(&calc).Error; err != nil {
			log.Printf("Error persisting calculation: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to save calculation")
			return
		}

		respondJSON(w, http.StatusCreated, calc)
	}
}

// ListCalculationsHandler returns an http.HandlerFunc that lists all calculations
// ordered by created_at descending.
func ListCalculationsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var calculations []Calculation
		if err := db.Order("created_at desc").Find(&calculations).Error; err != nil {
			log.Printf("Error listing calculations: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to list calculations")
			return
		}

		respondJSON(w, http.StatusOK, calculations)
	}
}

// GetCalculationHandler returns an http.HandlerFunc that retrieves a single
// calculation by its ID.
func GetCalculationHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid calculation ID")
			return
		}

		var calc Calculation
		if err := db.First(&calc, id).Error; err != nil {
			respondError(w, http.StatusNotFound, "Calculation not found")
			return
		}

		respondJSON(w, http.StatusOK, calc)
	}
}

// DeleteCalculationHandler returns an http.HandlerFunc that deletes a calculation
// by its ID.
func DeleteCalculationHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid calculation ID")
			return
		}

		result := db.Delete(&Calculation{}, id)
		if result.Error != nil {
			log.Printf("Error deleting calculation: %v", result.Error)
			respondError(w, http.StatusInternalServerError, "Failed to delete calculation")
			return
		}
		if result.RowsAffected == 0 {
			respondError(w, http.StatusNotFound, "Calculation not found")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// respondJSON writes a JSON response with the given status code and payload.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// respondError writes a JSON error response matching the FastAPI HTTPException
// format: {"detail": "<message>"}.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Detail: message})
}
