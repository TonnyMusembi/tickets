package config

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// DBConnection establishes and returns a MySQL database connection
func DBConnection() (*sql.DB, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Construct DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Open database connection
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		slog.Error("Error opening database connection", "error", err)
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		slog.Error("Error pinging database", "error", err)
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	slog.Info("Database connection established successfully")
	
	return db, nil

}