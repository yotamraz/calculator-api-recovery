package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// getenvWithDefault returns the value of the environment variable named by the
// key, or defaultVal if the variable is not set or empty.
func getenvWithDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func main() {
	dbPath := getenvWithDefault("DB_PATH", "calculator.db")
	addr := getenvWithDefault("SERVER_ADDR", "0.0.0.0:8000")

	db, err := InitDB(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	r := gin.Default()

	// Inject the database handle into every request context.
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	registerRoutes(r)

	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
