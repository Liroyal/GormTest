package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
)

// MigrationRunner handles database migrations
type MigrationRunner struct {
	db *sql.DB
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *sql.DB) *MigrationRunner {
	return &MigrationRunner{db: db}
}

// RunMigrations executes all pending migrations
func (m *MigrationRunner) RunMigrations(migrationsPath string) error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to list migration files: %w", err)
	}

	sort.Strings(files)

	// Run each migration
	for _, file := range files {
		if err := m.runMigration(file); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file, err)
		}
	}

	return nil
}

func (m *MigrationRunner) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := m.db.Exec(query)
	return err
}

func (m *MigrationRunner) runMigration(filePath string) error {
	filename := filepath.Base(filePath)
	
	// Check if migration already applied
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", filename).Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return nil // Already applied
	}

	// Read and execute migration
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if _, err := m.db.Exec(string(content)); err != nil {
		return err
	}

	// Record migration as applied
	_, err = m.db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", filename)
	return err
}
