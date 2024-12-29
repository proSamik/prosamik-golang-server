package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
)

var DB *sql.DB

// InitDB initializes the database connection and applies migrations
func InitDB() error {
	// Construct the connection string from environment variables
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	// Open a connection to the PostgreSQL database
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25) // Maximum number of open connections
	DB.SetMaxIdleConns(5)  // Maximum number of idle connections

	// Verify connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Apply migrations
	if err = applyMigrations(); err != nil {
		return fmt.Errorf("error applying migrations: %v", err)
	}

	return nil
}

// applyMigrations applies any pending migrations
func applyMigrations() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %v", err)
	}

	migrationsPath := filepath.Join(wd, "internal", "database", "migrations")
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)

	// Create a database driver instance
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %v", err)
	}

	// Create migration instance using driver
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres", // database name
		driver,
	)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %v", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error applying migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
