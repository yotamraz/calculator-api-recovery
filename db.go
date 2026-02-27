package main

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// NewDB opens a GORM SQLite connection at the given DSN and runs AutoMigrate
// for all application models. The DSN is typically a file path (e.g. "calculator.db")
// or ":memory:" for in-memory databases used in testing.
func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Calculation{}); err != nil {
		return nil, err
	}

	return db, nil
}
