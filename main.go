package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Read configuration from environment with defaults.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "calculator.db"
	}

	// Initialize database.
	db, err := NewDB(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Build router and register routes.
	r := chi.NewRouter()
	RegisterRoutes(r, db)

	// Start HTTP server.
	addr := "0.0.0.0:" + port
	log.Printf("Starting Calculator API server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
