package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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

// Close closes the database connection
func (d *SQLiteDriver) Close() error {
	return d.db.Close()
}

// Ping checks the database connection
func (d *SQLiteDriver) Ping() error {
	return d.db.Ping()
}

// BeginTx starts the a transaction
func (d *SQLiteDriver) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, nil)
}

