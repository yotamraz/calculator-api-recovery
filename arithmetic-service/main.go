package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("ARITHMETIC_SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/add", HandleAdd)
	mux.HandleFunc("/subtract", HandleSubtract)
	mux.HandleFunc("/multiply", HandleMultiply)
	mux.HandleFunc("/divide", HandleDivide)
	mux.HandleFunc("/health", HandleHealth)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Arithmetic service starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
