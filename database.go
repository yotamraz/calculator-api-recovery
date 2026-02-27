package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB opens the SQLite database at dbPath, runs migrations, and returns
// a configured *gorm.DB instance. It creates the database file if it does
// not already exist.
func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Calculation{}); err != nil {
		return nil, err
	}

	return db, nil
}
