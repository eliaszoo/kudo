package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Migration represents a database migration
type Migration struct {
	Name string
	SQL  string
}

var migrations = []Migration{
	{
		Name: "001_initial_schema",
		SQL:  readMigrationFile("001_initial_schema.sql"),
	},
	{
		Name: "002_test_data",
		SQL:  readMigrationFile("002_test_data.sql"),
	},
}

func readMigrationFile(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read migration file %s: %v", filename, err)
	}
	return string(content)
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get database connection string
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/reward_system?charset=utf8mb4&parseTime=True&loc=Local"
	}

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Check connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Create migrations table if not exists
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("All migrations completed successfully")
}

func createMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err := db.Exec(query)
	return err
}

func runMigrations(db *sql.DB) error {
	for _, migration := range migrations {
		// Check if migration already applied
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", migration.Name).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %v", migration.Name, err)
		}

		if count > 0 {
			log.Printf("Migration %s already applied, skipping", migration.Name)
			continue
		}

		log.Printf("Applying migration: %s", migration.Name)

		// Execute migration
		_, err = db.Exec(migration.SQL)
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %v", migration.Name, err)
		}

		// Record migration as applied
		_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", migration.Name)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %v", migration.Name, err)
		}

		log.Printf("Migration %s applied successfully", migration.Name)
	}

	return nil
}