package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB

// InitDB initializes the database connection and ensures schema exists
func InitDB() error {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
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

	// Initialize schema
	if err = initializeSchema(); err != nil {
		return fmt.Errorf("error initializing schema: %v", err)
	}

	return nil
}

// initializeSchema creates necessary tables if they don't exist
func initializeSchema() error {
	// Begin transaction
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback() // Will rollback if commit doesn't happen

	createTableSQL := `
    -- Create sequence if it doesn't exist
    CREATE SEQUENCE IF NOT EXISTS newsletter_subscriptions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MAXVALUE
    NO CYCLE;
    
    -- Create newsletter subscriptions table
    CREATE TABLE IF NOT EXISTS newsletter_subscriptions (
        id INTEGER PRIMARY KEY DEFAULT nextval('newsletter_subscriptions_id_seq'),
        email VARCHAR(255) UNIQUE NOT NULL,
        registration_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        verified BOOLEAN DEFAULT FALSE
    );
    
    -- Reset sequence to max id or 1 if table is empty
    SELECT setval('newsletter_subscriptions_id_seq', 
        COALESCE((SELECT MAX(id) FROM newsletter_subscriptions), 1), false);

    -- Create index on email for faster lookups
    CREATE INDEX IF NOT EXISTS idx_newsletter_email ON newsletter_subscriptions(email);
    `

	// Execute schema creation within transaction
	if _, err := tx.Exec(createTableSQL); err != nil {
		return fmt.Errorf("error creating schema: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing schema changes: %v", err)
	}

	log.Println("Successfully initialized database schema")
	return nil
}
