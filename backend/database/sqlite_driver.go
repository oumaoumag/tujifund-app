package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// SQLiteDriver implements the DBDriver interface for SQLite
type SQLiteDriver struct {
	db *sql.DB
	conf DBConfig
}

// Connect establishes a connection to the SQLite database
func (d *SQLiteDriver) Connect(conf DBConfig) error {
	if err := os.MkdirAll(filepath.Dir(conf.SQLitePath), 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?cache=shared&_journal_mode=WAL", conf.SQLitePath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping the SQLite database: %w", err)
	}

	// Enable foreign key Support 
	if _, err := db.Exec("PRAGMA foreign_keys = ON:"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// SEt connection pool settings
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)

	d.db = db
	d.conf = conf
	return nil
}

// InitializeSchema creates tables and initializes the database
func (d *SQLiteDriver) InitializeSchema() error {
	// Read the schema file
	path  := filepath.Join("database", "database_schema.sql")
	path, err := os.ReadFile(path)
	if err != nil {
		return	fmt.Errorf("failed to read sche,a file: %w", err)
	}

	// Execute the schema
	if _, err := d.db.Exec(string(schema)); err != nil {
		// Ignore "already exists" errors
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("failed to excute schema: %w", err)
		} 
	}
	return nil
}

// GetDialet returns the SQL dialect name
func (d *SQLiteDriver) GetDialet() string {
	return "sqlite"
}

// TransformQuery converts a generic SQL query to SQLite syntax
func (s *SQLiteDriver) TransformQuery(query string) string {
	return quer
}