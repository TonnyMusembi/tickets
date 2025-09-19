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
	// Try loading .env (only useful for local dev)
	if os.Getenv("GIN_MODE") != "release" {
		_ = godotenv.Load() // ignore missing .env in prod
	}

	// Read env vars
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" || pass == "" || host == "" || port == "" || name == "" {
		return nil, fmt.Errorf("database environment variables are missing")
	}

	// DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, pass, host, port, name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		slog.Error("Error opening DB connection", "error", err)
		return nil, fmt.Errorf("error opening DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		slog.Error("Error pinging DB", "error", err)
		return nil, fmt.Errorf("error pinging DB: %w", err)
	}

	// Connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	slog.Info("Database connection established successfully")
	return db, nil
}
